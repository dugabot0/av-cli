# av-cli

DMM アフィリエイト API v3 と DUGA API をラップする CLI ツール。
検索結果を JSON で標準出力するため、AI エージェントとの連携に最適です。

## 機能

- **DMM**: フロア一覧 / 商品検索 / 女優検索 / ジャンル・メーカー・シリーズ・作者検索
- **DUGA**: 商品検索
- 出力は常に JSON（`--pretty` で整形、デフォルトはコンパクト）
- エラーメッセージは stderr、JSON は stdout に分離
- 明確な終了コードでスクリプト・エージェントから扱いやすい設計

## API 登録

各サービスのアフィリエイト会員登録と API 申請が必要です。

- **DMM**: [https://affiliate.dmm.com/](https://affiliate.dmm.com/) から会員登録・API ID を取得してください。
- **DUGA**: [https://click.duga.jp/aff/40413-01](https://click.duga.jp/aff/40413-01) から登録してください。このリンク経由で登録すると、あなたがアフィリエイト収益を得た際に作者への紹介報酬が発生する場合があります。

## インストール

### バイナリをダウンロード（推奨）

[GitHub Releases](https://github.com/dugabot0/av-cli/releases) から OS・アーキテクチャに合ったバイナリをダウンロードしてください。

```bash
# 例: Linux amd64
curl -L https://github.com/dugabot0/av-cli/releases/latest/download/av-cli_linux_amd64.tar.gz | tar xz
sudo mv av-cli /usr/local/bin/
```

### Homebrew

```bash
brew install dugabot0/tap/av-cli
```

### ソースからビルド

```bash
git clone https://github.com/dugabot0/av-cli.git
cd av-cli
go build -o av-cli .
```

## 認証情報の設定

**優先順位:** CLI フラグ > 環境変数 > 設定ファイル

### 設定ファイル

`~/.config/av-cli/config.yaml` を作成してください。

```yaml
dmm:
  api_id: "your-api-id"
  affiliate_id: "your-affiliate-990"
duga:
  app_id: "your-duga-app-id"
  agent_id: "your-agent-id"
  banner_id: "01"
```

### 環境変数

```bash
export DMM_API_ID=your-api-id
export DMM_AFFILIATE_ID=your-affiliate-990
export DUGA_APP_ID=your-duga-app-id
export DUGA_AGENT_ID=your-agent-id
export DUGA_BANNER_ID=01
```

## 使い方

### グローバルフラグ

```
--pretty             JSON を整形して出力（デフォルト: コンパクト）
--quiet              stderr へのログ出力を抑制
--timeout duration   HTTP タイムアウト（デフォルト: 30s）
```

---

### `av-cli dmm floors`

フロア一覧を取得します。他の `dmm` コマンドで必要な `floor_id` はここで確認できます。

```bash
av-cli dmm floors --pretty
```

<details>
<summary>レスポンス例</summary>

```json
{
  "site": [
    {
      "name": "FANZA（アダルト）",
      "code": "FANZA",
      "service": [
        {
          "name": "動画",
          "code": "digital",
          "floor": [
            { "id": "43", "name": "ビデオ", "code": "videoa" },
            { "id": "44", "name": "素人",   "code": "videoc" }
          ]
        }
      ]
    }
  ]
}
```
</details>

---

### `av-cli dmm items`

商品を検索します。

```bash
av-cli dmm items --keyword "松本いちか" --hits 5 --pretty
av-cli dmm items --site FANZA --service digital --floor videoa --sort date --hits 10 --pretty
av-cli dmm items --article actress --article-id 1054998 --pretty
```

| フラグ | 説明 | デフォルト |
|--------|------|-----------|
| `--keyword` | 検索キーワード | |
| `--site` | `FANZA` または `DMM.com` | `FANZA` |
| `--service` | サービスコード（`dmm floors` で確認） | |
| `--floor` | フロアコード（`dmm floors` で確認） | |
| `--hits` | 取得件数（1〜100） | `20` |
| `--offset` | 開始位置 | `1` |
| `--sort` | `rank` / `price` / `-price` / `date` / `review` / `match` | |
| `--article` | 絞り込み種別: `actress` / `author` / `genre` / `series` / `maker`（複数指定可） | |
| `--article-id` | 絞り込み ID（複数指定可） | |
| `--cid` | コンテンツ ID で直接取得 | |

---

### `av-cli dmm actress`

女優を検索します。

```bash
av-cli dmm actress --keyword "松本いちか" --pretty
av-cli dmm actress --gte-bust 90 --lte-waist 60 --sort -bust --hits 10 --pretty
```

| フラグ | 説明 |
|--------|------|
| `--keyword` | 名前キーワード |
| `--actress-id` | 女優 ID で直接取得 |
| `--initial` | 頭文字（あ〜ん） |
| `--gte-bust` / `--lte-bust` | バスト下限 / 上限（cm） |
| `--gte-waist` / `--lte-waist` | ウエスト下限 / 上限（cm） |
| `--gte-hip` / `--lte-hip` | ヒップ下限 / 上限（cm） |
| `--gte-height` / `--lte-height` | 身長下限 / 上限（cm） |
| `--gte-birthday` / `--lte-birthday` | 生年月日（`yyyy-mm-dd`） |
| `--hits` | 取得件数（1〜100）、デフォルト `20` |
| `--offset` | 開始位置 |
| `--sort` | `name` / `-name` / `bust` / `-bust` / `waist` / `-waist` / `hip` / `-hip` / `height` / `-height` / `birthday` / `-birthday` / `id` / `-id` |

---

### `av-cli dmm genre` / `maker` / `series` / `author`

ジャンル・メーカー・シリーズ・作者を検索します。`--floor-id` は必須です（`dmm floors` で確認）。

```bash
av-cli dmm genre  --floor-id 43 --pretty
av-cli dmm maker  --floor-id 43 --initial あ --pretty
av-cli dmm series --floor-id 43 --hits 50 --pretty
av-cli dmm author --floor-id 82 --pretty
```

| フラグ | 説明 | デフォルト |
|--------|------|-----------|
| `--floor-id` | フロア ID（**必須**） | |
| `--initial` | 頭文字フィルター（あ〜ん） | |
| `--hits` | 取得件数（1〜500） | `100` |
| `--offset` | 開始位置 | `1` |

---

### `av-cli duga search`

DUGA の商品を検索します。

```bash
av-cli duga search --keyword "素人" --hits 5 --pretty
av-cli duga search --sort rating --adult 1 --hits 20 --pretty
av-cli duga search --release-start 20240101 --release-end 20241231 --pretty
```

| フラグ | 説明 | デフォルト |
|--------|------|-----------|
| `--keyword` | 検索キーワード | |
| `--hits` | 取得件数（1〜100） | `10` |
| `--offset` | 開始位置 | `1` |
| `--sort` | `favorite` / `release` / `new` / `price` / `rating` / `mylist` | |
| `--category` | カテゴリ ID | |
| `--performer-id` | 出演者 ID で絞り込み | |
| `--series-id` | シリーズ ID で絞り込み | |
| `--label-id` | レーベル ID で絞り込み | |
| `--adult` | `1`=アダルト / `0`=一般 | `1` |
| `--target` | `ppv` / `sd` / `rental` / `hd` / `hdrental` | |
| `--open-start` / `--open-end` | 公開日範囲（`YYYYMMDD`） | |
| `--release-start` / `--release-end` | 発売日範囲（`YYYYMMDD`） | |

---

## AI エージェントとの連携

stdout は常に JSON のため、`jq` やエージェントのツール呼び出し結果としてそのまま利用できます。

```bash
# jq でパイプ処理
av-cli dmm actress --keyword "松本いちか" | jq '.[0].id'

# エージェント向け: エラー時は stderr にメッセージ、終了コードで判別
av-cli dmm items --keyword "test" --quiet
echo $?  # 0=成功 2=入力エラー 3=認証エラー 4=ネットワークエラー
```

### 終了コード

| コード | 意味 |
|--------|------|
| `0` | 成功 |
| `1` | 予期しないエラー |
| `2` | 入力エラー（必須フラグ未指定・認証情報未設定） |
| `3` | 認証エラー（API が 401/403 を返した） |
| `4` | ネットワーク／API エラー（タイムアウト・5xx） |

## ライセンス

MIT
