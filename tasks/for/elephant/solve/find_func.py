#!/usr/bin/env python
import glob
import sys
import re


def find(b: bytes) -> list[tuple[str, int]]:
    res = []
    for fn in glob.glob("./dmp/*.dmp"):
        s_h, e_h = re.findall(r"pid\.\d+\.vma\.0x([0-9a-f]+)-0x([0-9a-f]+)\.dmp$", fn)[
            0
        ]
        s, e = int(s_h, 16), int(e_h, 16)
        with open(fn, "rb") as f:
            if (idx := f.read().find(b)) > -1:
                res.append((fn, s + idx))
    return res


if __name__ == "__main__":
    import sys

    if len(sys.argv) < 2:
        print("usage: find_func.py signature(int,text)", file=sys.stderr)

    sig = sys.argv[1]
    for _, offi in find(sig.encode()):
        print(_, offi, file=sys.stderr)
        for _, offj in find(offi.to_bytes(8, "little")):
            print(f"  // {sig}\n  0x{offj:x},")
