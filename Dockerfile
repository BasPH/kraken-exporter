# docker build -t basph/kraken-exporter .
# docker run -d -p 8080:8080 -e "KEY=xxxxx" -e "SECRET=xxxxx" --name kraken-exporter basph/kraken-exporter

FROM golang:1.9.2 AS builder
WORKDIR /go/src/github.com/BasPH/kraken-exporter
COPY . .
RUN apt-get update && apt-get install -y upx-ucl
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' && \
    upx --brute kraken-exporter

FROM scratch
COPY --from=builder /go/src/github.com/BasPH/kraken-exporter .
EXPOSE 8080
ENTRYPOINT ["./kraken-exporter"]