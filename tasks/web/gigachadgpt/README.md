# web | GigaChadGPT

## Information

> We improved security of our LLM by asking the LLM how to make it secure!
>
> https://gigachadgpt-647150b7b2303d42.brics-ctf.ru

## Deploy

```sh
cd deploy
docker-compose -p web_gigachadgpt up --build -d
```

Setup http reverse proxy to port `5000`.

## Public

Provide zip file: [public/web-gigachadgpt.zip](public/web-gigachadgpt.zip).

## TLDR

Exploit the SQL-injection protected by WAF by providing JSON that will be parsed by `GoJay` library and
standard `encoding/json` library differently.

## Writeup

This is the follow-up challenge to `chadgpt`. In this challenge instead of escaping the characters WAF will ban all the requests that have characters outside the whitelist.
However, the WAF will not check the strings inside the JSON body if it failed to read the JSON body, it will proxy it as is.

And we can notice that the target application uses [gojay](github.com/francoispqt/gojay) library to parse the JSON body, but the WAF uses standard `encoding/json` library.

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
[Exploit](solve/solve.py)

## Domain
`gigachadgpt-647150b7b2303d42.task.brics-ctf.ru`

## Flag

`brics+{play1ngWith1ChadBotAndmySQL1sFun}`

## Cloudflare
Yes
