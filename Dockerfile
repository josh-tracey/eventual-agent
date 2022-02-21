FROM docker.io/golang:1.17.0-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN sh ./build

FROM scratch AS bin
COPY --from=build /src/build/eventual-agent /
