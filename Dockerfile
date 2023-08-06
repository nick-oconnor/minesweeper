FROM golang@sha256:29a2b239b10daa321685e58dd58ca7f5e04639f7fc99ba5af84975e2796b7cd7 AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY field/ ./field
COPY solver/ ./solver
COPY *.go ./

RUN go build -v -o /minesweeper

FROM scratch

COPY --from=build /minesweeper /

ENTRYPOINT ["/minesweeper"]
