# repo-updater:config arch=amd64
# repo-updater:container tag_include=^1\.25\.\d+-alpine3\.23$ version=1.25.7-alpine3.23
FROM index.docker.io/library/golang@sha256:724e212d86d79b45b7ace725b44ff3b6c2684bfd3131c43d5d60441de151d98e AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY field/ ./field
COPY matrix/ ./matrix
COPY solver/ ./solver
COPY *.go ./
COPY *.pgo ./

# repo-updater:static container_tag=true version=v1.1.1
RUN go build -v -o /minesweeper

# repo-updater:container tag_include=^3\.23\.\d+$ version=3.23.3
FROM index.docker.io/library/alpine@sha256:59855d3dceb3ae53991193bd03301e082b2a7faa56a514b03527ae0ec2ce3a95

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
