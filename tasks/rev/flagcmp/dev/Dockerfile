FROM rust:1.72-slim-bookworm

RUN apt update && apt install -y curl

RUN cargo install wasm-pack@0.12.1

WORKDIR /wabt
RUN curl https://github.com/WebAssembly/wabt/releases/download/1.0.33/wabt-1.0.33-ubuntu.tar.gz -L -o wabt.tar.gz
RUN tar xzvf wabt.tar.gz

WORKDIR /build
COPY wasm .

RUN \
  --mount=type=cache,target=/usr/local/cargo/registry \
  --mount=type=cache,target=/build/target \
  wasm-pack build --target nodejs --out-name impl --no-typescript --no-pack

RUN /wabt/wabt-1.0.33/bin/wasm-strip ./pkg/impl_bg.wasm -o ./pkg/impl_bg.wasm
