FROM golang:1.23.9-alpine AS builder

WORKDIR /usr/src/app

ARG COMMIT_SHA
ARG COMMIT_TAG
ARG BUILD_TIMESTAMP

RUN apk add --no-cache ca-certificates git

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

#RUN go run github.com/google/go-licenses@latest report --template licenses.tpl --include_tests . 2>/dev/null > files/statics/licenses.json

RUN go build -v -o bin/quantum-go -ldflags="-s -w" cmd/main.go
RUN go build -v -o bin/quantum-server-go -ldflags="-s -w" cmd/quantum-server-go/main.go
RUN go build -v -o bin/quantum-client-go -ldflags="-s -w" cmd/quantum-client-go/main.go


#FROM builder as test
#
#RUN apk add --no-cache git
#RUN go install gotest.tools/gotestsum@latest

FROM scratch AS quantum-go

WORKDIR /usr/local/bin

COPY --from=tzdata /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/src/app/bin/quantum-go /usr/local/bin/quantum-go
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["quantum-go"]

FROM scratch AS quantum-go-server

WORKDIR /usr/local/bin

COPY --from=tzdata /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/src/app/bin/quantum-server-go /usr/local/bin/quantum-server-go
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["quantum-server-go"]

FROM alpine:3.20.6 AS tzdata
ENV DATE=20250512
RUN apk add --no-cache tzdata

FROM scratch AS quantum-client-go

WORKDIR /usr/local/bin

COPY --from=tzdata /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/src/app/bin/quantum-client-go /usr/local/bin/quantum-client-go
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["quantum-client-go"]
