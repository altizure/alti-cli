FROM scratch
COPY alti-cli /
ENTRYPOINT ["/alti-cli"]
CMD ["help"]
