#include <array>
#include <vector>
#include <random>
#include <bitset>
#include <m4ri/m4ri.h>
uint8_t func(uint64_t x)
{
	x *= 5;
	x = (x << 7) | (x >> (64 - 7));
	x *= 9;
	return x >> 56;
}
constexpr size_t N_MONO = 1 + 256 + 256 * 255 / 2;
using deg2_rel = std::bitset<1 + 256 + 256 * 255 / 2>;
using deg2_rel64 = std::bitset<1 + 64 + 64 * 63 / 2>;
using lin_rel = std::bitset<256>;
struct sym_u64
{
	lin_rel st[64] {};
	sym_u64& operator^=(const sym_u64& rhs)
	{
		for(size_t i = 0; i < 64; i++)
			st[i] ^= rhs.st[i];
		return *this;
	}
	sym_u64 operator^(const sym_u64& rhs) const
	{
		sym_u64 ret = *this;
		ret ^= rhs;
		return ret;
	}
	sym_u64 operator<<(int by) const
	{
		sym_u64 ret {};
		for(size_t i = 0; i < 64 - by; i++)
			ret.st[i + by] = st[i];
		return ret;
	}
	sym_u64 operator>>(int by) const
	{
		sym_u64 ret {};
		for(size_t i = 0; i < 64 - by; i++)
			ret.st[i] = st[i + by];
		return ret;
	}
	sym_u64 rotl(int by) const
	{
		sym_u64 ret {};
		for(size_t i = 0; i < 64; i++)
			ret.st[(i + by)%64] = st[i];
		return ret;
	}
};
struct sym_xs256
{
	sym_u64 s0, s1, s2, s3;
	sym_xs256()
	{
		for(size_t i = 0; i < 64; i++)
		{
			s0.st[i][i] = true;
			s1.st[i][i+64] = true;
			s2.st[i][i+128] = true;
			s3.st[i][i+192] = true;
		}
	}
	sym_u64 step()
	{
		sym_u64 res_s1 = s1;
		sym_u64 t = s1 << 17;
		s2 ^= s0;
		s3 ^= s1;
		s1 ^= s2;
		s0 ^= s3;
		s2 ^= t;
		s3 = s3.rotl(45);
		return res_s1;
	}
};
std::array<std::vector<deg2_rel64>, 256> krels;
deg2_rel to_deg2(lin_rel lhs)
{
	deg2_rel ret {};
	for(size_t i = 0; i < 256; i++)
		ret[i+1] = lhs[i]; // 0 is the constant term
	return ret;
}
deg2_rel mul(lin_rel lhs, lin_rel rhs)
{
	deg2_rel ret {}; // compute the product coefficients directly
	size_t mi = 1;
	for(size_t i = 0; i < 256; i++)
	{
		ret[mi] = ret[mi] ^ (lhs[i] & rhs[i]);
		mi++;
	}
	for(size_t i = 0; i < 256; i++)
	for(size_t j = 0; j < i; j++)
	{
		ret[mi] = ret[mi] ^ (lhs[i] & rhs[j]);
		ret[mi] = ret[mi] ^ (lhs[j] & rhs[i]);
		mi++;
	}
	return ret;
}
int main()
{
	FILE* f = fopen("list_rels", "r");
	while(true)
	{
		char buf[4096];
		int res = 0;
		int ok = fscanf(f, "%d%s", &res, buf);
		if(ok < 2)
			break;
		deg2_rel64 rel {};
		for(size_t i = 0; i < rel.size(); i++)
			rel[i] = buf[i] - '0';
		krels[res].push_back(std::move(rel));
	}
	fclose(f);
	sym_xs256 rng {};
	// we're targeting the state before the flag is encrypted so skip some bytes now
	f = fopen("out_bin", "rb");
	constexpr int FLAG_LEN = 41;
	uint8_t encflag[FLAG_LEN];
	for(size_t i = 0; i < 41; i++)
	{
		encflag[i] = fgetc(f);
		rng.step();
	}
	std::vector<deg2_rel> all_rels;
	for(size_t i = 0; i < 2000; i++)
	{
		fprintf(stderr, "i=%d\n", i);
		uint8_t outb = fgetc(f);
		sym_u64 s1 = rng.step();
		for(const deg2_rel64& big_rel : krels[outb])
		{
			deg2_rel res {};
			size_t mi = 1;
			if(big_rel[0])
				res[0] = res[0] ^ 1;
			for(size_t i = 0; i < 64; i++)
			{
				if(big_rel[mi])
					res ^= to_deg2(s1.st[i]);
				mi++;
			}
			for(size_t i = 0; i < 64; i++)
			for(size_t j = 0; j < i; j++)
			{
				if(big_rel[mi])
					res ^= mul(s1.st[i], s1.st[j]);
				mi++;
			}
			all_rels.push_back(res);
		}
	}
	fclose(f);
	fprintf(stderr, "%zu\n", all_rels.size());
	mzd_t *A = mzd_init(all_rels.size(), N_MONO);
	for(size_t i = 0; i < all_rels.size(); i++)
	{
		for(size_t j = 0; j < N_MONO; j++)
			mzd_write_bit(A, i, j, all_rels[i][j]);
	}
	// we solve this system by linearization
	fprintf(stderr, "starting solve\n");
	fprintf(stderr, "mat dim: %zu rows, %zu cols\n", A->nrows, A->ncols);
	mzd_t* ker = mzd_kernel_left_pluq(A, 0);
	mzd_t* tker = mzd_transpose(nullptr, ker);
	fprintf(stderr, "tker dim: %zu rows, %zu cols\n", tker->nrows, tker->ncols);
	mzd_free(A);
	mzd_free(ker);
	assert(tker->nrows == 1); // only 1 solution
	sym_xs256 rng2 {};
	lin_rel ist; // the initial state. dot product with symbolic output equals concrete output
	for(size_t j = 0; j < 256; j++)
		ist[j] = mzd_read_bit(tker, 0, j + 1);
	for(size_t i = 0; i < FLAG_LEN; i++)
	{
		sym_u64 s1 = rng2.step();
		uint64_t x = 0;
		for(size_t j = 0; j < 64; j++)
			x |= uint64_t((s1.st[j] & ist).count() % 2) << j;
		uint8_t outb = func(x);
		printf("%c", encflag[i] ^ outb);
	}
	printf("\n");
}
