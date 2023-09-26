import random
import base64
from Crypto.Cipher import AES

part = "perlovkaOchenVkusnAya"
flag = "brics+{" + part + "}"

xored = 0
pos = 0
for i, c in enumerate(part):
    r = random.randint(1, 40960)
    prog = f'''my ${c}kek = 1; if ($%skek == 1) ''' + '{ return %d; } else { return 0; } //''' % (r)
    to_fill = len(prog) % 16
    prog += ' ' * (16 - to_fill)
    key = part[i // 2] * 16
    # key = c * 16
    
    cipher = AES.new(key.encode(), AES.MODE_ECB)
    encrypted = cipher.encrypt(prog.encode())
    # encrypted = encrypted.hex()
    encrypted = base64.b64encode(encrypted).decode()
    print("'" + encrypted + "',")
    xored ^= r
print(xored)




