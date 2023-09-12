FROM alpine:latest

ADD  main /bin/
RUN apk add --no-cache git ca-certificates openssh
ENTRYPOINT /bin/main