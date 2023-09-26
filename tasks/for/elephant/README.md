# for | Elephant

## Information
Sometimes, people become obsessed with things they like. This home server's owner was in love with a certain DBMS. So much so, that he entrusted all of his life's secrets to it. 

Why are we telling you this? Well, there are some secrets on that server which are of interest to us. 

We went through with breaking into his home to access the server. It was kept behind a door with a passcode lock, but the code was '12345' - maybe he didn't expect anybody to ever find out about the server? Anyways, we could only manage to swipe a RAM dump, but we can't be sure it's byte-to-byte perfect.

Help us recover the secret!

## Public
Provide 7z file: [for-elephant.7z](public/for-elephant.7z)

## Writeup
We are given a RAM dump of a Linux machine. Volatility 3's `banners` plugin finds the kernel version: `Ubuntu 5.15.0-84.93-generic`. [This script](solve/prepare_vol3.sh) will download this kernel with debug symbols from Ubuntu's ddeb repository and generate symbols for Volatility.

Now we can analyze the dump further. Notice two running processes: `psql` and `postgres`. Unfortunately, Volatility3's `linux.pslist.PsList --dump` didn't work for me, so I had to write a simple IDAPython [script](solve/load_segs.py) which would populate my IDB with chunks that I was able to obtain via `linux.proc.Maps --dump`.

It's wise to first look at `psql`. There we see that `SELECT secret('test', 'passwrod')` and `SELECT secret('test', 'password')` were performed. This is obviously some kind of a function, most likely (okay, definitely) a function in PL/pgSQL - Postgres' built-in procedural language.

Moving on to `postgres` processes (there are multiple) - the one to look at has PID 1002. It's the one that handles the connection and, spoiler, compiles PL/pgSQL sources.

Again, I loaded the process' memory into IDA. The 'Find binary' command wasn't working for me, so I wrote another [script](solve/bgrep.py) to scan for a byte sequence in all `.dmp` files. The string we're interested in is `secret(text,text)` - function's signature. These segments contain it:
```
./dmp/pid.1002.vma.0x56361fc27000-0x56361fe72000.dmp	0x000056361fcb6e80
./dmp/pid.1002.vma.0x56361fbc3000-0x56361fc27000.dmp	0x000056361fbd730d
```

Now, one is tempted to search for the sources. But one won't find them. At least, far their full form. An example:
```
((dub+1)::bit(128)), kek);
    uwu := uwu || substring(aba(substring(lol FROM dub*16+1), uou) FOR qop);

    RETURN uwu;
END;
```
Another example:
```
yay[i] := aaa(unu[i*4] || unu[i*4+1] || unu[i*4+2] || unu[i*4+3])
```

This is bad news. Do we have to find all of those chunks and somehow piece them together into original code? That's too much! Is there any other way to know what these functions do?

There is! But we have to dig into Postgres' source code. Particularly, [this directory](https://github.com/postgres/postgres/tree/REL_16_0/src/pl/plpgsql/src). Looking into `plpgsql.h`:
```c
/*
 * Complete compiled function
 */
typedef struct PLpgSQL_function
{
	char	   *fn_signature;
```
The function is compiled (and cached) before execution! Also:
```c
/*
 * SQL Query to plan and execute
 */
typedef struct PLpgSQL_expr
{
	char	   *query;			/* query string, verbatim from function body */
// ...
/*
 * Assign statement
 */
typedef struct PLpgSQL_stmt_assign
{
	PLpgSQL_stmt_type cmd_type;
	int			lineno;
	unsigned int stmtid;
	int			varno;
	PLpgSQL_expr *expr;
} PLpgSQL_stmt_assign;
```

It wouldn't be wrong to assume that the string from the second example above is actually some `PLpgSQL_expr`'s `query`. What if we look for
`0x000056361fcb6e80` in the dump?
```
 $ ./bgrep.py `python -c "print(0x000056361fcb6e80.to_bytes(8, 'little').hex())"`
./dmp/pid.1002.vma.0x56361fc27000-0x56361fe72000.dmp	0x000056361fc65b18
```
So, we found the only xref to `secret(text,text)` in the entire process memory. Taking into account what we learned from the C sources, we can deduce that `0x000056361fc65b18` is actually... a pointer to `PLpgSQL_function`.

