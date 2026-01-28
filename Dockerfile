FROM index.docker.io/library/golang@sha256:72613572613050e1c6da0a585daf8f7d4174eb17a0c9c3a539633d677e3e2979 AS build

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

FROM index.docker.io/library/alpine@sha256:59855d3dceb3ae53991193bd03301e082b2a7faa56a514b03527ae0ec2ce3a95

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
