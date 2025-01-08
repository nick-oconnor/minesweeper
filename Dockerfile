FROM golang@sha256:b4743faf9518405c68649c29f1c9e29f43872a5e882c61e411347e73ef64a0b5 AS build

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

FROM alpine@sha256:483f502c0e6aff6d80a807f25d3f88afa40439c29fdd2d21a0912e0f42db842a

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
