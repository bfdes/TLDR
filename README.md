# url-shortener

URL shortener and redirect service.

## Background

Redirect services work by mapping slugs in requests they receive to urls. The number of unique urls they can redirect is bounded by 

* the size of the slug's character set `b`
* the length of the longest slug `n`

and is given by the sum of a power-series that evaluates to `pow(b, n+1)-b / (b-1)`.

This service uses a character set of `"[A-Z][a-z][0-9]-_"` to generate about 930 million links before reaching slugs of six symbols. 

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

* Go 1.11.5
* Docker Engine 19.03.8
* Docker Compose 1.25.4

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

## Architecture

Theoretical capacity cannot be met by using uuids or url hashing to generate the slug. Instead, each link's url is persisted to an SQL database instance and assigned an integer ID from a auto-incrementing sequence. The ID is then transformed into a slug by base-62 encoding.

| id   | url                  |
|------|----------------------|
| 1    | "http://example.com" |
| ...  | ...                  |
| 1504 | "http://google.com"  |

To action a redirect request, the request slug is decoded to obtain an ID and its url is sought in a key-value cache followed by the database. 

| slug | url                  |
|------|----------------------|
| "1"  | "http://example.com" |
| ...  | ...                  |
| "go" | "http://google.com"  |

Since links have to be immutable, cache-invalidation is not an issue. Eviction only needs to be carried out when memory constraints are met.

Multiple links for the same url are accommodated to avoid wasting index space, compute and sequence values for every duplicate insert. In practice duplicate links are unlikely to be made, so the theoretical capacity should be met.
