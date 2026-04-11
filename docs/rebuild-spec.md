# pummit 再実装仕様

## 1. 目的

この文書は、`docs/spec.md` に記録された現行挙動を踏まえつつ、`pummit` を一から再実装するための正規仕様を定義する。

この文書は次の 3 点を目的とする。

- プロダクトとして外部に約束する挙動を固定する
- 実装者が迷わないディレクトリ構成と技術選定を与える
- エージェント実装時に曖昧さが出ない受け入れ条件を定義する

競合時の優先順位は次のとおりとする。

1. この文書
2. `docs/config.schema.cue`
3. `docs/agent-handoff.md`
4. `docs/spec.md`
5. README

## 2. プロダクト方針

`pummit` は、Git のステージ済み変更から一貫したコミットメッセージを作る CLI ツールである。

初期リリースの優先順位は次のとおりとする。

1. 明確で壊れにくい CLI 契約
2. 設定ファイルの安全な読み書き
3. オフラインでも破綻しない emoji 解決
4. テストしやすい内部構造
5. MCP からの再利用性

## 3. リリース範囲

### 3.1 初期リリースに含めるもの

- `commit` コマンドと互換ルート呼び出し
- `config` サブコマンド群
- `alias` サブコマンド群
- `doctor`
- `migrate` と `migrate status` と `migrate rollback`
- `mcp`
- 設定ファイルの JSON -> TOML 移行
- `branchMapping` の通常 CLI での利用
- 設定読み込み時の既定値マージ

### 3.2 初期リリースに含めないもの

- `ai` コマンド
- 本格的な対話 TUI
- 複雑なテンプレート DSL
- 履歴学習ベースの scope 推定
- 多言語文言の完全対応

補足:

- `interactive` は将来追加候補だが、初期リリースでは必須にしない
- `templates` と `scope` は設定構造を保持してよいが、未対応機能は明示的に `not implemented` を返す

## 4. CLI 正規仕様

### 4.1 コマンド体系

正規コマンドは次を採用する。

```bash
pummit commit --emoji <emoji-or-alias> <message...>
pummit commit --auto-emoji <message...>
```

互換モードとして、次の呼び出しも維持する。

```bash
pummit <emoji-or-alias> <message...>
```

互換モードは内部的に `pummit commit --emoji <arg0> <rest...>` と同等に扱う。

### 4.2 グローバルフラグ

| フラグ | 意味 |
| --- | --- |
| `-v`, `--version` | バージョン表示 |
| `--offline` | Gitmoji API 通信を無効化 |
| `--config <path>` | 任意の設定ファイルを使う |

### 4.3 `commit` の挙動

`commit` は次の順で処理する。

1. 設定を読み込む
2. 設定ファイルの未知キーと制約違反を検証する
3. ステージ済みファイル一覧を取得する
4. 対象が 0 件なら終了コード `0` で終了し、何もコミットしない
5. emoji 解決を行う
6. コミットメッセージを組み立てる
7. `git commit -m` を実行する

### 4.4 emoji 解決の優先順位

明示 emoji の場合:

1. alias shortcut 一致
2. 組み込み emoji カタログ一致
3. オンライン Gitmoji API 一致
4. `:name:` フォールバック

`--auto-emoji` の場合:

1. `branchMapping.rules` の先頭一致
2. `branchMapping.fallback`
3. それでも解決できなければ `construction`

### 4.5 コミットメッセージ形式

初期リリースでは 1 行メッセージのみを正式サポートする。

```text
<prefix> <message> (<files>)
```

ルール:

- `<prefix>` は raw emoji または `:name:`
- `<message>` は空文字不可
- `<files>` はステージ済みファイルを `, ` で連結したもの
- `base.filesLength = 0` は無制限
- `base.filesLength > 0` の場合は rune 単位で切り詰める

### 4.6 終了コード

| 条件 | 終了コード |
| --- | --- |
| 正常終了 | `0` |
| コミット対象なし | `0` |
| 設定エラー | `2` |
| Git 実行エラー | `3` |
| ユーザー入力エラー | `4` |
| 内部エラー | `1` |

### 4.7 エラー出力ポリシー

- 成功時の主要結果は標準出力
- エラー内容は標準エラー出力
- ライブラリ層では `os.Exit` を使わない
- 終了コードの決定は CLI 層だけで行う

## 5. 設定仕様

### 5.1 保存先

設定ディレクトリは `os.UserConfigDir()` を基準に決定する。

| OS | 保存先 |
| --- | --- |
| Windows | `%AppData%/pummit` |
| macOS | `~/Library/Application Support/pummit` 相当 |
| Linux | `$XDG_CONFIG_HOME/pummit` または `~/.config/pummit` |

### 5.2 読み込みルール

1. 既定値をメモリ上に構築する
2. `config.toml` があればその内容を上書きマージする
3. `config.toml` がなければ `config.json` から移行する
4. どちらもなければ既定の `config.toml` を生成する

### 5.3 検証ルール

起動時と `config validate` で次を検証する。

- 未知キーが存在しないこと
- `meta.version` が空でないこと
- `base.filesLength >= 0`
- `scope.historyLimit >= 0`
- `templates.enabled = true` のとき `defaultTemplate` が存在すること
- `branchMapping.rules[].pattern` が正規表現として妥当であること
- alias shortcut が重複しないこと

