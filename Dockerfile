# bundle the admin webapp
FROM node:16.9.0-alpine3.14 AS admin_builder
RUN apk update && apk upgrade --no-cache
RUN apk add git

RUN git clone https://github.com/superseriousbusiness/gotosocial-admin
WORKDIR /gotosocial-admin

RUN npm install
RUN node index.js

# ----
# golang build of gotosocial application binary
FROM golang:1.16-buster as application_builder
ARG GOARCH=

RUN mkdir /build
WORKDIR /build

# we can't do this trick because its vendored.
# COPY go.mod go.sum /build
# RUN go mod download  

COPY . /build

RUN scripts/build.sh


# ----
# put the build results from previous steps together into a lean docker image for distribution 
FROM alpine:3.14.2 AS executor
RUN apk update && apk upgrade --no-cache

# copy over the binary from the first stage
RUN mkdir -p /gotosocial/storage
COPY gotosocial /gotosocial/gotosocial

# copy over the web directory with templates etc
COPY web /gotosocial/web

# copy over the admin directory
COPY --from=admin_builder /gotosocial-admin/public /gotosocial/web/assets/admin

# make the gotosocial group and user
RUN addgroup -g 1000 gotosocial
RUN adduser -HD -u 1000 -G gotosocial gotosocial

# give ownership of the gotosocial dir to the new user
RUN chown -R gotosocial gotosocial /gotosocial

# become the user
USER gotosocial

WORKDIR /gotosocial
ENTRYPOINT [ "/gotosocial/gotosocial", "server", "start" ]
