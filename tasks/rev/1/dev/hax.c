#include "src/kissat.h"
#include "stddef.h"
#include "stdlib.h"
#include "stdint.h"
#include "string.h"
#include "stdio.h"
#include "assert.h"
#include "sys/random.h"
struct num64
{
	int64_t vars[64];
};
static int64_t cur_var = 1;
static kissat* sol_ctx;
int64_t get_var()
{
	return cur_var++;
}
//struct llist_nums
//{
//	int64_t x;
//	struct llist_nums* next;
//};
//struct llist_nums* big_list_head = NULL;
//struct llist_nums* big_list_tail = NULL;
//void add_num(int64_t x)
//{
//	if(!big_list_tail)
//	{
//		big_list_head = malloc(sizeof(struct llist_nums));
//		big_list_tail = big_list_head;
//		big_list_head->x = x;
//		big_list_head->next = NULL;
//	}
//	else
//	{
//		struct llist_nums* new_node = malloc(sizeof(struct llist_nums));
//		new_node->x = x;
//		new_node->next = NULL;
//		big_list_tail->next = new_node;
//		big_list_tail = new_node;
//	}
//}
#define ADD(x) do kissat_add(sol_ctx, x); while(0)
//#define ADD(x) do add_num(x); while(0)
#define ADD2(a, b) do { ADD(a); ADD(b); } while(0)
#define ADD3(a, b, c) do { ADD(a); ADD(b); ADD(c); } while(0)
#define ADD4(a, b, c, d) do { ADD(a); ADD(b); ADD(c); ADD(d); } while(0)
#define ADD5(a, b, c, d, e) do { ADD(a); ADD(b); ADD(c); ADD(d); ADD(e); } while(0)
struct num64* add(struct num64* lhs, struct num64* rhs)
{
	struct num64* ret = malloc(sizeof(struct num64));
	int64_t cout = get_var();
	int64_t r0 = get_var();
	ADD4(-lhs->vars[0], -rhs->vars[0], -r0, 0);
	ADD4(-lhs->vars[0], cout, r0, 0);
	ADD4(lhs->vars[0], -rhs->vars[0], r0, 0);
	ADD4(lhs->vars[0], rhs->vars[0], -r0, 0);
	ADD3(rhs->vars[0], -cout, 0);
	ADD3(-cout, -r0, 0);
	ret->vars[0] = r0;
	int64_t cin = cout;
	for(size_t i = 1; i < 63; i++)
	{
		cout = get_var();
		int64_t r = get_var();
		int64_t a = lhs->vars[i];
		int64_t b = rhs->vars[i];
		ADD5(-a, -b, -cin, r, 0);
		ADD5(-a, -b, cin, -r, 0);
		ADD5(-a, b, -cin, -r, 0);
		ADD4(-a, cout, r, 0);
		ADD5(a, -b, cin, r, 0);
		ADD5(a, b, -cin, r, 0);
		ADD5(a, b, cin, -r, 0);
		ADD4(a, -cout, -r, 0);
		ADD4(-b, -cin, cout, 0);
		ADD4(b, cin, -cout, 0);
		ret->vars[i] = r;
		cin = cout;
	}
	int64_t a = lhs->vars[63];
	int64_t b = rhs->vars[63];
	int64_t r = get_var();
	ADD5(-a, -b, -cin, r, 0);
	ADD5(-a, -b, cin, -r, 0);
	ADD5(-a, b, -cin, -r, 0);
	ADD5(-a, b, cin, r, 0);
	ADD5(a, -b, -cin, -r, 0);
	ADD5(a, -b, cin, r, 0);
	ADD5(a, b, -cin, r, 0);
	ADD5(a, b, cin, -r, 0);
	ret->vars[63] = r;
	return ret;
}
struct num64* xor(struct num64* lhs, struct num64* rhs)
{
	struct num64* ret = malloc(sizeof(struct num64));
	for(size_t i = 0; i < 64; i++)
	{
		int64_t r = get_var();
		int64_t a = lhs->vars[i];
		int64_t b = rhs->vars[i];
		ret->vars[i] = r;
		ADD4(-a, -b, -r, 0);
		ADD4(-a, b, r, 0);
		ADD4(a, -b, r, 0);
		ADD4(a, b, -r, 0);
	}
	return ret;
}
struct num64* rotl(struct num64* lhs, uint64_t x)
{
	struct num64* ret = malloc(sizeof(struct num64));
	for(size_t i = 0; i < 64 - x; i++)
		ret->vars[i+x] = lhs->vars[i];
	for(size_t i = 0; i < x; i++)
		ret->vars[i] = lhs->vars[i+64-x];
	return ret;
}
struct num64* shl(struct num64* lhs, uint64_t x)
{
	struct num64* ret = malloc(sizeof(struct num64));
	for(size_t i = 0; i < 64 - x; i++)
		ret->vars[i+x] = lhs->vars[i];
	for(size_t i = 0; i < x; i++)
	{
		ret->vars[i] = get_var();
		ADD2(-ret->vars[i], 0);
	}
	return ret;
}
void assert_const(struct num64* lhs, uint64_t val)
{
	for(size_t i = 0; i < 64; i++)
	{
		int64_t r = get_var();
		if(val >> i & 1)
			ADD2(lhs->vars[i], 0);
		else
			ADD2(-lhs->vars[i], 0);
	}
}
struct num64* get_num()
{
	struct num64* ret = malloc(sizeof(struct num64));
	for(size_t i = 0; i < 64; i++)
		ret->vars[i] = get_var();
	return ret;
}
uint64_t get_sol(struct num64* x)
{
	uint64_t ret = 0;
	for(size_t i = 0; i < 64; i++)
		ret |= (uint64_t)(kissat_value(sol_ctx, x->vars[i]) > 0 ? 1 : 0) << i;
	return ret;
}
int main()
{
	sol_ctx = kissat_init();
	struct num64* a0 = get_num(), *b0 = get_num(), *c0 = get_num(), *d0 = get_num();
	struct num64* a = a0, *b = b0, *c = c0, *d = d0;
	for(size_t i = 0; i < 10; i++)
	{
		b = add(b, a);
		c = add(c, b);
		d = add(d, c);
		a = add(a, d);
		a = rotl(a, 32);
		b = rotl(b, 7);
		c = rotl(c, 15);
		d = rotl(d, 24);
		a = xor(a, c);
		b = xor(b, d);
		c = xor(c, b);
		d = xor(d, a);
	}
	char string[40] = {};
	fgets(&string, 40, stdin);
	uint64_t at; memcpy(&at, string, 8);
	uint64_t bt; memcpy(&bt, string+8, 8);
	uint64_t ct; memcpy(&ct, string+16, 8);
	uint64_t dt; memcpy(&dt, string+24, 8);
	assert_const(a, at);
	assert_const(b, bt);
	assert_const(c, ct);
	assert_const(d, dt);
	assert_const(a0, 0x591fd26dcdd4ebf2);
	assert_const(b0, 0x0509dc63bd8415ed);
	assert_const(c0, 0x8de748e0a2dde3f5);
	assert_const(d0, 0x92c79a9c875694f8);
	// the commented code below was actually present during the final build. Unfortunately, `big_list_head` is always null in the final version, so this block was patched out manually
	//struct llist_nums* cur = big_list_head;
	//do
	//{
	//	kissat_add(sol_ctx, cur->x);
	//	cur = cur->next;
	//}
	//while(cur != nullptr);
	int x = kissat_solve(sol_ctx);
	puts(x == 10 ? "Correct" : "Incorrect");
	kissat_release(sol_ctx);
}
