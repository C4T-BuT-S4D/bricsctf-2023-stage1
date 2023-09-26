#!/usr/bin/env python3
import asyncio
import hashlib
import io
import string
import struct
import subprocess
import sys
from collections import Counter
from typing import Iterator, Iterable

import PIL
from PIL import Image, ImageSequence
from grpclib.client import Channel

import gifparse
from restore_grpc import RestoreServiceStub
from restore_pb2 import Answer

PIL.Image.MAX_IMAGE_PIXELS = 2112322750 + 1


def parse_gif_file(gif: bytes):
    return gifparse.read_file(io.BytesIO(gif))


def restore_gif(gif: bytes) -> bytes:
    cropped_last_ind = parse_gif_file(gif) + 1
    cropped = bytearray(gif[:cropped_last_ind])
    trailer = bytearray(gif[cropped_last_ind:])

    header_end = 0
    for i, b in enumerate(cropped):
        if b == 0xf9 and cropped[i - 1] == 0x21:
            header_end = i - 1
            break
    gif_header = cropped[:header_end]

    pos = len(trailer) - 1
    if trailer[pos] != 0x3b:
        raise ValueError("expecting gif EOF")

    frames = []
    frame_end = pos
    while pos >= 0:
        if trailer[pos] == 0 and trailer[pos + 1] == 0x21 and trailer[pos + 2] == 0xf9:
            frames.append(trailer[pos + 1:frame_end])
            frame_end = pos + 1
        pos -= 1

    mw = 0
    mh = 0
    for frame in frames:
        wh_offset = 13
        w = struct.unpack('<H', frame[wh_offset:wh_offset + 2])[0]
        h = struct.unpack('<H', frame[wh_offset + 2:wh_offset + 4])[0]
        mw = max(mw, w)
        mh = max(mh, h)

    # Change new width and height.
    gif_header[6:8] = struct.pack('<H', mw)
    gif_header[8:10] = struct.pack('<H', mh)

    out = []
    out += gif_header
    for frame in reversed(frames):
        out += frame
    out.append(0x3b)
    return bytes(out)


def prepare_ocr(img: Image) -> Image:
    return noise_remove(img)


def noise_remove(img: Image) -> Image:
    img = img.convert('L')
    w, h = img.size
    c = Counter()
    for x in range(w):
        for y in range(h):
            pxl = img.getpixel((x, y))
            c[pxl] += 1
    pop = c.most_common()[1][0]

    return img.point(lambda x: 0 if x == pop else 255)


def options_(c: str) -> Iterable[str]:
    if c == '_':
        alpha = list(string.ascii_letters + string.digits)
        # 'g' hack (g is most common unrecognized char)
        g_ind = alpha.index('g')
        alpha[0], alpha[g_ind] = alpha[g_ind], alpha[0]
        return list(alpha)

    uq = set()
    sets = ({'0', 'o', 'O'},
            {'1', 'l', 'i', 'I', 'j'},
            {'z', '2'},
            {'9', 'g', 'q'},
            {'b', 'j', 'J'},
            {'v', 'y', 'V', 'Y'},
            {'r', 'f'},
            {'h', 'n'},
            {'5', 's', 'S'},
            {'8', 'B', 'g'})
    for s in sets:
        if c not in s:
            continue
        for v in s:
            uq.add(v)
    uq.add(c.upper())
    uq.add(c.lower())
    return uq


def gen_options(x: str, pos: int) -> Iterator[str]:
    if pos >= len(x):
        yield ''
        return

    opts = options_(x[pos])
    for opt in opts:
        for v in gen_options(x, pos + 1):
            yield opt + v
    return


def brute_answer(s: str, md5h: str, salt: str) -> str:
    for opt in gen_options(s, 0):
        if hashlib.md5(f'{opt}{salt}'.encode()).hexdigest() == md5h:
            return opt
    raise ValueError(f'Failed to find hash for {s}')


def restore(gif: bytes, md5h: str, salt: str) -> str:
    frames_count = Image.open(io.BytesIO(gif)).n_frames
    print(frames_count)
    restored = restore_gif(gif)
    with open('restored.gif', 'wb') as f:
        f.write(restored)
    pil_image = Image.open('restored.gif')
    re_cropped = ImageSequence.all_frames(pil_image, func=prepare_ocr)
    rec_frames_cnt = len(re_cropped)
    print("Recovered frames number: ", rec_frames_cnt)

    rec = ""
    for i, img in enumerate(re_cropped):
        img.save(f'frames/temp_frame_{i}.png')
        frame_rec = subprocess.check_output(
            ["gocr", "-c", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
             f"frames/temp_frame_{i}.png"]).decode()
        frame_rec = frame_rec.strip().replace('\n', '')
        if len(frame_rec) != 1:
            replacement = "_"
            frame_rec = frame_rec.lower()
            if frame_rec == "vv" or frame_rec == 'vn':
                replacement = "w"
            elif frame_rec == 'l11' or frame_rec == 'lll' or frame_rec == 'i11' or frame_rec == 'l1l':
                replacement = 'm'
            elif frame_rec == 'll' or frame_rec == '11':
                replacement = 'n'
            elif frame_rec == 'mj' or frame_rec == 'v1i':
                replacement = 'W'
            print(f"Bad = {frame_rec} replaced with = {replacement}")
            rec += replacement
        else:
            rec += frame_rec
    print("Recognized: ", rec)
    widths, heights = zip(*(i.size for i in re_cropped))
    letter_space = 0
    left_offset = 20
    total_width = sum(widths) + rec_frames_cnt * letter_space + left_offset  # space between images
    max_height = max(heights)
    combined_img = Image.new('RGB', (total_width, max_height + 20), color=(255, 255, 255))
    x_offset = left_offset
    for im in re_cropped:
        combined_img.paste(im, (x_offset, 0))
        x_offset += im.size[0]
        x_offset += letter_space
    combined_img.save('restored_crp_combined.png')

    to_guess = frames_count - len(rec) + rec.count('_')
    if to_guess > 3:
        raise ValueError(f"Need to brute to many values: {to_guess} for rec = {rec}")

    rec = '_' * (frames_count - len(rec)) + rec
    return brute_answer(rec, md5h, salt)


async def main():
    host, port = sys.argv[1].split(':')
    port = int(port)
    async with Channel(host, port) as channel:
        restore_cli = RestoreServiceStub(channel)

        async with restore_cli.Restore.open() as stream:
            # stream.
            await stream.send_request()

            for i in range(101):
                restore_req = await stream.recv_message()

                with open('crop.gif', 'wb') as f:
                    f.write(restore_req.gif)

                restored = restore(restore_req.gif, restore_req.hash, restore_req.salt)
                print(restored)

                await stream.send_message(Answer(answer=restored))
                print("Round = ", i)


if __name__ == '__main__':
    asyncio.run(main())
