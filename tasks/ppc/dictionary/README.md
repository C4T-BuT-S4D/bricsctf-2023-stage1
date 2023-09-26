# ppc | Dictionary

## Information

> No dictionary as useless as this one has ever existed! More languages are planned to be supported soon.
> 
> https://dictionary-36b1430d4a6f1ab4.brics-ctf.ru

## Deploy

```sh
cd deploy
docker-compose -p dictionary up --build -d
```

Setup Caddy with HTTPS to proxy to port `2229` **without CF** (because task requires lots of requests which CF might count as DOS), enable ratelimiting on the reverse proxy.

## Public

Provide zip file: [public/dictionary.zip](public/dictionary.zip).

## TLDR

A simple SQL injection exists in the only API endpoint of the task, however the DB used is GenjiDB, which doesn't contain a whole lot of functions, leaving us only with the ability to exploit a boolean-based SQLi. The task is to automate this SQLi in some way since the flag is split into lots of parts.

## Writeup

The only endpoint of the server - [existsHandler](deploy/internal/web/handlers.go) queries the DB without properly escaping/sanitizing/preparing the 'word' query parameter. The mentioned handler, however, returns only 2 status codes without any body: 200 and 404 (well, technically, 500 is also returned), leaving no choice but to exfiltrate data using a boolean-based SQLi.

The flag is contained in a separate table, split into words character-by-character, where each character is prefixed with a string of random length, which complicates the attack, since GenjiDB contains no function like `substr` which would allow extraction of the flag characters. Flag extraction could've been done using a simple bruteforce of the words' characters, taking a long time, or using binary search over the word and flag alphabet (`[a-zA-Z0-9_-]`).

[Exploit](solve/solve.py)

## Domain

dictionary-36b1430d4a6f1ab4.brics-ctf.ru

## Cloudflare

No

## Flag

brics+{sqL_1nJect1on5_ar3_s1Mpl3_4nD_fUn_To_eXpl01t_16b08c0731324599}
