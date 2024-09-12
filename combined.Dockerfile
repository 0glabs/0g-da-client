FROM golang:1.21.1-alpine3.18 AS builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# build dispersal api server with local monorepo go modules
WORKDIR /
COPY . 0g-da-client
WORKDIR /0g-da-client
RUN make build

FROM alpine:3.18

COPY --from=builder /0g-da-client/disperser/bin/combined /usr/local/bin/combined

CMD ["combined"]