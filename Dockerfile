FROM golang:alpine as builder
WORKDIR /app
COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o artefact .

FROM alpine:latest
ENV TZ=America/Sao_Paulo
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone
WORKDIR /app
COPY --from=builder /app/artefact .
CMD ["/app/artefact"]