### 5.4 設定互換方針

- 読み込みは `config.toml` を唯一の正式形式とする
- `config.json` は移行入力専用とする
- `meta.version` は設定スキーマの識別子でありアプリ本体のバージョンではない

## 6. 推奨ディレクトリ構成

次の構成を推奨する。

```text
main.go

internal/
  app/
    app.go
  cli/
    root.go
    commit.go
    version.go
    doctor.go
    migrate.go
    mcp.go
    alias/
      add.go
      list.go
      delete.go
      reset.go
    config/
      list.go
      get.go
      set.go
      edit.go
      validate.go
      reset.go
  domain/
    config/
      model.go
      defaults.go
      validate.go
    commit/
      model.go
      format.go
    emoji/
      model.go
    alias/
      model.go
  usecase/
    commit_service.go
    config_service.go
    alias_service.go
    doctor_service.go
    migrate_service.go
  infra/
    git/
      client.go
    config/
      store.go
      toml_codec.go
      json_migrator.go
    emoji/
      catalog.go
      gitmoji_client.go
    mcp/
      server.go
      tool_git.go
      tool_config.go
      tool_alias.go
      tool_doctor.go
  ui/
    prompt/
      confirm.go
      editor.go
  testutil/
    gitrepo/
    configfile/

tests/
  e2e/
  integration/
  acceptance/

legacy/
  cli/
  config/
  git/
  alias/
  emojis/
  doctor/
  mcp/
  prompt/
  utils/
  gitmoji/
  logger/

docs/
  spec.md
  rebuild-spec.md
  agent-handoff.md
  config.schema.cue
```

### 6.1 構成方針

- ルート `main.go` は当面維持し、`go install github.com/HidemaruOwO/pummit@latest` の導線を壊さない
- `legacy/` は旧実装の退避場所であり、新機能の実装先ではない
- `domain` には副作用のない型とルールだけを置く
- `usecase` でユースケース単位の処理をまとめる
- `infra` で Git、設定ファイル、外部 API を吸収する
- `cli` はフラグ解析と終了コード変換だけに寄せる
- `ui` は対話コンポーネントだけに限定する
- `pkg` は初期リリースでは作らない

## 7. 推奨ライブラリ選定

| 領域 | 推奨 | 理由 |
| --- | --- | --- |
| CLI | `github.com/spf13/cobra` | 既存資産と親和性が高く、サブコマンド構成が安定している |
| TOML | `github.com/BurntSushi/toml` | Go で安定しており、明示的デコードがしやすい |
| Git 実行 | 標準ライブラリ `os/exec` | 実際の `git` を叩く方が hooks, config, porcelain と整合する |
| ログ | 標準ライブラリ `log/slog` | 依存を増やさず構造化ログを出せる |
| HTTP | 標準ライブラリ `net/http` | Gitmoji API 程度なら十分 |
| MCP | `github.com/mark3labs/mcp-go` | 現行実装と互換があり、型付きで扱える |
| 対話入力 | `github.com/AlecAivazis/survey/v2` | 軽量で CLI 対話に十分。初期リリース向き |
| 差分比較 | `github.com/google/go-cmp/cmp` | 設定や構造体比較のテストが書きやすい |
| CLI 受け入れ試験 | `github.com/rogpeppe/go-internal/testscript` | CLI の stdout, stderr, exit code を仕様化しやすい |
| 機械可読仕様 | `cuelang.org/go` と CUE ファイル | 設定スキーマと制約を人間と機械で共有できる |

### 7.1 採用しない方がよいもの

- `go-git`
  - 実 Git との挙動差が出やすい
- `viper`
  - 設定探索が暗黙的になりやすく、今回の明示仕様と相性が悪い
- 重い DI フレームワーク
  - この規模では過剰

## 8. 実装マイルストーン

### M1. 基盤

- `main.go`
- `internal/app`
- `internal/cli/root.go`
- 終了コードの統一

### M2. 設定

- `config.toml` 読み込み
- 既定値マージ
- 厳格バリデーション
- JSON から TOML への移行

### M3. コミット

- ステージ済みファイル取得
- alias と emoji 解決
- `branchMapping` による `--auto-emoji`
- コミットメッセージ生成

### M4. 管理コマンド

- `config`
- `alias`
- `doctor`
- `migrate`

### M5. MCP

- `git.status`
- `git.commit`
- `git.smart_commit`
- `config.get`
- `alias.list`
- `doctor.check`

## 9. 受け入れ条件

最低限、次を満たすこと。

1. ルート互換呼び出しと `commit` サブコマンドが両方動く
2. 設定ファイルの部分指定で不足項目が既定値で補完される
3. 設定ファイルの未知キーはエラーになる
4. `branchMapping` が通常 CLI でも使われる
5. Git エラー時に正しい終了コードを返す
6. `git commit` 失敗時に内部層から `os.Exit` しない
7. テストで stdout, stderr, exit code を固定できる

## 10. ドキュメント運用

再実装期間中の各文書の役割は次のとおりとする。

- `docs/spec.md`: 現行実装の参照資料
- `docs/rebuild-spec.md`: 再実装の正規仕様
- `docs/agent-handoff.md`: エージェント作業契約
- `docs/config.schema.cue`: 設定スキーマの機械可読定義

この 4 つが揃っていれば、実装者が人でもエージェントでも判断基準を共有しやすい。
