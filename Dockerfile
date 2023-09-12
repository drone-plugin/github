FROM alpine:latest

ADD  main /bin/
RUN apk add --no-cache git openssh
ENTRYPOINT /bin/main