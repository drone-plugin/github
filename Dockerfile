FROM alpine:latest

ADD  main /bin/
RUN apk add --no-cache git ca-certificates
ENTRYPOINT /bin/main