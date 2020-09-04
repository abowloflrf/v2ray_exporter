FROM golang:1.15 as builder
WORKDIR /code
COPY . .
RUN GOPRIVATE=github.com/v2ray/v2ray-core CGO_ENABLED=0 go build -o v2ray_exporter .

FROM alpine:3.12
COPY --from=builder /code/v2ray_exporter /usr/bin/v2ray_exporter
ENTRYPOINT [ "/usr/bin/v2ray_exporter" ]
CMD [ "--help" ]