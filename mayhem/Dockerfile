FROM golang:1.17 as builder

COPY . /sveltin
WORKDIR /sveltin/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN cd fuzz_sveltin_common && go-fuzz-build -libfuzzer -o fuzz_sveltin_common.a && \
    clang -fsanitize=fuzzer fuzz_sveltin_common.a -o fuzz_sveltin_common.libfuzzer

RUN cd fuzz_sveltin_utils && go-fuzz-build -libfuzzer -o fuzz_sveltin_utils.a && \
    clang -fsanitize=fuzzer fuzz_sveltin_utils.a -o fuzz_sveltin_utils.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /sveltin/mayhem/fuzz_sveltin_common/fuzz_sveltin_common.libfuzzer /
COPY --from=builder /sveltin/mayhem/fuzz_sveltin_utils/fuzz_sveltin_utils.libfuzzer /