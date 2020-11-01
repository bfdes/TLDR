FROM golang as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
LABEL maintainer="Bruno Fernandes <bfdes@users.noreply.github.com>"
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
