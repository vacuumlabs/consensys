#################################################
# STEP 1 build cache with Go modules cache
#################################################
FROM golang:1.16.2-alpine3.13 AS builder_cache
RUN apk update && apk add --no-cache git mercurial make build-base
WORKDIR /go/src/server
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
# "go mod download" downloads dependencies only when something changes in the go.mod or go.sum file
# which are cached via Docker's layer.
RUN go mod download

#################################################
# STEP 2 build the server
#################################################
FROM builder_cache AS builder
WORKDIR /go/src/server
COPY . .

# build the server
ARG project
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server $project

#################################################
# STEP 3 build a small image
#################################################
FROM scratch
COPY --from=builder /go/src/server/server /
WORKDIR /

# Run the server as a non-root/non-privileged user.
# We use just a "random" UID without a username and group because our server does not need them.
USER 10001

EXPOSE 8080
ENTRYPOINT ["/server"]
