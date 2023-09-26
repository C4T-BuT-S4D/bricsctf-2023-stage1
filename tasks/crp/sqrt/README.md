# crp | sqrt

## Information

> In this challenge you have to compute the square root of a group element.

## Public

Provide 2 files: public/4.py and public/output.txt

## Writeup

Any permutation can be represented as a product of disjoint cycles. When it is squared, cycles of odd length get reordered reversibly, and cycles of even length are replaced by 2 cycles of half their length (e.g. `1->2->3->4->1` becomes `1->3->1` and `2->4->2`). This process is invertible as well: it is possible to combine any 2 cycles of equal length into a larger one that gives them as its square (though there's often more than 1 way to do so). Implementing all of this is surprisingly non-trivial, the exploit is in solve/5.py.

## Flag

`brics+{ab99943f6dae4f20595c8645fcf8289e}`
