# web | ChadGPT

## Information

> We all love innovations! You have an opportunity to test the state-of-the-art LLM and play with it!
>
> https://chadgpt-647150b7b2303d42.task.brics-ctf.ru

## Deploy

```sh
cd deploy
docker-compose -p web_chadgpt up --build -d
```

Setup http reverse proxy to port `5000`.

## Public

Provide zip file: [public/web-chadgpt.zip](public/web-chadgpt.zip).

## TLDR

Exploit the SQL-injection protected by WAF by providing JSON that will be parsed by `GoJay` library and
standard `encoding/json` library differently.

## Writeup

We can find an SQL-injection inside the target application: the payload 'q' parameter is not properly escaped.
But there is a WAF that escapes all strings inside the JSON body.

Hovewer, the are two facts that can lead to the solution:
1. Quotes are not properly escaped.
2. WAF will proxy the request if it failed to decode the JSON body.

This allow two solutions:

### Solution 1
We can provide the `\'`, as `'` will be replaced with `''` we will get the `\''` as a result. Which will allow us to use the SQL-injection.

Use this as an input: `\' union select flag from flags -- `

[Exploit](solve/solve_sql.py)

### Solution 2
We can notice that the target application uses [gojay](github.com/francoispqt/gojay) library to parse the JSON body, but the WAF uses standard `encoding/json` library.

Having this in mind, we can find a JSON that will be parsed by the `gojay` library, but won't be parsed by the standard library.

One of the key differences is that GoJay library allows anything (including comment) after it parses all JSON fields needed.
So we can provide a JSON like this as a body:

```json
{
  "maxTokens": 123,
  "q": "asd'OR 1=2 UNION SELECT flag FROM flags -- ",
  asd
}
```
[Exploit](solve/solve_gojay.py)

## Domain
`chadgpt-647150b7b2303d42.task.brics-ctf.ru`

## Flag

`brics+{p14y1Ng_bM7H_C4N_y0u_F33l_mY_H34r7}`

## Cloudflare
Yes