Now comes the part where we frantically press <kbd>q</kbd><kbd>q</kbd><kbd>q</kbd><kbd>q</kbd> on any 8 bytes that look like a pointer. And we find other strings, like `_name` and `_password`.
```
seg007:000056361FC65D0E                 db    0
seg007:000056361FC65D0F                 db    0
seg007:000056361FC65D10                 dq offset off_56361FCBA3E0
seg007:000056361FC65D18                 db  68h ; h
seg007:000056361FC65D19                 db    1
seg007:000056361FC65D1A                 db    0
seg007:000056361FC65D1B                 db    0
```
```
seg007:000056361FCBA3DE                 db    0
seg007:000056361FCBA3DF                 db    0
seg007:000056361FCBA3E0 off_56361FCBA3E0 dq offset unk_56361FCB7FC8
seg007:000056361FCBA3E8                 dq offset unk_56361FCB8108
seg007:000056361FCBA3F0                 dq offset unk_56361FCB8550
seg007:000056361FCBA3F8                 dq offset unk_56361FCB85D8
seg007:000056361FCBA400                 dq offset unk_56361FCB9B20
seg007:000056361FCBA408                 db    0
seg007:000056361FCBA409                 db    0
seg007:000056361FCBA40A                 db    0
```
```
seg007:000056361FCB8107                 db    0
seg007:000056361FCB8108 unk_56361FCB8108 db    0                ; DATA XREF: seg007:000056361FCBA3E8↓o
seg007:000056361FCB8109                 db    0
seg007:000056361FCB810A                 db    0
seg007:000056361FCB810B                 db    0
seg007:000056361FCB810C                 db    1
seg007:000056361FCB810D                 db    0
seg007:000056361FCB810E                 db    0
seg007:000056361FCB810F                 db    0
seg007:000056361FCB8110                 dq offset aPassword     ; "_password"
seg007:000056361FCB8118                 db    0
seg007:000056361FCB8119                 db    0
```

Now, it would be great to somehow automate this. IDA has a C parser, so adapting `plpgsql.h` to 'compile' without other headers is one way to simplify the task. [Here](solve/c-dumper/plpgsql_ida.h)'s what I came up with. After setting the type:
```
seg007:000056361FC65B18 stru_56361FC65B18 dq offset aSecretTextText; fn_signature
seg007:000056361FC65B18                                         ; DATA XREF: seg007:000056361FCB9CE8↓o
seg007:000056361FC65B20                 dd 4019h                ; fn_oid
seg007:000056361FC65B24                 dd 2F3h                 ; fn_xmin
seg007:000056361FC65B28                 dw 0                    ; fn_tid.ip_blkid.bi_hi
seg007:000056361FC65B2A                 dw 2Bh                  ; fn_tid.ip_blkid.bi_lo
seg007:000056361FC65B2C                 dw 0Dh                  ; fn_tid.ip_posid
seg007:000056361FC65B2E                 db 2 dup(0)
seg007:000056361FC65B30                 dd PLPGSQL_NOT_TRIGGER  ; fn_is_trigger
seg007:000056361FC65B34                 dd 64h                  ; fn_input_collation
seg007:000056361FC65B38                 dq offset unk_56361FCBEBB0; fn_hashkey
seg007:000056361FC65B40                 dq offset unk_56361FCB6D70; fn_cxt
seg007:000056361FC65B48                 dd 19h                  ; fn_rettype
seg007:000056361FC65B4C                 dd 0FFFFFFFFh           ; fn_rettyplen
seg007:000056361FC65B50                 db 0                    ; fn_retbyval
seg007:000056361FC65B51                 db 0                    ; fn_retistuple
seg007:000056361FC65B52                 db 0                    ; fn_retisdomain
seg007:000056361FC65B53                 db 0                    ; fn_retset
seg007:000056361FC65B54                 db 0                    ; fn_readonly
seg007:000056361FC65B55                 db 66h                  ; fn_prokind
seg007:000056361FC65B56                 db 2 dup(0)
seg007:000056361FC65B58                 dd 2                    ; fn_nargs
seg007:000056361FC65B5C                 dd 0, 1, 62h dup(0)     ; fn_argvarnos
seg007:000056361FC65CEC                 dd 0FFFFFFFFh           ; out_param_varno
seg007:000056361FC65CF0                 dd 2                    ; found_varno
seg007:000056361FC65CF4                 dd 0                    ; new_varno
seg007:000056361FC65CF8                 dd 0                    ; old_varno
seg007:000056361FC65CFC                 dd PLPGSQL_RESOLVE_ERROR; resolve_option
seg007:000056361FC65D00                 db 0                    ; print_strict_params
seg007:000056361FC65D01                 db 3 dup(0)
seg007:000056361FC65D04                 dd 0                    ; extra_warnings
seg007:000056361FC65D08                 dd 0                    ; extra_errors
seg007:000056361FC65D0C                 dd 5                    ; ndatums
seg007:000056361FC65D10                 dq offset off_56361FCBA3E0; datums
seg007:000056361FC65D18                 dq 168h                 ; copiable_size
seg007:000056361FC65D20                 dq offset unk_56361FCBA398; action
seg007:000056361FC65D28                 dd 6                    ; nstatements
seg007:000056361FC65D2C                 db 0                    ; requires_procedure_resowner
seg007:000056361FC65D2D                 db 3 dup(0)
seg007:000056361FC65D30                 dq 0                    ; cur_estate
seg007:000056361FC65D38                 dq 0                    ; use_count
```

