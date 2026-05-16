# Step 1: Modules caching
FROM golang:1.26-alpine AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Step 2: Builder
FROM golang:1.26-alpine AS builder

# ビルド時点の最新証明書をインストール
RUN apk add --no-cache ca-certificates
# 【RDS用：必要であれば追加】AWS公式からRDS用ルート証明書をダウンロード
# RUN wget https://amazonaws.com -O /etc/ssl/certs/rds-ca-bundle.pem

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch

COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# 【RDS用：必要であれば追加】RDS用の証明書もコピー
# COPY --from=builder /etc/ssl/certs/rds-ca-bundle.pem /etc/ssl/certs/

ENTRYPOINT ["/app"]
CMD ["--help"]
