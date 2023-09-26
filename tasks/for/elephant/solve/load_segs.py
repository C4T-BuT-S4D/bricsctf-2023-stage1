import re
import sys
import glob

for fn in glob.glob(f"./dmp/*.dmp"):
    print(fn)
    s_h, e_h = re.findall(r"pid\.\d+\.vma\.0x([0-9a-f]+)-0x([0-9a-f]+)\.dmp$", fn)[0]
    s, e = int(s_h, 16), int(e_h, 16)

    add_segm_ex(s, e, 0, 2, 8, 0, 0)
    LoadFile(fn, 0, s, e - s)
