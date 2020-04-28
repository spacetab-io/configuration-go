# Dockerfile for test image
FROM spacetabio/docker-test-golang:1.14-1.0.2

COPY . /app
RUN make tests