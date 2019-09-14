FROM golang:1.13 as builder
WORKDIR /code
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn \ 
    && go env -w GOSUMDB=sum.golang.google.cn \ 
    && go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o v2ray-exporter .

FROM alpine:latest
COPY --from=builder /code/v2ray-exporter /usr/bin/v2ray-exporter
ENTRYPOINT [ "/usr/bin/v2ray-exporter" ]
CMD [ "--help" ]