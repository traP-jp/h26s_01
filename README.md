# 2026年度春ハッカソン01班

タスクランナーとして [`xc`](https://xcfile.dev/) を使っています。使わなくても問題ありません。`Directory` に書かれているディレクトリに移動した後、その下にあるコマンドを実行してください。

インストールは次のコマンドで行えます（Goが必要）。詳しくは [ドキュメント](https://xcfile.dev/getting-started/#installation) を参照してください。

```
go install github.com/joerdav/xc/cmd/xc@latest
```


## Tasks

### setup

クライアントの依存関係をインストールします。

Directory: client

```sh
npm install
```

### dev

クライアントを通常の開発サーバーで起動します。

Directory: client

```sh
npm run dev
```

### dev-mock

mock サーバーとクライアントを同時に起動します。

Directory: client

```sh
npm run dev:mock
```

### check

フォーマット、lint、型チェック、mock サーバーの型チェックを実行します。

Directory: client

```sh
npm run check
```

### build

ビルドが成功するか確認します。

Directory: client

```sh
npm run build
```

### gen-mock-api

OpenAPI から TypeScript の型定義を生成します。

Directory: client

```sh
npm run gen:mock-api
```
