# 整理番号確認システム (tickets-system-v1)

店舗やイベントブース等で、現在呼び出し中の整理番号をリアルタイムに表示するWebシステムです。<br />
元々サークルの文化祭で使っていたものを公開用に書き直したものです。<br />
https://blog.kinokoroom.net/%e6%a8%a1%e6%93%ac%e5%ba%97%e3%81%ae%e3%81%9f%e3%82%81%e3%81%ab%e6%95%b4%e7%90%86%e7%95%aa%e5%8f%b7%e8%a1%a8%e7%a4%ba%e3%82%b7%e3%82%b9%e3%83%86%e3%83%a0%e3%82%92%e4%bd%9c%e3%81%a3%e3%81%9f%e8%a9%b1/

## 技術スタック

Go言語をバックエンドとして、HTML, CSS, JavaScriptをフロントエンドとして使用しています。データはJSON形式で保存されます。

## 起動方法（ローカル環境）

1. リポジトリをクローンします。
2. ルートディレクトリに `.env` ファイルを作成し、管理者のパスワードを設定します（未設定の場合は `password` になります）。
   ```env
   ADMIN_PASSWORD="your_secure_password"
   PORT="3000"
   ```
3. Go環境がインストールされていることを確認し、以下のコマンドで起動します。
   ```bash
   go run ./cmd/server
   ```
4. ブラウザでアクセスします。
   - 確認ページ: `http://localhost:3000/`
   - 管理者用ページ: `http://localhost:3000/admin/login`

## 起動方法（Docker）

Docker環境をお持ちの場合は、以下のコマンドで一発で起動できます。

```bash
docker-compose up -d --build
```

## カテゴリの設定方法

ルートディレクトリにある `config.json` を書き換えることで、扱う商品やブースごとの番号帯を変更できます。

```json
{
  "categories": [
    {
      "id": "cotton",
      "name": "綿あめ",
      "rangeStart": 1000,
      "rangeEnd": 1999,
      "color": "#e91e63"
    },
    {
      "id": "tapioca",
      "name": "タピオカ",
      "rangeStart": 2000,
      "rangeEnd": 2999,
      "color": "#9c27b0"
    }
  ]
}
```
