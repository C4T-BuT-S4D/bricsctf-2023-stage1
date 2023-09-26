# ppc | gif0day

## Information

> SOS! Some crazy malware cropped all my gif collection!
> Hopefully, people say that this malware had some 0-day flaw.
>
> Could you please recover my gif collection ?
> `gif0day-f34fed283522eb7e.task.brics-ctf.ru:50051`

## Deploy

```sh
cd deploy
docker-compose -p ppc_gif0day up --build -d
```

## Public

Provide zip file from [public/](public/) and IP:PORT of the deployed app.

## TLDR

Recover original GIF frames from GIF cropped with acropalypse style vulnerability.
Recognize chars from recovered frames, brute missing ones.
Do this 100 times and get the flag.

## Writeup

We are given the gRPC service, that will 100 times do these steps:

1. Generate random string (will be the flag at step = 100).
2. Generate random salt.
3. Calculate md5 hash of this string and salt.
4. Create animated GIF image from the string, where 1 GIF frame = 1 character.
5. Crop resulting GIF to 30x30.
6. Send the cropped GIF bytes, hash and salt to the client.
7. Wait for client reply with timeout=10s.
8. Compare the client result with the string.
9. Repeated step 1 until step > 100, or client wasn't able to answer.

From the source code we can see that cropping is done with the acropalypse-style bug: the file is not truncated after
the crop, so the part of the original (not cropped) GIF file is still there.

As the cropped image is ~50 times smaller than original, most of the GIF frames can be safely recovered.
To recover the gif, we can follow this steps:

1. Read the cropped GIF from file to find where the 'trailer' (original GIF bytes) part begins.
2. Read the trailer, find fully-recoverable frames by reading the frame information (extension block).
3. Find the original image size, by reading each recovered frame size.
4. Write the cropped GIF header with modified size.
5. Write recovered frames.

This will give us the recovered GIF image. This image will have most of the frames (characters) recovered, but not all
of them, so we need to brute-force the rest (knowing the hash, salt).

The solution algorithm:

1. Get recovery request.
2. Recover the GIF.
3. Split the GIF into multiple frames.
4. For each frame do:
    1. Remove image noise.
    2. Run OCR.
5. Combine OCR results to get the answer suffix.
6. Brute force the prefix and poorly recognized chars using the fact we know the answer salted hash and the salt.
7. Repeat step '1' 100 times (last recovery request will have a flag).
8. Guess the flag prefix — it's brics+{ (flag format).

![Recovered flag GIF](solve/restored.gif)
[Exploit](solve/solve.py)

## Domain

`gif0day-f34fed283522eb7e.task.brics-ctf.ru`

## Flag

`brics+{gfQF2aN_g1f_Cr0P_brut3_xgRubfE}`

## Cloudflare
No
