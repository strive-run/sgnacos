FROM golang:latest as builder

WORKDIR /app
COPY . .
ENV GOPROXY="https://goproxy.cn"
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o nacos-server .


FROM scratch as runtime
COPY --from=builder /app/nacos-server /usr/local/bin/nacos-server
COPY --from=builder /app/config /etc/mock/server/data/
COPY --from=builder /app/conf/config.yaml /etc/mock/server/config.yaml

ENV SERVER_CONFIG "/etc/mock/server/config.yaml"
ENTRYPOINT ["/usr/local/bin/nacos-server"]
EXPOSE 8848
