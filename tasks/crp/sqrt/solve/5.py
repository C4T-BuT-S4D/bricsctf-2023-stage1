from sage.all import *
import hashlib
import multiprocessing
import itertools
import ast
import sys
S = Permutations(256)(ast.literal_eval(input()))
eflag = bytes(ast.literal_eval(input()))
def check(A):
    dig = hashlib.sha512(str(A).encode()).digest()
    flag = bytes([x^y for x,y in zip(dig, eflag)])
    if flag.startswith(b'brics'):
        print(flag, file=sys.stderr)
        raise ZeroDivisionError # to hopefully exit from pool.map
def sqrt_odd(P):
    return [P[pow(2,-1,len(P))*i%len(P)] for i in range(len(P))]
def get_ways_even(group):
    if len(group) == 0:
        return [()]
    assert len(group) % 2 == 0
    ways = []
    for perm in itertools.permutations(group, r=len(group)):
        ok = True
        for x,y in [perm[i:i+2] for i in range(0,len(perm),2)]:
            if x > y:
                ok = False
                break
        if not ok:
            continue
        if any(perm[i] > perm[i+2] for i in range(0, len(perm)-2, 2)):
            continue
        prod_elts = []
        for x,y in [perm[i:i+2] for i in range(0,len(perm),2)]:
            joins = []
            for rot in range(len(y)):
                ry = y[rot:] + y[:rot]
                p = [a for pair in zip(x, ry) for a in pair]
                joins.append(p)
            prod_elts.append(joins)
        ways += list(itertools.product(*prod_elts))
    return ways
cs = S.to_cycles()
cs = sorted(cs, key=len)
group_sols = []
for sz, (*group,) in itertools.groupby(cs, key=len):
    if sz % 2 == 0:
        res = get_ways_even(group)
    else:
        res = []
        for nskip in range(len(group)%2, len(group)+1, 2):
            for skip in itertools.combinations(range(len(group)), nskip):
                take = [i for i in range(len(group)) if i not in skip]
                sqrts = [sqrt_odd(group[i]) for i in skip]
                for w in get_ways_even([group[i] for i in take]):
                    res.append(w + tuple(sqrts))
    group_sols.append(res)
print([len(x) for x in group_sols])
def try_1(sol):
    sol = sum(sol, ())
    assert sum(len(x) for x in sol) == 256
    A = sage.combinat.permutation.from_cycles(256, sol)
    assert A**2 == S
    check(A)
with multiprocessing.Pool(16) as pool:
    try:
        pool.map(try_1, itertools.product(*group_sols))
    except ZeroDivisionError:
        exit()
