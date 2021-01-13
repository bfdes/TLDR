# url-shortener

![GitHub Actions](https://github.com/bfdes/url-shortener/workflows/Build/badge.svg)
[![Codecov](https://codecov.io/gh/bfdes/url-shortener/branch/master/graph/badge.svg)](https://codecov.io/gh/bfdes/url-shortener)

URL shortener and redirect service [designed for](https://www.notion.so/URL-shortening-8272c692648143698859d9f3524a8b5e#a2becb53582444cfb3e9cce1dd8978ba) low-latency, read-heavy use.

## Usage

A local copy of the application can be built and depolyed with `docker-compose up`.

### Shorten a URL

```
curl http://localhost:8080/api/links \
  --request POST \
  --data '{"url": "http://example.com"}'
# {"url": "http://example.com", "slug": "<SLUG>" }
```

### Redirect a link

```
curl http://localhost:8080/<SLUG>
```

## Local development

### Requirements

* Go 1.11.x
* Docker Engine 19.x
* Docker Compose 1.x

Run the following commands within the repository root:

```bash
docker-compose up -d cache database
# Starts container dependencies in the background

go install
# Installs library dependencies

go run .
# Starts the server
```

Run unit and integration tests with `go test`, after starting the container dependencies.
