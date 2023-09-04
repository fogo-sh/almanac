FROM scratch
WORKDIR /config
ENTRYPOINT ["/almanac", "serve", "--content-dir", "/content"]
COPY almanac /
