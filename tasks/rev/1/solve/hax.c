#include <kissat.h>
#include <stdint.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <stdio.h>
typedef void(*func)(void*);
void inject_jmp(uint8_t* from, uint8_t* to)
{
	uint32_t d = to - from - 5;
	*from++ = 0xe9; // jmp rel imm32
	memcpy(from, &d, 4);
}
kissat** kissat_ctx_loc = (kissat**)0x13172a0;
int write_handle(void* a, void* b, void* c)
{
	printf("Write hooked %p %p %p\n", a, b, c);
	//for(size_t i = 1; i < 10000; i++)
	//	printf("%d %d\n", i, kissat_value(*kissat_ctx_loc, i) > 0 ? 1 : -1);
	uint8_t flag[32] = {};
	for(size_t i = 0; i < 256; i++)
	{
		int x = kissat_value(*kissat_ctx_loc, i + 7641) > 0 ? 1 : 0; // obviously the flag is in variables 7641 to 7897
		flag[i/8] |= x << (i % 8);
	}
	for(size_t i = 0; i < 32; i++)
		putchar(flag[i]);
	fflush(stdout);
	exit(0);
	return -1;
}
int main()
{
	size_t skip = 65536;
	int fd = open("dump", 0, O_RDONLY); // ./dump is produced by "s 0; wtf dump 33554432" in radare2
	void* ret = mmap(skip, 33554432 - skip, PROT_READ|PROT_WRITE|PROT_EXEC, MAP_PRIVATE | MAP_FIXED, fd, skip);
	printf("%p\n", ret);
	func f = (func)0x20ee50;
	uint64_t x = 0;
	// skip adding conditions on the input, so that the solver instead computes the input given the output
	inject_jmp((void*)0x116c3f8, (void*)0x117bb18);
	inject_jmp((void*)0x117bb2e, (void*)0x118b3dc);
	inject_jmp((void*)0x118b3f2, (void*)0x119ab86);
	//inject_jmp((void*)0x119ab9c, (void*)0x11aa3a6); // segfaults due to missing an instruction that is on the normal path
	inject_jmp((void*)0x119ab9c, (void*)0x119abc8);
	inject_jmp((void*)0x119abd0, (void*)0x11aa3a6);
	// hook the output function (called when solving is complete) to print the solver's state instead
	inject_jmp((void*)0x1314e90, (void*)write_handle);
	f(&x);
}
