#!/usr/bin/env python
from secrets import token_bytes

with open("./sources", "rb") as f:
    sources = f.read().split(b"---ooo---")

print(*sources, sep="\n")

with open("./ram1.img", "rb") as f:
    image = f.read()

for src in sources:
    print(src[16:32])
    image = image.replace(src, token_bytes(len(src)))

with open("./ram_clean.img", "wb") as f:
    f.write(image)
