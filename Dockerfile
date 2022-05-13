FROM golang:1.18 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

# copy other sources & build
COPY . /app
RUN CGO_ENABLED=0 go build -a -o /app/display_manager cmd/main.go

FROM alpine AS runtime-env
COPY --from=build-env /app/display_manager /usr/local/bin/display_manager
ENTRYPOINT ["/usr/local/bin/display_manager"]