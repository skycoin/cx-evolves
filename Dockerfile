FROM alpine:3.4

WORKDIR /

COPY ./go.* ./

COPY ./cx-evolves .

COPY ./server .

COPY ./scripts/maze_benchmark.sh .
COPY ./scripts/constants_benchmark.sh .
COPY ./scripts/evens_benchmark.sh .
COPY ./scripts/evens_v2_benchmark.sh .

ENTRYPOINT ["sh","/maze_benchmark.sh"]
