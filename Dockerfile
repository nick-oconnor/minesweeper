FROM index.docker.io/library/golang@sha256:301db1784c200c7ec221a84672d0bd3ffdcf6c873bc9806482bbe730910924fb AS build

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

FROM index.docker.io/library/alpine@sha256:1882fa4569e0c591ea092d3766c4893e19b8901a8e649de7067188aba3cc0679

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]
