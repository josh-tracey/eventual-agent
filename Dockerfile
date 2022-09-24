FROM docker.io/golang:1.18.0-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN sh ./scripts/generate-protos.sh
RUN sh ./script/build.sh

FROM scratch AS bin
COPY --from=build /src/build/eventual-agent /
