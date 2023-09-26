## crp | arc3

## Information

> \> RC3 was broken before ever being used.

## Public

Provide 2 files: public/2.py and public/data.bin.

## Writeup

The script encrypts a string with a known 1000-byte prefix with RC4 200000 times. While RC4 is quite weak, as far as I know it isn't weak enough to be breakable in this scenario. The (intended) vulnerability is the unusual initialization algorithm: there is a 1/65536 probability for it to result in a weak state, known as a Finney cycle. For reference, here is the RC4 output generation algorithm:

```
while True:
    i = (i + 1) % 256
    j = (j + S[i]) % 256
    S[i], S[j] = S[j], S[i]
    yield S[(S[i] + S[j]) % 256]
```

If at the start of an iteration `S[i + 1] == 1` and `j == i + 1`, then at the end the same condition will repeat with `i` and `j` increased by 1 and `S[i+1]` moved cyclically by 1. After 256 such iterations, `i`, `j` and the position of `1` in `S` will be the same, but all other elements of `S` will be rotated by 1 relative to `S[j]`. This means that the period of a weak state is 256 * 255. The transition function is invertible, therefore this type of weak states can never occur with correct initialization (which sets `i = 0`)

Only 1000 bytes of output are given, which is significantly less than 65280, but it's still sufficient to recover the state. It's likely that there are multiple solutions here, I describe the one I found. Knowing `S[i + 1] == 1` and `j == i + 1` simplifies the output byte formula to `S[S[i+2]+1]`. Observe that after 255 iterations, the entire array is cyclically rotated by 1, meaning that the value of `S[i+2]+1` does not change and that output bytes `n`, `n+255`, `n+510`, ... correspond to the same position in `S`. So if we consider `[output[i::255] for i in range(255)]`, we have 255 taps into the rotating internal state (i.e. we're given some cyclic substrings of the initial state of length 3 or 4, chosen at random positions).

The reality is slightly more complicated: while computing `S[i+2]+1` at the start of an iteration always correctly predicts the array position that will be read, the read happens after the swap. This sometimes causes an `1` to be inserted in the output. I wasn't able to understand this part fully, but this is sufficiently close to the end to switch from thinking to guessing.

We can discard all `1`s from every tap's output. Then we have to reconstruct an array of elements in [0,...,254] given some (cyclic) substrings of it. In general, this is a NP-complete problem (though easy to approximate well enough), but a very simple solution  exists in this case: every character occurs exactly once in the solution, so we can compute all bigrams, walk the resulting cycle and reconstruct the initial state up to a cyclic rotation and the position of the `1`. The exploit is in `solve/expl.py`.

## Flag

`brics+{551b6a3f4bfd7ff0}`
