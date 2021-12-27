Fireboard Exporter
==================

Export current Fireboard data in a way that Prometheus can scrape it.

Data is collected from the Fireboard cloud API rather than from the Fireboard directly because the data isn't exposed locally.

The Fireboard API rate limits requests to 200/hour. A scrape interval shorter than 18 seconds will cause you to exceed the rate limit.

The Fireboard API doesn't have a documented way to tell when data might be stale. At least one temperature probe needs to be connected or the `fireboard_up` metric will report as `0` and no other metrics will be returned (because the data might be stale).

## Installing

To install using go, run `go install github.com/pcfens/fireboard-exporter@latest`

To run with Docker, skip down to below.

## Running

Before running, you'll need an API token. The easiest way to retrieve it is with cURL
```bash
curl https://fireboard.io/api/rest-auth/login/ \            
    -X POST  \  
    -H 'Content-Type: application/json' \
    -d '{"username":"user@example.com","password":"password"}'
```

The response will look like `{"key":"292f783349256413248b7a132d34ba60d9c0faca"}`

To run the exporter, use `fireboard-exporter -key 292f783349256413248b7a132d34ba60d9c0faca` (with your key).

By default the exporter listens on port 8080.

### Demo

To run a demo using docker-compose, get an API token, then edit docker-compose.yml to add your token. Next, run `docker-compose up`.

After things start, you can [view probe temperatures](http://localhost:9090/graph?g0.expr=fireboard_probe_temperature_degrees&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=15m) in the Prometheus UI.

### Using Docker

To run in Docker use 

`docker run --rm -it -p 8080:8080 ghcr.io/pcfens/fireboard-exporter:main -key 292f783349256413248b7a132d34ba60d9c0faca`