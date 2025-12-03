FROM index.docker.io/library/golang@sha256:18a5f71675ef62af731b00ac0bd22dd54133365ec0558bd93e203c578afc7e18 AS build

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
