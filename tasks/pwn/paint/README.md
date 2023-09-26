# pwn | paint

## Information

> nc paint-71ae86dc10a3fe17.task.brics-ctf.ru 13003

## Deploy

```
cd deploy/
docker-compose -p paint up --build -d
```

## Public

Provide zip file from [public/](./public/) and IP:PORT of the deployed app

## TLDR

Rewrite `l_addr` in `link_map` so that address of `free` will be written to `memset`. Using this free fake chunk on stack and write rop

## Writeup

There isn't any check of index in this program. This gives us several ways to leak libc, but all freed pointers are zeroed, so we can't get arbitrary write through heap. Instead of this we will rewrite two last bytes in `l_addr` in `link_map` and call `delete_canvas`. `free` will be called at first time, and because of corrupted `l_addr` pointer to `free` will be written to `memset`. Then on calling `rate`, `memset` will try to free buffer on stack. Before buffer there are variables `rate` and `index`. Set them to get valid header and get freed chunk on stack. Now we can leak some addresses and write rop.

[Exploit](./solve/sploit.py)

## Domain
`paint-71ae86dc10a3fe17.task.brics-ctf.ru`

## Flag

```
brics+{f024b6861b74986efc564bbbe89da0ab}
```

## Cloudflare
No
