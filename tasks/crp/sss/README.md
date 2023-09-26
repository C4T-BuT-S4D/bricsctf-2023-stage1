# crp | SSS

## Information

> Finally an implementation of Shamir's secret sharing that remains secure even if K shares are leaked

## Public

Provide 2 files: public/3.py and public/output.txt.xz

## TLDR

Interpolate the 3 missing points using the Lagrange polynomial, then use the inverse discrete Fourier transform to recover the polynomial's coefficients.

## Writeup

The challenge implements Shamir's secret sharing with 20000159 parties and polynomial degree 20000156. The entire polynomial's hash is used as the secret. The given Python program works in quadratic time and will probably take several months to complete; a trivial exploit would probably be even slower.

Fortunately for the challenge developer, there is a fast algorithm for evaluating a polynomial at a large set of points: the discrete Fourier transform. It takes as input a polynomial of degree N and an Nth root of unity, and evaluates the polynomial at all N powers of the root. We're working modulo a prime, so the multiplicative group is cyclic and there's a root of unity such that its powers include *all* field elements except 0. This means that a single DFT can compute all shares except for the 0th one (which can be computed directly). (of course, the DFT output has to be reordered)

The most common FFT algorithms do not work here, as the FFT size is not a smooth number (and MOD-1 is not smooth either). But it's possible to compute the DFT indirectly as a convolution using Bluestein's algorithm. The convolution can be computed with a conventional FFT algorithm.

To solve the challenge, the inverse DFT can be used (it interpolates the polynomial given its values at the powers of a root of unity). To apply it, it is necessary to recover the polynomial's values at 69, 420 and 1337. This can be done using the Lagrange polynomial. Applying the formula directly would require quadratic time, but it can be simplified by noting that the numerator and denominator include the product of almost all field elements, and instead computing the elements that are *not* included.

The generator and exploit are in `dev/src/main.rs`. My implementation does convolution using power-of-two iterative Cooley-Tukey FFT; to handle overflows two 64-bit moduli with CRT are used. Running with argument `gen` generates output.txt, `solve` reads output.txt from stdin and finds the flag.

## Flag

`brics+{94b8da67f4b9496276bbf1ad41650a46}`
