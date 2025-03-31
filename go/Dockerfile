# 使用するGoのイメージを指定
FROM golang:1.24-alpine AS builder

# 環境変数を設定
ENV GO111MODULE=on
ENV GOPATH=/go

# 作業ディレクトリを指定
WORKDIR /app

# 必要なGoモジュールをインストール
COPY go.mod go.sum ./
RUN go mod tidy

# 依存関係のパッケージをダウンロード
RUN go get github.com/mattn/go-sqlite3 go.uber.org/mock/gomock

# プロジェクトのソースコードをコピー
# 必要なファイルだけ正しい場所にコピー
COPY go/cmd/api /app/cmd/api
COPY go/app /app/app
COPY go/db /app/db
COPY go/items.json /app/items.json

# ビルド
RUN go build -o main /app/cmd/api/main.go

# 軽量な実行用のイメージを作成
FROM alpine:latest

# 必要なライブラリをインストール
# Dockerfileの途中でインストール結果を確認
RUN apk update --no-cache && apk add --no-cache sqlite && sqlite3 --version


# 作業ディレクトリを /root に設定
WORKDIR /root

# ビルドしたバイナリをコンテナにコピー
COPY --from=builder /app/main .
COPY go/items.json ./items.json
COPY go/db ./db

# ポート9000を開放
EXPOSE 9000

# 実行時に main バイナリを実行
CMD ["./main"]

# コンテナ内での画像保存場所を指定
VOLUME ["/app/images"]

