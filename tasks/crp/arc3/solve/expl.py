def gen_ks(S, i, j):
    S = list(S) # copy
    while True:
        i = (i + 1) % 256
        j = (j + S[i]) % 256
        S[i], S[j] = S[j], S[i]
        yield S[(S[i] + S[j]) % 256]
def try_atk(keystream, enc_flag):
    # ignore all 1s, they appear to be inserted randomly?
    taps = [bytes([c for c in keystream[i::255] if c != 1]) for i in range(255)]
    # this is a shortest common superstring problem, but every char occurs exactly once, so we can solve by computing unique bigrams
    bigrams = set([taps[j][i:i+2] for j in range(255) for i in range(len(taps[j])-1)])
    if len(bigrams) != 255:
        return # not a Finney cycle
    print('FOUND')
    nextc = {}
    for a,b in bigrams:
        nextc[a] = b
    rs = b''
    i = 0
    for _ in range(255):
        rs += bytes([i])
        i = nextc[i]
    # rs is the initial state, but with the 1 removed and cyclically rotated by an unknown amount
    # we reach this point only for weak states so just trying all possibilities is fast enough
    for cshift in range(255):
        s1 = rs[cshift:] + rs[:cshift]
        for pos1 in range(256):
            s2 = s1[:pos1] + bytes([1]) + s1[pos1:]
            # we now check if s2 is the correct initial state. i,j can be computed using the Finney cycle condition
            j = pos1
            i = j - 1
            pt = bytes([a^b for a,b in zip(keystream + enc_flag, gen_ks(s2, i, j))])
            if pt[:1000] != b'\x00'*1000:
                continue
            print(cshift, pos1, pt[1000:])
import sys
with open(sys.argv[1],'rb') as f:
    for i in range(200000):
        dat = f.read(1024)
        try_atk(dat[:1000], dat[1000:])
