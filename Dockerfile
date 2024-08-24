FROM golang:alpine AS builder


ADD go.mod .
RUN go mod download

COPY . .
RUN go build -o /build/service ./cmd/todo/main.go

FROM scratch

LABEL authors="Ildar Galeev"


COPY --from=builder /build/service /app/service
COPY template.config.yml /app/config.yml

#CMD ["/app/service"]