# cardvalidate

Go service for credit card validation.

## Known credit card issuers

Table of IIN ranges used for validation/identification:
| Issuer | IIN ranges | Card number length |
| --- | --- | --- |
| American Express | 34, 37 | 15 |
| Diners Club | 30, 36, 38, 39 | 14 |
| Discover | 6011, 644-649, 65 | 16 |
| JCB | 3528–3589 | 16 |
| MasterCard | 51-55, 2221–2720 | 16 |
| UnionPay | 62 | 16-19 |
| Visa | 4 | 16 |

## Setup

You'll need to have Go v1.22 or newer and Docker installed to set up the API. Once you've got
all of the requirements, run:
```bash
# Run tests
make test  # or go test -v ./...

# Build the executable
make build  # or go build -o ./bin/server ./cmd/server

# Run the server
./bin/server  # That's it!
```

Running inside Docker is as simple as the previous example:
```bash
docker build -t cardvalidate .
docker run --rm --name cardvalidate -p 8000:8000 cardvalidate
```

## API

If everything went successfully, the server should be listening on local port 8000.
Visit `http://localhost:8000/docs` page to see the Swagger UI documentation for the REST API.
