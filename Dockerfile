FROM golang@sha256:8ee9b9e11ef79e314a7584040451a6df8e72a66712e741bf75951e05e587404e AS build

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

FROM alpine@sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