We could continue further with analyzing in IDA, but here's a cooler way to recover sources:

1. Wander into [pl_funcs.c](https://github.com/postgres/postgres/blob/REL_16_0/src/pl/plpgsql/src/pl_funcs.c);
2. Notice `plpgsql_dumptree(PLpgSQL_function *func)`;
3. Somehow run it on structs we just found
4. ???
5. PROFIT!

(These)[solve/c-dumper] sources do exactly that. When I first thought of this solution, I expected that it would be hard to 'rip out' the requied code. Actually, it was quite easy. Anyhow, the trick is to `mmap` `pid.1002.vma.0x56361fc27000-0x56361fe72000.dmp` to the right address and call `plpgsql_dumptree(0x000056361fc65b18)`. Voilà!
```

Execution tree of successfully compiled PL/pgSQL function secret(text,text):

Function's data area:
    entry 0: VAR _name            type text (typoid 25) atttypmod -1
    entry 1: VAR _password        type text (typoid 25) atttypmod -1
    entry 2: VAR found            type bool (typoid 16) atttypmod -1
    entry 3: VAR key              type bytea (typoid 17) atttypmod -1
                                  DEFAULT 'l1l(_name, _password)'
    entry 4: VAR _content         type bytea (typoid 17) atttypmod -1

Function's statements:
  5:BLOCK <<*unnamed*>>
  6:  ASSIGN var 4 := '_content := ('\x' || (SELECT content FROM secrets WHERE name = _name))::bytea'
  7:  ASSIGN var 4 := '_content := lll(_content, key)'
  9:  IF 'bbb(_password::bytea) <> ('\x' || (SELECT password FROM secrets WHERE name = _name))::bytea OR 
        bbb(substring(_content FROM 5)) <> substring(_content FOR 4)' THEN
 11:    RAISE level=21 message='Check the password!'
      ENDIF
 14:  RETURN 'convert_from(substring(_content FROM 5), 'UTF8')'
    END -- *unnamed*

End of execution tree of function secret(text,text)

```

Now we repeat the process for each new function we encounter. This [script](solve/find_func.py) finds a pointer (and formats it to paste into C code) given the function signature (we'll have to guess the argument types).

Having [dumped](solve/dump.txt) all the functions, we can notice the Rijndael S-box and the CRC-32 table. Reversing further, we discover that, in fact, AES-CTR is performed using a key derived via CRC-32 from string `name || password || 'p3pp3r'`.

The rows of the encrypted table can be found with `bgrep.py`:
```
 $ ./bgrep.py `python -c "print(b'secret'.hex())"`
./dmp/pid.1002.vma.0x7f77bb5ec000-0x7f77bb6ae000.dmp	0x00007f77bb5f9d04
./dmp/pid.1002.vma.0x7f77c699a000-0x7f77c69f2000.dmp	0x00007f77c69b3841
./dmp/pid.1002.vma.0x56361e1d7000-0x56361e46f000.dmp	0x000056361e1e9503
./dmp/pid.1002.vma.0x56361fc27000-0x56361fe72000.dmp	0x000056361fca4650
./dmp/pid.1002.vma.0x56361db4c000-0x56361dc15000.dmp	0x000056361dbc8577
./dmp/pid.1002.vma.0x56361fbc3000-0x56361fc27000.dmp	0x000056361fbd730d
./dmp/pid.1002.vma.0x7f77bb7d6000-0x7f77c4668000.dmp	0x00007f77bbf48a44
./dmp/pid.1002.vma.0x7f77c67dd000-0x7f77c6805000.dmp	0x00007f77c67f4061
./dmp/pid.1002.vma.0x56361e4a9000-0x56361e4df000.dmp	0x000056361e4aa50c
```

In `./dmp/pid.1002.vma.0x7f77bb7d6000-0x7f77c4668000.dmp` we find following:

```
00000000  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000010  f4 02 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000020  05 00 03 00 02 08 18 00  11 67 6f 6f 64 62 79 65  |.........goodbye|
00000030  63 31 33 38 31 32 63 35  30 63 36 38 30 38 37 31  |c13812c50c680871|
00000040  64 31 30 65 66 61 35 32  37 65 39 34 64 31 37 38  |d10efa527e94d178|
00000050  63 62 66 38 32 63 38 63  36 30 61 31 65 39 30 38  |cbf82c8c60a1e908|
00000060  36 13 65 65 36 61 39 62  66 66 00 00 00 00 00 00  |6.ee6a9bff......|
00000070  f4 02 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000080  04 00 03 00 02 08 18 00  0b 66 6c 61 67 00 00 00  |.........flag...|
00000090  08 05 00 00 66 30 37 62  32 39 32 39 63 66 39 61  |....f07b2929cf9a|
000000a0  39 37 39 39 65 30 64 66  34 38 65 38 37 32 37 30  |9799e0df48e87270|
000000b0  65 36 39 36 30 30 31 30  31 63 62 37 63 65 66 33  |e69600101cb7cef3|
000000c0  66 35 35 31 30 32 65 30  36 63 66 30 64 31 34 30  |f55102e06cf0d140|
000000d0  32 64 31 33 38 38 35 35  63 66 32 62 63 61 61 31  |2d138855cf2bcaa1|
000000e0  31 38 35 61 39 31 64 31  35 35 66 64 37 64 39 34  |185a91d155fd7d94|
000000f0  30 66 62 37 33 62 36 38  34 63 61 38 66 65 61 64  |0fb73b684ca8fead|
00000100  31 33 31 61 62 35 37 35  36 31 62 63 62 62 32 39  |131ab57561bcbb29|
00000110  36 65 37 31 66 34 62 39  33 61 31 63 62 65 30 66  |6e71f4b93a1cbe0f|
00000120  33 64 61 66 31 61 36 34  33 32 34 64 38 39 37 65  |3daf1a64324d897e|
00000130  61 30 62 34 39 65 62 63  66 33 66 38 38 65 37 65  |a0b49ebcf3f88e7e|
00000140  65 34 32 65 63 39 31 32  34 33 66 32 66 34 65 65  |e42ec91243f2f4ee|
00000150  35 39 39 33 39 31 62 65  39 35 64 61 62 33 64 62  |599391be95dab3db|
00000160  64 32 65 62 63 35 30 39  38 35 31 36 31 35 61 34  |d2ebc509851615a4|
00000170  61 63 66 36 64 63 66 62  61 37 65 35 35 61 34 63  |acf6dcfba7e55a4c|
00000180  31 39 34 39 39 32 33 65  37 63 65 65 30 36 31 37  |1949923e7cee0617|
00000190  65 37 31 65 33 35 30 31  31 61 64 65 39 36 32 36  |e71e35011ade9626|
000001a0  36 63 34 35 66 36 33 30  34 39 35 38 33 66 30 38  |6c45f63049583f08|
000001b0  31 64 61 64 30 63 64 62  30 62 36 39 34 35 32 34  |1dad0cdb0b694524|
000001c0  32 65 61 61 38 61 63 61  64 62 38 32 35 62 65 31  |2eaa8acadb825be1|
000001d0  31 30 13 38 39 38 64 65  30 30 39 00 00 00 00 00  |10.898de009.....|
000001e0  f4 02 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
000001f0  03 00 03 00 02 08 18 00  0f 73 65 70 74 31 33 00  |.........sept13.|
00000200  38 03 00 00 34 66 39 36  63 35 64 65 35 61 32 30  |8...4f96c5de5a20|
00000210  33 38 30 36 30 34 31 36  66 31 37 37 65 34 35 35  |38060416f177e455|
00000220  31 37 33 61 63 64 65 66  64 63 63 66 65 38 39 37  |173acdefdccfe897|
00000230  33 64 38 34 65 31 39 38  64 39 31 30 63 33 63 37  |3d84e198d910c3c7|
00000240  61 30 33 33 32 37 61 61  65 31 66 34 37 37 38 36  |a03327aae1f47786|
00000250  33 63 33 34 30 31 64 39  66 61 36 63 63 34 39 66  |3c3401d9fa6cc49f|
00000260  66 34 30 61 32 30 61 36  62 66 61 65 37 37 63 37  |f40a20a6bfae77c7|
00000270  61 66 35 31 63 63 37 62  63 30 34 64 34 30 62 37  |af51cc7bc04d40b7|
00000280  66 30 61 37 39 33 36 64  39 33 39 34 66 37 38 64  |f0a7936d9394f78d|
00000290  62 33 63 33 66 33 36 64  34 33 63 63 37 39 36 64  |b3c3f36d43cc796d|
000002a0  32 65 61 31 39 30 61 62  66 31 36 38 39 39 66 64  |2ea190abf16899fd|
000002b0  64 37 61 34 33 39 61 35  65 61 30 62 31 64 36 63  |d7a439a5ea0b1d6c|
000002c0  30 34 61 35 31 35 64 63  32 33 61 33 65 62 13 34  |04a515dc23a3eb.4|
000002d0  33 66 66 39 62 62 64 00  f4 02 00 00 00 00 00 00  |3ff9bbd.........|
000002e0  00 00 00 00 00 00 00 00  02 00 03 00 02 08 18 00  |................|
000002f0  0d 74 72 6f 6c 6c 87 38  38 62 32 35 39 61 36 39  |.troll.88b259a69|
00000300  33 35 33 32 31 64 63 64  64 33 31 32 63 37 32 65  |35321dcdd312c72e|
00000310  31 36 66 38 39 63 35 63  35 39 36 61 31 34 35 61  |16f89c5c596a145a|
00000320  39 63 61 65 35 62 35 33  36 66 30 66 32 63 63 30  |9cae5b536f0f2cc0|
00000330  65 33 66 66 30 61 39 33  38 13 38 35 37 61 65 36  |e3ff0a938.857ae6|
00000340  34 38 00 00 00 00 00 00  f4 02 00 00 00 00 00 00  |48..............|
00000350  00 00 00 00 00 00 00 00  01 00 03 00 02 09 18 00  |................|
00000360  0b 74 65 73 74 63 38 61  65 33 39 66 65 64 37 31  |.testc8ae39fed71|
00000370  64 64 35 66 36 66 34 64  62 38 39 36 34 62 32 39  |dd5f6f4db8964b29|
00000380  38 64 38 38 65 63 63 64  30 62 64 35 37 30 30 38  |8d88eccd0bd57008|
00000390  38 37 39 33 64 35 13 33  35 63 32 34 36 64 35 00  |8793d5.35c246d5.|
000003a0  00 00 00 00 f0 2b 51 01  00 00 04 00 c4 00 18 01  |.....+Q.........|
000003b0  00 20 04 20 00 00 00 00  48 9f 68 01 90 9e 68 01  |. . ....H.h...h.|
```

`35c246d5` = `crc32('password')`. Taking a hint from the task description (about the passcode being '12345'), let's try to [break](solve/find_pass.py) hash `898de009`. The password is `soundblaster54`, found in rockyou.txt.

Now, using this [script](solve/decrypt.py), we can get our flag.
```
 $ ./decrypt.py flag soundblaster54 f07b2929cf9a9799e0df48e87270e69600101cb7cef3f55102e06cf0d1402d138855cf2bcaa1185a91d155fd7d940fb73b684ca8fead131ab57561bcbb296e71f4b93a1cbe0f3daf1a64324d897ea0b49ebcf3f88e7ee42ec91243f2f4ee599391be95dab3dbd2ebc509851615a4acf6dcfba7e55a4c1949923e7cee0617e71e35011ade96266c45f63049583f081dad0cdb0b6945242eaa8acadb825be110
Yesterday, I could't remember the 12th letter of the fl4g. So I'll put it here. brics+{7d46b73adab228de671ac5ef64444ea50a8eb6dee65ed2c8414ccd4c08_el3ph4nt}
```

## Domains
None

## Cloudflare
N/A

## Flag
`brics+{7d46b73adab228de671ac5ef64444ea50a8eb6dee65ed2c8414ccd4c08_el3ph4nt}`

