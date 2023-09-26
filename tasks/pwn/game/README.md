# pwn | game

## Information

Playing the Fool.
> nc game-51463dfee2ff6388.task.brics-ctf.ru 13001

## Deploy

```
cd deploy/
docker-compose -p game up --build -d
```

## Public

Provide zip file from [public/](./public/) and IP:PORT of the deployed app

## TLDR

Rewrite `user->name_sz` field, leak heap, tcache poisoning to leak `stdin` from binary, rewrite `strlen` in libc `.got` to `system`.

## Writeup

When we add portal in `add_maze()` function, there's no check that the portal will be not on the side of the maze. So we can create the maze like this:
```
##############A#
#@A            #
######## #######
```
and if we go right and up we will escape from the maze. Now we can go and rewrite `user->name_sz` to the bigger one.

To get libc leak, we will do tcache poisoning on struct `user`. Function `cut` in `play_cards` calls `split`, which several times calls `malloc`. If we construct the following tree
```
     node1 (root)
(l) /     \ (r)
  0      node2
    (l) /     \ (r)
       0      node3
         (l) /     \ (r)
            0       0
```
and call `cut(node1, node3->k)`, `cut` will take 4 chunks from tcache and free them again. Also notice that after calling `init_state` we have several chunks in tcache, usually about 4-6. We can rewrite pointer in 3rd chunk from tcache to address of `user` struct in the heap, so after calling `cut` `user` will be freed. And now we can create maze with size that equals to user's and rewrite `user->name` pointer to address of `stdin` in binary and leak libc.

After `stdin` there is a `maze` structure and pointer to `user`. Call `user_info` and rewrite pointer to `user` to a fake structure that we will create during this editing. In this fake structure we will set `user->name` to point to address of `strlen-8` in libc `.got`. Do one more edit and write `b'/bin/sh\0' + p64(libc.sym.system)`, so next calling of `printf(user->name)` will trigger `system('/bin/sh')`

[Exploit](./solve/spl.py)

## Domain
`game-51463dfee2ff6388.task.brics-ctf.ru`

## Flag

```
brics+{4c760273d8880a807d0c802781ba8793}
```

## Cloudflare
No
