#
# File: Dockerfile
# Project: kube-informer
#

# Builder stage
FROM golang:1.17 AS builder
WORKDIR /buildsource
COPY . /buildsource
RUN set -eux; \
    make linux-binary

# Final stage
FROM registry.access.redhat.com/ubi8/ubi-micro
COPY --from=builder /buildsource/bin/inform /usr/local/bin/inform
CMD ["inform", "--in-cluster"]
