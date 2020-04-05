# url-shortener

URL shortener and redirect service.

## Background

Redirect services work by mapping slugs in requests they receive to urls. The number of unique urls they can redirect is bounded by 

* the size of the slug's character set `b`
* the length of the largest slug `n`

and is given by `math.pow(b, n)`.

This service uses a character set of `"[A-Z][a-z][0-9]-_"` to generate about 900 million links before reaching slugs of six symbols. 

## Local development

### Requirements

* Go 1.11.5
* Docker 19.03.8
* Docker Compose 1.25.4

Install library dependencies with `go install`. In the project directory run `docker-compose up` to start the Postgres container, and `go run .` to start the server.

Default connection settings are defined in `main.go`.

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

Run unittests with `go test`. 

## Architecture

Each url to be shortened is persisted to an SQL database instance and assigned an integer ID from a auto-incrementing sequence.

To action a redirect request, the request slug is decoded to obtain an ID and its url is sought in a key-value cache followed by the SQL database. 

If viewing figures are desired, "LINK_VIEWED" events can be dispatched to an unordered queue for asynchronous processing by serverless-functions.

### Advantages

* Since links have to be immutable, cache-invalidation is not an issue. Eviction only needs to be carried out when memory constraints are met. 

* Multiple links can be created for the same url. 

### Shortcomings

* For performance reasons, ID sequences are not gapless. So multiple links for the same url are accommodated to avoid wasting index space, compute and sequence values for every duplicate insert.

* Viewing figures may be stale because they are written in an asynchronous manner in order to avoid hitting the database for every redirect request.
