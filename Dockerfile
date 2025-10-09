FROM index.docker.io/library/golang@sha256:e89daba6691cee7e338c17467e3d8777ff18d80ee62ff544808971fd51189282 AS build

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

FROM index.docker.io/library/alpine@sha256:85f2b723e106c34644cd5851d7e81ee87da98ac54672b29947c052a45d31dc2f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
