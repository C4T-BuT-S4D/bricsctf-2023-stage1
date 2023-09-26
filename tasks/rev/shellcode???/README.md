# rev | shellcode???

## Information

> Write shellcode in reverse??? The author is clearly out of his mind...


## Deploy

No need

## Public

Provide zip file: [public/flagcmp.zip](public/shellcode).

## TLDR
Simple printf vm

## Writeup
It can be understood that in the task, the printf modifiers are being overridden, and each of them is responsible for a specific operation in the VM. The same operation is repeated 52 times:
```
# PUSH X 
# PUSH Y 
# POW Y^2 
# OVER
# POW X^3
# SUB
# OVER
# mul
# ADD FROM VALUE
# sub
```
Values from the stack are placed into the equation of an elliptic curve with specific coefficients, and then the result is compared with a given value.

```
%p 98
%p 125
%p 106
%p 111
%p 100
%p 95
%p 48
%p 111
%p 95
%p 71
%p 117
%p 84
%p 95
%p 98
%p 49
%p 115
%p 116
%p 104
%p 51
%p 95
%p 114
%p 115
%p 112
%p 65
%p 49
%p 95
%p 117
%p 115
%p 95
%p 106
%p 48
%p 116
%p 95
%p 110
%p 97
%p 110
%p 95
%p 99
%p 48
%p 85
%p 95
%p 121
%p 48
%p 108
%p 123
%p 76
%p 115
%p 43
%p 105
%p 99
%p 98
%p 114
```
By iterating through the values of variables X and Y, we obtain the necessary values on the stack, which will then be passed to the VM.

```
b}jod_0o_GuT_b1sth3_rspA1_us_j0t_nan_c0U_y0l{Ls+icbr
```
To restore the desired order, let's swap every 2 characters and reverse the string.

## Domain

No

## Cloudflare

No

## Flag

brics+{L0l_y0U_can_n0t_jus1_pArs3_th1s_buT_G0od_job}