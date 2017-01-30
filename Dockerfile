FROM scratch
COPY service /
ENTRYPOINT ["/service"]