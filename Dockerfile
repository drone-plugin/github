FROM bitnami/git

ADD  main /bin/
ENTRYPOINT /bin/main