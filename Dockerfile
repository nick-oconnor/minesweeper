FROM golang@sha256:53443fdc64453f971b5c82374d86945b90c053880131b13ee50704001e77f23a AS build

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
