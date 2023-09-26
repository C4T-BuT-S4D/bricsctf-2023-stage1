#include <vector>
#include <random>
#include <m4ri/m4ri.h>
// we need to find degree 2 relations that hold certainly
uint8_t f(uint64_t x)
{
	x *= 5;
	x = (x << 7) | (x >> (64 - 7));
	x *= 9;
	return x >> 56;
}
void process(uint8_t tgt)
{
	std::mt19937_64 mt;
	size_t n_cols = 1;
	n_cols += 64;
	n_cols += 64 * 63 / 2;
	fprintf(stderr, "%zu monomials\n", n_cols);
	mzd_t *A = mzd_init(20000, n_cols);
	for(size_t ri = 0; ri < A->nrows; ri++)
	{
		uint64_t x = 0;
		do
			x = mt();
		while(f(x) != tgt);
		//printf("%zu\n", x);
		size_t mi = 0;
		mzd_write_bit(A, ri, mi++, 1);
		for(size_t i = 0; i < 64; i++)
			mzd_write_bit(A, ri, mi++, (x >> i) & 1);
		for(size_t i = 0; i < 64; i++)
		for(size_t j = 0; j < i; j++)
			mzd_write_bit(A, ri, mi++, (x >> i) & (x >> j) & 1);
		assert(mi == n_cols);
	}
	mzd_t* ker = mzd_kernel_left_pluq(A, 0);
	fprintf(stderr, "mat dim: %zu rows, %zu cols\n", A->nrows, A->ncols);
	fprintf(stderr, "ker dim: %zu rows, %zu cols\n", ker->nrows, ker->ncols);
	mzd_t* tker = mzd_transpose(nullptr, ker);
	fprintf(stderr, "tker dim: %zu rows, %zu cols\n", tker->nrows, tker->ncols);
	mzd_free(A);
	mzd_free(ker);
	/*
	// check again that the relations are certain (not very fast)
	for(size_t ri = 0; ri < tker->nrows; ri++)
	{
		for(size_t _ = 0; _ < 100000; _++)
		{
			uint64_t x = 0;
			do
				x = mt();
			while(f(x) != tgt);
			//printf("%zu\n", x);
			size_t nm = 0;
			size_t mi = 0;
			nm ^= mzd_read_bit(tker, ri, mi++);
			for(size_t i = 0; i < 64; i++)
				nm ^= mzd_read_bit(tker, ri, mi++) & (x >> i) & 1;
			for(size_t i = 0; i < 64; i++)
			for(size_t j = 0; j < i; j++)
				nm ^= mzd_read_bit(tker, ri, mi++) & (x >> i) & (x >> j) & 1;
			//assert(mi == n_cols);
			assert(nm == 0);
		}
	}
	*/
	for(size_t ri = 0; ri < tker->nrows; ri++)
	{
		printf("%d ", (int)tgt);
		for(size_t ci = 0; ci < tker->ncols; ci++)
			printf("%d", (int)mzd_read_bit(tker, ri, ci));
		printf("\n");
	}
	mzd_free(tker);
}
int main()
{
	for(size_t i = 0; i < 256; i++)
		process(i);
}
