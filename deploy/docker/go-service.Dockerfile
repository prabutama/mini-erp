ARG SERVICE

FROM golang:1.24-alpine AS build
ARG SERVICE
WORKDIR /src
RUN apk add --no-cache git
COPY go.work go.work.sum ./
COPY services ./services
RUN cd services/${SERVICE} && go mod download
RUN cd services/${SERVICE} && CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${SERVICE}

FROM alpine:3.22 AS migrate
ARG TARGETARCH
RUN apk add --no-cache ca-certificates curl tar
RUN curl -fsSL "https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-${TARGETARCH}.tar.gz" \
    | tar -xz -C /usr/local/bin migrate

FROM alpine:3.22
ARG SERVICE
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates
COPY --from=build /out/app /app/app
COPY --from=migrate /usr/local/bin/migrate /usr/local/bin/migrate
COPY services/${SERVICE}/migrations /app/migrations
USER app
ENTRYPOINT ["/app/app"]
