# pummit エージェント実装契約

## 1. 目的

この文書は、実装エージェントに `pummit` の再実装を依頼する際の作業契約を定義する。

この文書の狙いは次のとおり。

- エージェントが勝手に仕様を広げないようにする
- 1 タスクごとの完了条件を固定する
- PR 単位で安全に進められるようにする

## 2. 読む順番

エージェントは必ず次の順で文書を読むこと。

1. `docs/rebuild-spec.md`
2. `docs/config.schema.cue`
3. `docs/spec.md`
4. `legacy/` 配下の既存コード

## 3. 非交渉ルール

エージェントは次を守ること。

- 仕様に書かれていない新機能を追加しない
- 既存コマンド互換を壊す変更は、互換レイヤーなしでは入れない
- ライブラリ層で `os.Exit` を使わない
- 不明点は README ではなく `docs/rebuild-spec.md` を優先する
- `legacy/` は参照用であり、新規実装の追加先にしない
- 1 PR で 1 マイルストーン、または 1 マイルストーン内の 1 サブ機能だけを扱う
- すべての非自明ロジックにユニットテストを追加する
- CLI 契約変更には受け入れテストを追加する

## 4. 実装順序

### Phase 1. 基盤

作るもの:

- `main.go`
- `internal/app/app.go`
- `internal/cli/root.go`
- 終了コード定義

完了条件:

- `pummit --version` が動く
- `pummit help` が動く
- 主要エラーが終了コードに変換される

### Phase 2. 設定

作るもの:

- 設定モデル
- TOML codec
- JSON migrator
- validation

完了条件:

- 設定の部分指定で既定値が補完される
- 未知キーで失敗する
- `config validate` が制約違反を検出する

### Phase 3. コミット

作るもの:

- Git client
- commit service
- emoji resolver
- alias service

完了条件:

- `pummit commit --emoji sparkles "Add x"` が成功する
- `pummit sparkles "Add x"` が互換動作する
- `--auto-emoji` が `branchMapping` を使う
- ステージ済み変更がなければ終了コード `0` で終了する

### Phase 4. 管理コマンド

作るもの:

- `config` サブコマンド群
- `alias` サブコマンド群
- `doctor`
- `migrate`

完了条件:

- 主要サブコマンドにヘルプがある
- 破壊的操作に確認または `--force` がある
- `migrate rollback` が `config.toml` を削除し `config.json` を戻す

### Phase 5. MCP

作るもの:

- `mcp` サーバー
- Git 系ツール
- Config 系ツール
- Alias 系ツール
- Doctor 系ツール

完了条件:

- `git.smart_commit` が staged / unstaged を扱える
- `config.get` は正規仕様のキーだけ返す
- MCP 側も通常 CLI と同じ emoji 解決ルールを使う

## 5. 変更単位

1 回の変更で扱ってよい対象は次のいずれかだけとする。

1. 1 つのコマンド
2. 1 つのユースケース
3. 1 つの設定レイヤー改善
4. 1 つのテスト基盤改善

巨大な横断変更は禁止する。

## 6. 実装時のチェックリスト

各タスクで必ず確認すること。

1. 仕様差分を 3 行以内で要約したか
2. 変更対象レイヤーを 1 つに絞ったか
3. ユニットテストを追加したか
4. CLI 契約が変わるなら受け入れテストを追加したか
5. `gofmt -w .` を実行したか
6. `go test -v ./...` を実行したか
7. `go build . -o pummit` または新しいエントリポイント相当の build を実行したか

## 7. テスト戦略

### 7.1 ユニットテスト対象

- 設定マージ
- 設定バリデーション
- alias 解決
- branchMapping 解決
- コミットメッセージ整形
- ファイル一覧切り詰め

### 7.2 統合テスト対象

- Git リポジトリ上での staged files 取得
- JSON から TOML への移行
- doctor の診断結果

### 7.3 受け入れテスト対象

- `--version`
- `commit --emoji`
- 互換ルート呼び出し
- `--auto-emoji`
- `config validate`
- `alias add/list/delete`

## 8. エージェントへの依頼テンプレート

次の形式で依頼するとよい。

```text
Task: Implement Phase 2 - config loading and validation.

Read first:
1. docs/rebuild-spec.md
2. docs/config.schema.cue
3. docs/agent-handoff.md

Scope:
- Implement only config loading, default merging, strict validation, and config validate command.
- Do not implement commit, MCP, or interactive flows.

Required outputs:
- Code changes
- Unit tests
- Acceptance test updates if CLI behavior changes
- Short summary of what was implemented

Verification:
- gofmt -w .
- go test -v ./...
- go build . -o pummit
```

## 9. 実装完了の定義

再実装が完了したとみなしてよい条件は次のとおり。

1. `docs/rebuild-spec.md` の初期リリース範囲がすべて実装済み
2. `docs/config.schema.cue` と実装の差分がない
3. 主要コマンドの受け入れテストがある
4. 既存互換呼び出しが維持されている
5. README が正規仕様と一致している
