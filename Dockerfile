FROM golang:alpine as build
RUN apk --no-cache add ca-certificates

FROM scratch
COPY alti-cli /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/alti-cli"]
CMD ["help"]
