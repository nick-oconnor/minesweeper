FROM index.docker.io/library/golang@sha256:6d6d1e4e530e8512543843504590c86b30524dd8644953c3435fa5b3396ae39c AS build

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

FROM index.docker.io/library/alpine@sha256:a107a3c031732299dd9dd607bb13787834db2de38cfa13f1993b7105e4814c60

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
