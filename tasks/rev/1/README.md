# rev | 1

## Information

[none]

## Public

Provide binary `public/1`.

## Writeup

When executed, the binary reads a string from stdin, checks whether it equals the flag and prints "Correct" or "Incorrect" accordingly. The algorithm relies on a simple ad hoc pseudorandom permutation on 32-byte strings, based on some additions, rotations and XORs on 4 64-bit numbers. The check is represented as the constraint system `F(X) == input`, where `X` is a constant embedded in the program. To efficiently determine if this system is consistent, the Kissat SAT solver (https://github.com/arminbiere/kissat) is called. It's possible to solve the challenge by removing the constraint that `input` matches what has been read on stdin -- then the system will have a solution and the obtained value of `input` will equal the flag. The exploit is in file solve/hax.c.

## Flag

`brics+{ef2b18c79dc0c4fe17ba77aa}`
