#!/usr/bin/env python3
import tqdm

from zlib import crc32


def find_crc(wordlist: list[str], hash: int) -> str:
    for w in tqdm.tqdm(wordlist):
        if crc32(w.encode()) == hash:
            return w
    raise Exception("password not found")


if __name__ == "__main__":
    import sys

    if len(sys.argv) != 3:
        print("usage: findpass.py wordlist-path password-hash", file=sys.stderr)
        sys.exit(1)

    with open(sys.argv[1], "r", encoding="latin-1") as f:
        words = list(map(str.strip, f))
        print(f"wordlist has {len(words)}")

    print(find_crc(words, int(sys.argv[2], 16)))
