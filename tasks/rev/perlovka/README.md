# rev | perlovka

## Information

> We found a very old malware. Could you reverse it ?

## Deploy

No deploy

## Public

Provide the binary: [public/perlovka](public/perlovka).

## TLDR

Reverse engineering the Perl script compiled to binary using [B:CC](https://perldoc.perl.org/5.8.7/B::CC).

## Writeup

1. Perl script is checking flags by evaling string formatted using user-input.
2. Each string is decrypted using key = `16 * user_input[i // 2]`, so we can brute-force the key and decrypt the Perl strings.
3. Each string has the flag letter inside it.


## Cloudflare

No

## Flag

`brics+{perlovkaOchenVkusnAya}`

