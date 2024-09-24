FROM scratch
WORKDIR /config
COPY --from=alpine:3 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/almanac", "serve", "--content-dir", "/content"]
COPY almanac /
