FROM index.docker.io/library/golang@sha256:217cb265b15f1bce711dda88e3dd2302da61b8b1096d6afe15727a7e96719dd1 AS build

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
