FROM docker.io/golang:1.17.0-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN  go build -o build/server src/main.go

FROM scratch AS bin
COPY --from=build /src/build/server /