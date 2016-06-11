FROM scratch
COPY udp-proxy /
ENTRYPOINT ["/udp-proxy"]
