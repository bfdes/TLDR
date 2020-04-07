# url-shortener

URL shortener and redirect service.

## Background

Redirect services work by mapping slugs in requests they receive to urls. The number of unique urls they can redirect is bounded by 

* the size of the slug's character set `b`
* the length of the longest slug `n`

and is given by the sum of a power-series that evaluates to `pow(b, n+1)-b / (b-1)`.

This service uses a character set of `"[A-Z][a-z][0-9]-_"` to generate about 930 million links before reaching slugs of six symbols. 

## Local development

### Requirements

* Go 1.11.5
* Docker 19.03.8
* Docker Compose 1.25.4

Install library dependencies with `go install`, and start the container dependencies with `docker-compose up`.

`go run .` starts the server.

## Usage

### Shorten a URL

```
curl http://localhost:<PORT>/api/links \
  --request POST \
  --data '{"url": "http://example.com"}'
# {"url": "http://example.com", "slug": "<SLUG>" }
```

### Redirect a link

```
curl http://localhost:<PORT>/<SLUG>
```

## Testing

Run unit and integration tests with `go test`, after starting the container dependencies.  

## Architecture

Each url to be shortened is persisted to an SQL database instance and assigned an integer ID from a auto-incrementing sequence.

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

Multiple links for the same url are accommodated to avoid wasting index space, compute and sequence values for every duplicate insert.
