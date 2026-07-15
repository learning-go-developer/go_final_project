FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -tags netgo \
    -ldflags="-s -w" \
    -o /bin/scheduler main.go

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/scheduler .
COPY --from=builder /src/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db

EXPOSE 7540

ENTRYPOINT ["./scheduler"]