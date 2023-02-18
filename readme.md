# Access Token Broker

AWS S3 access token broker written in Go and optimized for performance and scalability

## Build & Run

Build the go project

```bash
make build
```

```bash
cd bin
./atb
```

## API

### `/s3/sign`

query:

- `s3`: full path of the object including bucket name (`mybucket/my/full/path.txt`)
- `ttl`: time to live for the token in minutes

returns:

URL as string
