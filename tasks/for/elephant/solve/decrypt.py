#!/usr/bin/env python3
import warnings

from zlib import crc32 as _crc32
from Crypto.Cipher import AES


def crc32(b: bytes) -> bytes:
    return _crc32(b).to_bytes(4, "big")


def derive_key(name: str, password: str) -> bytes:
    key = crc32((name + password).encode() + b"p3pp3r")
    for _ in range(1, 4 + (key[2] % 3) * 2):
        sk = max(1, key[-2] % 4)
        key = key + crc32(key[sk - 1 :])
    return key


def decrypt(name: str, password: str, content: bytes) -> str:
    key = derive_key(name, password)

    aes = AES.new(key, AES.MODE_CTR, nonce=b"", initial_value=1)
    dec = aes.decrypt(content)

    if dec[:4] != crc32(dec[4:]):
        warnings.warn("crc check failed")

    return dec[4:].decode()


if __name__ == "__main__":
    import sys

    if len(sys.argv) != 4:
        print("usage: decrypt.py name password content-hex", file=sys.stderr)
        sys.exit(1)

    print(decrypt(sys.argv[1], sys.argv[2], bytes.fromhex(sys.argv[3])))
