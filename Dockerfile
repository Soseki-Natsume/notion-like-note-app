# ---------------------------------------------------
# Stage 1: TypeScript のビルド
# ---------------------------------------------------
FROM node:20-slim AS ts-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# ---------------------------------------------------
# Stage 2: Go アプリのビルド (CGO有効化: SQLite依存のため)
# ---------------------------------------------------
FROM golang:latest AS go-builder
WORKDIR /app

# CGOに必要なビルドツールのインストール
RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Stage 1 でビルドされた JS/CSS などの静的ファイルをコピー
COPY --from=ts-builder /app/dist ./static

# CGO_ENABLED=1 でビルド
RUN CGO_ENABLED=1 GOOS=linux go build -o app .

# ---------------------------------------------------
# Stage 3: 実行用軽量イメージ
# ---------------------------------------------------
FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

# ビルド成果物と静的ファイルをコピー
COPY --from=go-builder /app/app .
COPY --from=go-builder /app/dist ./static

# SQLite データを保持するボリュームマウント用ディレクトリを作成
RUN mkdir -p /data

EXPOSE 8080

CMD ["./app"]