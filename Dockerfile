FROM golang:1.14

WORKDIR /

COPY ./go.* ./

RUN go mod download

COPY . .

CMD ./scripts/maze_benchmark.sh
