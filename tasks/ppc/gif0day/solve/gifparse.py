# Patched copy of https://github.com/qalle2/pygif/blob/main/gifstruct.py

import os, struct, sys

# for Graphic Control Extension
DISPOSAL_METHODS = {
    0: "unspecified",
    1: "leave in place",
    2: "restore to background color",
    3: "restore to previous",
}


def error(descr):
    # print error and exit
    raise ValueError(f"Error: {descr}")


def printval(descr, value):
    # print a value with description
    assert isinstance(value, int) or isinstance(value, str) \
           or isinstance(value, bytes)
    if isinstance(value, bool):
        value = ("no", "yes")[value]
    elif isinstance(value, bytes):
        value = value.decode("ascii", errors="backslashreplace")
    # print(4 * " " + f"{descr}: {value}")


def getoffs(handle, adjust=0):
    # print file offset
    assert isinstance(adjust, int)
    return handle.tell() + adjust
    # printval("file offset", handle.tell() + adjust)


def getbytes(handle, length):
    # read bytes from file
    assert isinstance(length, int)
    data = handle.read(length)
    if len(data) < length:
        error("unexpected end of file")
    return data


def get_subblocks(handle):
    # generate data from GIF subblocks
    sbSize = getbytes(handle, 1)[0]  # subblock size
    while sbSize:
        chunk = getbytes(handle, sbSize + 1)  # subblock & size of next one
        yield chunk[:-1]
        sbSize = chunk[-1]


# -----------------------------------------------------------------------------

def read_header(handle):
    # read Header from current file position; return file version

    (id_, version) = struct.unpack("3s3s", getbytes(handle, 6))
    if id_ != b"GIF":
        error("not a GIF file")
    return version


def read_lsd(handle):
    # read Logical Screen Descriptor from current file position; return a dict

    (width, height, packedFields, bgIndex, aspectRatio) \
        = struct.unpack("<2H3B", getbytes(handle, 7))

    return {
        "width": width,
        "height": height,
        "gctFlag": bool(packedFields & 0b10000000),
        "colorResolution": ((packedFields >> 4) & 0b00000111) + 1,
        "sortFlag": bool(packedFields & 0b00001000),
        "gctSize": (packedFields & 0b00000111) + 1,
        "bgIndex": bgIndex,
        "aspectRatio": aspectRatio,
    }


def read_image(handle):
    # read information of one image in GIF file
    # handle position must be at first byte after ',' of Image Descriptor
    # return a dict

    (x, y, width, height, packedFields) \
        = struct.unpack("<4HB", getbytes(handle, 9))

    return {
        "x": x,
        "y": y,
        "width": width,
        "height": height,
        "lctFlag": bool(packedFields & 0b10000000),
        "interlaceFlag": bool(packedFields & 0b01000000),
        "sortFlag": bool(packedFields & 0b00100000),
        "lctSize": (packedFields & 0b00000111) + 1,
    }


def read_extension_block(handle):
    # read Extension block in GIF file;
    # handle position must be at first byte after Extension Introducer ('!')

    label = getbytes(handle, 1)[0]
    if label == 0x01:
        # TODO (low priority): print more info
        printval("type", "Plain Text")
        getbytes(handle, 13)  # skip bytes
        list(get_subblocks(handle))  # skip subblocks
    elif label == 0xf9:
        printval("type", "Graphic Control")
        (packedFields, delayTime, transparentIndex) \
            = struct.unpack("<xBHBx", getbytes(handle, 6))
        disposal = (packedFields & 0b00011100) >> 2
        userInput = bool(packedFields & 0b00000010)
        transparentFlag = bool(packedFields & 0b00000001)
        printval(
            "delay time in 1/100ths of a second",
            (delayTime if delayTime else "none")
        )
        printval("wait for user input", userInput)
        printval(
            "transparent color index",
            (transparentIndex if transparentFlag else "none")
        )
        printval("disposal method", DISPOSAL_METHODS.get(disposal, "?"))
    elif label == 0xfe:
        printval("type", "Comment")
        data = b"".join(get_subblocks(handle))
        printval("data", data)
    elif label == 0xff:
        # TODO (low priority): print more info
        printval("type", "Application")
        (identifier, authCode) = struct.unpack("x8s3s", getbytes(handle, 12))
        printval("identifier", identifier)
        printval("authentication code", authCode)
        list(get_subblocks(handle))  # skip subblocks
    else:
        error("unknown extension type")


def read_file(handle):
    handle.seek(0)

    # print("Header:")
    getoffs(handle)
    version = read_header(handle)
    printval("version", version)

    # print("Logical Screen Descriptor:")
    getoffs(handle)
    lsdInfo = read_lsd(handle)
    printval("width", lsdInfo["width"])
    printval("height", lsdInfo["height"])
    printval(
        "original color resolution in bits per RGB channel",
        lsdInfo["colorResolution"]
    )
    printval(
        "pixel aspect ratio in 1/64ths",
        (lsdInfo["aspectRatio"] + 15) if lsdInfo["aspectRatio"] else "unknown"
    )
    printval("has Global Color Table", lsdInfo["gctFlag"])

    if lsdInfo["gctFlag"]:
        # print("Global Color Table:")
        getoffs(handle)
        printval("colors", 2 ** lsdInfo["gctSize"])
        printval("sorted", lsdInfo["sortFlag"])
        printval("background color index", lsdInfo["bgIndex"])
        getbytes(handle, 2 ** lsdInfo["gctSize"] * 3)  # skip it

    # read rest of blocks
    while True:
        blockType = getbytes(handle, 1)
        if blockType == b",":
            # print("Image Descriptor:")
            getoffs(handle, -1)
            imageInfo = read_image(handle)
            printval("x position", imageInfo["x"])
            printval("y position", imageInfo["y"])
            printval("width", imageInfo["width"])
            printval("height", imageInfo["height"])
            printval("interlaced", imageInfo["interlaceFlag"])
            printval("has Local Color Table", imageInfo["lctFlag"])

            if imageInfo["lctFlag"]:
                # print(f"Local Color Table:")
                getoffs(handle)
                printval("colors", 2 ** imageInfo["lctSize"])
                printval("sorted", imageInfo["sortFlag"])
                getbytes(handle, 2 ** imageInfo["lctSize"] * 3)  # skip it

            # print("LZW data:")
            # TODO (low priority): print more info
            getoffs(handle)
            lzwPalBits = getbytes(handle, 1)[0]
            lzwDataLen = sum(len(d) for d in get_subblocks(handle))
            printval("palette bit depth", lzwPalBits)
            printval("data size", lzwDataLen)
        elif blockType == b"!":
            # print("Extension:")
            getoffs(handle, -1)
            read_extension_block(handle)
        elif blockType == b";":
            # print("Trailer:")
            return getoffs(handle, -1)
            # break
        else:
            error("unknown block type")
