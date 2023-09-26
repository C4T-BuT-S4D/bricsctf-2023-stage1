#include "sys/mman.h"
#include "sys/fcntl.h"

#include "pl_funcs.c"

const char segm_fname[] = "../dmp/pid.1002.vma.0x56361fc27000-0x56361fe72000.dmp";
const size_t segm_start = 0x56361fc27000;
const size_t segm_len = 0x56361fe72000 - segm_start;

PLpgSQL_function *funcs[] = {
  // aaa(bit varying)
  0x56361fc67380,
  // aab(bytea)
  0x56361fc66f78,
  // aba(bytea,bytea)
  0x56361fd62e50,
  // baa(bytea)
  0x56361fd62a48,
  // bab(bytea)
  0x56361fd62640,
  // bba(bytea)
  0x56361fd63258,
  // abb(bytea,bytea)
  0x56361fc67b90,
  // bbb(bytea)
  0x56361fc66b70,
  // lll(bytea,bytea)
  0x56361fc67788,
  // l1l(text,text)
  0x56361fc663b0,
  // secret(text,text)
  0x56361fc65b18,
};

int main() {
  int segm_fd = open(segm_fname, O_RDONLY);
  void *segm = mmap((void *) segm_start, segm_len, PROT_READ, MAP_PRIVATE | MAP_FIXED, segm_fd, 0); 

  for (size_t i = 0; i < sizeof(funcs) / sizeof(void *); i++) {
    plpgsql_dumptree(funcs[i]);
  }
}
