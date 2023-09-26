# rev | FlagCMP

## Information

> Thanks to modern technologies we finally have a crossplatform way to check that your flag is, indeed, correct.
>
> https://flagcmp-b8f9ceaf86d7846a.brics-ctf.ru

## Deploy

```sh
cd deploy
docker-compose -p flagcmp up --build -d
```

## Public

Provide zip file: [public/flagcmp.zip](public/flagcmp.zip).

## TLDR

Reverse engineering the WASM module shows that the user's guess is checked using a [Befunge-93](https://github.com/catseye/Befunge-93/blob/master/doc/Befunge-93.markdown) program which is constructed on the fly, and the construction mechanism allows the user input to escape the original string execution mode in order to execute arbitrary Befunge commands. Using a Befunge quine the original flag checking code can be dumped and then analyzed to retrieve the flag.

## Writeup

1. Use [wasm2ida](https://github.com/vient/wasm2ida) to convert the wasm module into a binary which is properly decompilable by IDA so that you can read the logic of the program. Since `run` and `set_flag` are exported symbols, it's easy to find and reverse engineer them. Specifically, `run` is seen to write lots of single-byte codes to some location, and then execute a loop with a big switch-case statement for all the different codes. It's not hard to guess that this challenge is a yet-another-vm-reverse-challenge.
2. To simplify reverse engineering, you can use the Chrome debugger and memory view to debug multiple interesting functions and see the data they're operating on, for example, the main interpreter function with the switch statement, and the transformation of user input.
3. By analyzing how the user input is prepared and inspecting it in Chrome raw memory view, you can see that for some reason some characters are transformed differently than others, and have the same codes as the main VM program.
4. In order to dump the main VM program, you can use the `wasm2wat` and `wat2wasm` tools from the [WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt) for writing a simple hook which would output instructions when placed in the main interpreter function:
   ```
   (func $hook (type 3) (param $a i32)
       local.get $a
       i32.const 2
       call 28
     )
   ```

   And by modifying one of the imports to console.log instead of whatever it does
5. Once you dump all the executing instructions, you will notice that the transformed user input is copied here in the same way as it was earlier, which allows you to pretty much write any of the VM commands. Either by reverse engineering the VM opcodes or by fuzzing possible character -> code transformations, you could find that `"` would escape the enabled StringMode and allow you to execute your own commands on the VM.
6. From here on you just need to write a Befunge-93 quine which would print the whole program, and then you can extract the flag comparisons from it to get the flag.

## Domain

flagcmp-b8f9ceaf86d7846a.brics-ctf.ru

## Cloudflare

Yes

## Flag

brics+{c3rtif1ed_es0l4ng_expl0it_d3vel0per_c65596e73d72cac7}
