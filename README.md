# Prometheus Kraken exporter

A Prometheus exporter for the [Kraken](https://www.kraken.com) bitcoin exchange.

## Usage
Run `make` to build the exporter Docker image.

To run with Prometheus and Grafana:
```
docker run -d -p 8080:8080 --name kraken-exporter -e "KEY=your_key" -e "SECRET=your_secret" basph/kraken-exporter:latest
docker run -d --net="host" -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml -p 9090:9090 prom/prometheus
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

Blog post about writing this exporter & deploying on Google Cloud: https://harenslak.nl/blog/monitoring-kraken-balance.
