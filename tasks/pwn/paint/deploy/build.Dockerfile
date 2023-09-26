FROM ubuntu:22.04@sha256:b492494d8e0113c4ad3fe4528a4b5ff89faa5331f7d52c5c138196f69ce176a6

RUN apt update && apt install -y gcc

WORKDIR build/
COPY ./vuln.c .
RUN gcc -Wl,-z,relro,-z,lazy -o vuln vuln.c
