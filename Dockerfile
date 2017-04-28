FROM scratch
COPY test-service /
ENTRYPOINT ["/test-service"]