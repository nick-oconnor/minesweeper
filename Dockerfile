FROM golang@sha256:44a637dca7f3632232d2a21e6d5091ba5e2bbdc4830be97f54eee76c204dee7f AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY field/ ./field
COPY matrix/ ./matrix
COPY solver/ ./solver
COPY *.go ./
COPY *.pgo ./

RUN go build -v -o /minesweeper

FROM alpine@sha256:eafc1edb577d2e9b458664a15f23ea1c370214193226069eb22921169fc7e43f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
