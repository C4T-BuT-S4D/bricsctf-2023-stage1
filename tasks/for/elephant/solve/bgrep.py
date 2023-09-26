#!/usr/bin/env python
import glob
import sys
import re

sb = bytes.fromhex(sys.argv[1])

for fn in glob.glob("./dmp/*.dmp"):
    s_h, e_h = re.findall(r"pid\.\d+\.vma\.0x([0-9a-f]+)-0x([0-9a-f]+)\.dmp$", fn)[0]
    s, e = int(s_h, 16), int(e_h, 16)

    with open(fn, "rb") as f:
        if (idx := f.read().find(sb)) > -1:
            print(f"{fn}\t0x{s+idx:016x}")
