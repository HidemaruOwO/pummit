# pummit 仕様書

> 注記:
> この文書は現行実装の挙動を記録するための `Current Behavior Spec` です。
> 一から再実装する際の正規仕様は `docs/rebuild-spec.md`、エージェント向け実装契約は `docs/agent-handoff.md`、設定スキーマの機械可読定義は `docs/config.schema.cue` を参照してください。

## 1. この文書の目的

この文書は、現行の `pummit` 実装をもとに、次の 2 点を日本語で整理した仕様書です。

- 設定ファイルの形式、保存場所、読み込み順、移行ルール
- コマンドラインインターフェースの引数、フラグ、挙動

この文書は README の紹介文ではなく、実装コードとテストから読み取れる現在の動作を正とします。

## 2. コードベースの全体像

`pummit` は Go 製の CLI ツールです。主要な責務は次のように分かれています。

| パッケージ | 役割 |
| --- | --- |
| `main.go` | CLI 起動 |
| `internal/cli` | Cobra ベースのコマンド定義 |
| `internal/config` | 設定ファイルの保存、読み込み、JSON から TOML への移行 |
| `internal/git` | Git 状態確認、変更ファイル取得、コミット実行 |
| `internal/alias` | 絵文字エイリアスの追加、削除、一覧、初期化 |
| `internal/emojis` | 埋め込み済み絵文字データと Gitmoji API を使った絵文字解決 |
| `pkg/gitmoji` | Gitmoji API 通信とオンライン判定 |
| `internal/doctor` | 診断機能 |
| `internal/mcp` | MCP サーバーと各ツール登録 |
| `internal/prompt` | 確認プロンプトとエディター選択 UI |
| `pkg/logger` | CLI 向けログ出力 |

## 3. アプリケーション情報

- アプリケーションバージョン: `1.3.0`
- 設定ファイルのスキーマ用メタ情報の既定値: `3.0`

この 2 つは別物です。`meta.version` はアプリ本体のバージョンではなく、設定スキーマの識別用に使われています。

## 4. 設定ファイル仕様

### 4.1 保存場所

設定ディレクトリは OS ごとに次の場所になります。

| OS | 保存先 |
| --- | --- |
| Windows | `%APPDATA%\pummit` |
| Windows で `%APPDATA%` が空 | `%USERPROFILE%\.pummit` 相当 |
| macOS / Linux / その他 Unix 系 | `~/.config/pummit` |

補足:

- 現行実装は `XDG_CONFIG_HOME` を見ていません
- 現行のアクティブ設定ファイルは `config.toml` です
- 旧形式の設定ファイルは `config.json` です

### 4.2 起動時の読み込み順

`pummit` 起動時は次の順で設定を決めます。

1. 設定ディレクトリを作成する
2. `config.toml` が存在すれば、それを読み込む
3. `config.toml` がなく `config.json` があれば、自動移行して `config.toml` を作る
4. どちらもなければ、既定値で `config.toml` を新規作成する

優先順位は常に `config.toml` が最優先です。`config.toml` と `config.json` が両方ある場合、通常利用では `config.toml` が使われ、`config.json` は旧形式として残ります。

### 4.3 現行 TOML 構造

設定のトップレベル構造は次のとおりです。

```toml
[meta]
[base]
[interactive]
[locale]
[templates]
[scope]
[alias]
[[alias.entries]]
[branchMapping]
[[branchMapping.rules]]
```

### 4.4 既定の TOML 設定

現行実装が新規作成する既定設定は次の内容です。

```toml
[meta]
version = "3.0"

[base]
emoji = true
filesLength = 50

[interactive]
enabled = true
defaultMode = "emoji"
showPreview = true
fuzzySearch = true

[locale]
language = "ja"
autoDetect = true

[templates]
enabled = true
defaultTemplate = "default"

[templates.definitions.default]
format = "{emoji} {message} ({files})"
description = "Default template"

[templates.definitions.feat]
format = "{emoji} {scope}: {message}\n\n{description}\n\n({files})"
description = "Feature addition template"
defaultEmoji = "sparkles"
scope = true
requireDescription = true

[scope]
enabled = true
autoDetect = true
suggestions = ["api", "ui", "core", "auth", "db"]
fromHistory = true
historyLimit = 50

[alias]
enabled = true

[[alias.entries]]
shortcuts = ["s", "feat", "feature"]
name = "sparkles"
emoji = "✨"

[[alias.entries]]
shortcuts = ["c", "wip"]
name = "construction"
emoji = "🚧"

[[alias.entries]]
shortcuts = ["t", "new", "init"]
name = "tada"
emoji = "🎉"

[[alias.entries]]
shortcuts = ["r", "pr", "pull", "merge"]
name = "recycle"
emoji = "♻️"

[[alias.entries]]
shortcuts = ["wb", "rm", "remove", "del", "delete"]
name = "wastebasket"
emoji = "🗑️"

[[alias.entries]]
shortcuts = ["b", "fix", "error"]
name = "bug"
emoji = "🐛"

[[alias.entries]]
shortcuts = ["e", "lint", "format", "refactor"]
name = "eyes"
emoji = "👀"

[[alias.entries]]
shortcuts = ["d", "doc", "docs", "document", "documents"]
name = "books"
emoji = "📚"

[[alias.entries]]
shortcuts = ["a", "ui", "design", "icon", "icons"]
name = "art"
emoji = "🎨"

[[alias.entries]]
shortcuts = ["h", "tune", "tuning", "perf", "perform", "performance"]
name = "rocket"
emoji = "🚀"

[[alias.entries]]
shortcuts = ["w", "change", "tool", "tools", "lib", "library"]
name = "wrench"
emoji = "🔧"

[[alias.entries]]
shortcuts = ["l", "test", "testing"]
name = "rotating_light"
emoji = "🚨"

[[alias.entries]]
shortcuts = ["p", "pack", "mod", "module"]
name = "package"
emoji = "📦️"

[branchMapping]
enabled = true
fallback = "construction"

[[branchMapping.rules]]
pattern = "^feature/.*"
emoji = "sparkles"
description = "Feature branch"

[[branchMapping.rules]]
pattern = "^(fix|bugfix|hotfix)/.*"
emoji = "bug"
description = "Bug fix branch"

[[branchMapping.rules]]
pattern = "^docs/.*"
emoji = "books"
description = "Documentation branch"

[[branchMapping.rules]]
pattern = "^refactor/.*"
emoji = "eyes"
description = "Refactoring branch"

[[branchMapping.rules]]
pattern = "^test/.*"
emoji = "rotating_light"
description = "Test branch"
```

### 4.5 各設定項目の意味

#### `[meta]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `version` | string | `"3.0"` | 設定スキーマのバージョン識別子 |

#### `[base]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `emoji` | bool | `true` | `true` のとき raw emoji を使う。`false` のとき `:name:` 形式を使う |
| `filesLength` | int | `50` | コミット末尾に付ける変更ファイル一覧の最大文字数。`0` は切り詰めなし |

補足:

- `filesLength > 0` のときだけ切り詰めが動きます
- 負の値は `config validate` では無効です
- README にある `-1 = unlimited` という説明は、現行コードの検証ルールとは一致していません

#### `[interactive]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `enabled` | bool | `true` | 将来向けの対話モード有効フラグ |
| `defaultMode` | string | `"emoji"` | 既定の対話モード名 |
| `showPreview` | bool | `true` | プレビュー表示の既定値 |
| `fuzzySearch` | bool | `true` | あいまい検索の既定値 |

現行コードでは、このセクションを参照して対話コミットを実行する CLI は実装されていません。

#### `[locale]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `language` | string | `"ja"` | 表示言語の設定値 |
| `autoDetect` | bool | `true` | 自動判定を使うかどうか |

現行コードでは、`locale.language` は `mcp` の設定取得結果には出ますが、通常 CLI の表示切り替えには使われていません。

#### `[templates]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `enabled` | bool | `true` | テンプレート機能の有効化フラグ |
| `defaultTemplate` | string | `"default"` | 既定テンプレート名 |
| `definitions` | map | 既定 2 件 | テンプレート定義の集合 |

各テンプレート定義は次のキーを持ちます。

| キー | 型 | 意味 |
| --- | --- | --- |
| `format` | string | フォーマット文字列 |
| `description` | string | テンプレート説明 |
| `defaultEmoji` | string | 既定絵文字名 |
| `scope` | bool | スコープ入力の要否 |
| `requireDescription` | bool | 説明文必須かどうか |

現行コードでは、テンプレート定義は保存と検証はされますが、通常のコミット生成処理では使われていません。

#### `[scope]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `enabled` | bool | `true` | スコープ機能の有効化フラグ |
| `autoDetect` | bool | `true` | 自動推定の有効化フラグ |
| `suggestions` | string 配列 | `api`, `ui`, `core`, `auth`, `db` | 候補一覧 |
| `fromHistory` | bool | `true` | 履歴から候補を作る意図のフラグ |
| `historyLimit` | int | `50` | 履歴参照件数 |

現行コードでは `historyLimit` の妥当性検証はありますが、通常 CLI のコミット処理ではこのセクションは使われていません。

#### `[alias]` と `[[alias.entries]]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `enabled` | bool | `true` | エイリアス解決を使うかどうか |
| `entries` | 配列 | 既定 13 件 | エイリアス定義 |

各 `alias.entries` は次の形です。

| キー | 型 | 意味 |
| --- | --- | --- |
| `shortcuts` | string 配列 | 省略入力に使う別名一覧 |
| `name` | string | 絵文字名。例: `sparkles` |
| `emoji` | string | 実際の絵文字。例: `✨` |

通常コミットで `base.emoji = true` かつ `alias.enabled = true` の場合、入力された先頭引数が `shortcuts` にあれば `emoji` を採用します。

#### `[branchMapping]` と `[[branchMapping.rules]]`

| キー | 型 | 既定値 | 意味 |
| --- | --- | --- | --- |
| `enabled` | bool | `true` | ブランチ名ベース判定の有効化フラグ |
| `fallback` | string | `"construction"` | 規則不一致時の絵文字名 |
| `rules` | 配列 | 既定 5 件 | ブランチ名マッチ規則 |

各 `branchMapping.rules` は次のキーを持ちます。

| キー | 型 | 意味 |
| --- | --- | --- |
| `pattern` | string | 正規表現 |
| `emoji` | string | 絵文字名 |
| `description` | string | 規則の説明 |

重要:

- この設定は現行の通常 CLI コミット処理では参照されていません
- MCP の `git.smart_commit` にはブランチ名から絵文字を決める処理がありますが、こちらは `config.toml` を読まず、コード内に固定された規則を使います

### 4.6 `config validate` の検証ルール

`pummit config validate` が見るルールは次のとおりです。

1. `meta.version` が空でないこと
2. `templates.enabled = true` のとき `templates.defaultTemplate` が空でないこと
3. `templates.defaultTemplate` が `templates.definitions` に存在すること
4. `scope.historyLimit >= 0` であること
5. `base.filesLength >= 0` であること

つまり、すべてのキーの意味的な妥当性を網羅的に検査するわけではありません。たとえば `locale.language` の値が `ja` か `en` かまでは見ていません。

### 4.7 TOML 読み込み時の注意

現行実装には次の特徴があります。

1. TOML の未知キーは読み込み時に無視される
2. 一部キーだけを書いた TOML を読み込んだ場合、不足分は既定値で補完されず、Go のゼロ値になる

例:

- `base.emoji = false` だけを書いた TOML を読むと、`base.filesLength` は既定の `50` ではなく `0` になります

このため、設定ファイルを手書きで最小化すると、意図せず一部機能が無効相当になる可能性があります。

### 4.8 旧 JSON 設定との互換と移行

旧形式の `config.json` は次の構造です。

| キー | 型 | 意味 |
| --- | --- | --- |
| `writeEmoji` | bool | raw emoji を使うか |
| `useAlias` | bool | エイリアスを使うか |
| `useLimitPathesLength` | bool | ファイル一覧の長さ制限を使うか |
| `limitPathesLength` | int | 長さ制限値 |
| `alias` | 配列の配列 | `[["shortcut1,shortcut2", "name", "emoji"]]` のような形式 |

移行時の変換規則は次のとおりです。

| JSON | TOML |
| --- | --- |
| `writeEmoji` | `base.emoji` |
| `useAlias` | `alias.enabled` |
| `useLimitPathesLength = true` と `limitPathesLength` | `base.filesLength` |
| `useLimitPathesLength = false` | `base.filesLength = 0` |
| `alias` | `alias.entries` |

自動移行または `pummit migrate` の通常実行では、元の `config.json` は `config.json.bak` にリネームされます。

## 5. CLI 仕様

### 5.1 ルートコマンド

基本形は次のとおりです。

```bash
pummit [emoji] [message...]
```

グローバルフラグは次の 2 つです。

| フラグ | 意味 |
| --- | --- |
| `-v`, `--version` | バージョン表示 |
| `--offline` | Gitmoji API 呼び出しを無効にし、オフラインモードで動かす |

### 5.2 通常コミットの動作

`pummit [emoji] [message...]` の挙動は次のとおりです。

1. 引数が 2 個未満ならヘルプを表示して終了する
2. `--version` があれば `pummit v1.3.0` を出して終了する
3. `--offline` があれば明示的にオフラインモードにする
4. `--offline` がなくても、簡易オンライン判定に失敗した場合は自動でオフラインモードに切り替える
5. `git diff --name-only --cached` でステージ済みファイル一覧を取得する
6. ステージ済みファイルが 1 件もなければ、情報を出して終了する
7. 先頭引数を絵文字またはエイリアスとして解決する
8. コミットメッセージ本体を組み立てて `git commit -m` を実行する

生成されるコミットメッセージの形式は固定です。

```text
<絵文字または:name:> <message> (<変更ファイル一覧>)
```

例:

```text
✨ Add config docs (docs/spec.md)
```

### 5.3 絵文字解決ルール

`base.emoji` と `alias.enabled` の組み合わせで結果が変わります。

| 条件 | 挙動 |
| --- | --- |
| `base.emoji = true` かつエイリアス一致 | `alias.entries[].emoji` を使う |
| `base.emoji = true` かつエイリアス不一致 | 絵文字名を実絵文字に変換する。失敗時は `:name:` に戻す |
| `base.emoji = false` かつエイリアス一致 | `:alias.entries[].name:` を使う |
| `base.emoji = false` かつエイリアス不一致 | `:入力値:` をそのまま使う |

補足:

- 既知の絵文字名は埋め込み JSON から解決されます
- オフラインで未知の名前だった場合は `:unknown:` のような形式にフォールバックします

### 5.4 変更ファイル一覧の扱い

- 対象はステージ済みファイルだけです
- 区切り文字は `, ` です
- `base.filesLength > 0` のときだけ、文字数上限を超えた部分を `...` で切り詰めます
- 切り詰めは rune 単位なので、日本語ファイル名でも途中破損しにくい実装です

### 5.5 エラーと終了コードの実装上の注意

現行コードには次の特徴があります。

1. ステージ済みファイルがない場合は成功扱いで終了します
2. `git commit` が失敗した場合、エラーメッセージは出ますが、ルートコマンドが `RunE` ではないため終了コードが 1 にならない可能性があります

この文書は現状の実装を説明するため、この挙動もそのまま記載しています。

## 6. サブコマンド仕様

### 6.1 `version`

```bash
pummit version
pummit --version
pummit -v
```

挙動:

- `pummit v1.3.0` を表示する

### 6.2 `config`

```bash
pummit config <subcommand>
```

利用できるサブコマンド:

- `list`
- `get [key]`
- `set [key] [value]`
- `edit`
- `validate`
- `reset`

#### `config list`

```bash
pummit config list
```

挙動:

- 現在の `CurrentTOMLConfig` を TOML 形式で標準出力に出す

#### `config get [key]`

```bash
pummit config get base.emoji
pummit config get alias.entries.0.emoji
```

キー指定ルール:

- `.` 区切りで下位キーに降りる
- 構造体は TOML タグ名で指定する
- 配列は数値インデックスで指定する
- map は文字列キーで指定する

出力ルール:

- 文字列、真偽値、整数などの単純値はその値だけを出す
- 構造体、配列、map など複合値は TOML 形式で出す

#### `config set [key] [value]`

```bash
pummit config set base.emoji false
pummit config set base.filesLength 80
pummit config set alias.entries.0.shortcuts s,feat,feature
```

挙動:

1. 指定キーに値を代入する
2. `ValidateRequiredFields` で最小限の検証を行う
3. `config.toml` に保存する

代入ルール:

- string: そのまま保存
- bool: `true` / `false`
- int 系: 10 進整数
- float 系: 小数
- string 配列: カンマ区切り文字列を分割して保存

制約:

- 存在しないキーはエラー
- 配列インデックス範囲外はエラー
- string 配列以外の複合型には直接代入できない

#### `config edit`

```bash
pummit config edit
```

挙動:

1. `EDITOR` 環境変数があれば、そのコマンドで `config.toml` を開く
2. `EDITOR` がなければ、TUI のエディター選択画面を出す

選択画面の特徴:

- 利用可能なコマンドだけを表示する
- 候補には `vim`, `nvim`, `nano`, `emacs`, `code`, `gedit`, `kate`, `vi`, `cat`, `bat`, `gat` が含まれる
- `c` でカスタムコマンド入力
- `q` または `Ctrl+C` で終了

#### `config validate`

```bash
pummit config validate
```

挙動:

- `config.toml` を読み直して TOML として解釈できるか確認する
- その後、4.6 の検証ルールを適用する
- 問題なければ `✅ Configuration is valid` を出す

#### `config reset`

```bash
pummit config reset
pummit config reset --force
```

挙動:

- `--force` がなければ確認プロンプトを出す
- 承認後、設定全体を既定の TOML 設定に戻して保存する

### 6.3 `alias`

```bash
pummit alias <subcommand>
```

利用できるサブコマンド:

- `add [name] [prefix]`
- `list`
- `delete [name]`
- `reset`

補足:

- コマンド定義上の使用法は `add [name] [prefix]` ですが、実際の意味は「第 1 引数が shortcut、第 2 引数が絵文字名」です

#### `alias add [name] [prefix]`

```bash
pummit alias add rs rocket
pummit alias add rs rocket --emoji 🚀
```

実際の意味:

- 第 1 引数: 新しく追加する shortcut
- 第 2 引数: 絵文字名。例: `rocket`
- `--emoji` を付けた場合: 実絵文字を直接指定する

挙動:

1. `--emoji` があれば、その絵文字を使う
2. `--emoji` がなければ第 2 引数の名前から絵文字を推定する
3. 同じ shortcut がすでにあれば失敗する
4. 同じ絵文字が既存エントリにあり、かつ `name` も同じなら、shortcut だけを追加する
5. 同じ絵文字でも `name` が異なる場合は失敗する

#### `alias list`

```bash
pummit alias list
```

挙動:

- `alias.entries` を名前順で表示する
- 絵文字、名前、shortcut 一覧を表形式で出す

#### `alias delete [name]`

```bash
pummit alias delete feat
pummit alias delete feat --confirm
```

実際の挙動:

- 指定するのは shortcut 名です
- その shortcut を `alias.entries` から取り除く
- その entry に shortcut が 1 つも残らなくなったら entry 自体を削除する
- `--confirm` がなければ確認プロンプトを出す

#### `alias reset`

```bash
pummit alias reset
pummit alias reset --confirm
```

挙動:

- alias 設定だけを既定値に戻す
- 他の設定セクションは変更しない

### 6.4 `migrate`

```bash
pummit migrate
pummit migrate --force
pummit migrate --dry-run
pummit migrate --preview
```

フラグ:

| フラグ | 意味 |
| --- | --- |
| `--force` | 既存の `config.toml` があっても上書きする |
| `--dry-run` | 実際には書き込まず、何をするかだけ表示する |
| `--preview` | 結果の `config.toml` をエディターで開く |

通常動作:

1. `config.toml` があり `--force` がなければ何もしない
2. `config.json` がなければ既定の `config.toml` を作る
3. `config.json` があれば TOML に変換して保存する
4. 変換元の `config.json` は `config.json.bak` にリネームする

補足:

- `--preview` がない場合でも、通常実行後は対話プロンプトでプレビューするか確認されます
- `--dry-run` 時は実ファイルを書き換えませんが、既存の `config.toml` があればその内容をプレビューできます

#### `migrate status`

```bash
pummit migrate status
```

挙動:

- 設定ディレクトリ内に `config.toml` と `config.json` があるかを調べ、状態を表示する

状態の種類:

- `both`
- `toml_only`
- `json_only`
- `none`

#### `migrate rollback [backup-file]`

```bash
pummit migrate rollback ~/.config/pummit/config.json.bak
```

挙動:

1. 指定バックアップを `config.json` として復元する
2. `config.toml` があれば削除する

### 6.5 `doctor`

```bash
pummit doctor
```

診断対象:

1. Git グローバル設定の `user.name`, `user.email`
2. 現在ディレクトリが Git リポジトリかどうか
3. ステージ済みファイル、未ステージファイルの状態
4. `config.toml` または `config.json` の読み込み可否
5. Gitmoji API への接続可否
6. 設定ディレクトリへの書き込み可否
7. OS、Go、Git、pummit のバージョン情報

終了条件:

- `ERROR` が 1 件でもあれば終了コード 1
- `WARNING` だけなら表示して継続終了

### 6.6 `mcp`

```bash
pummit mcp
```

挙動:

- stdio 上で MCP サーバーを起動する
- `PUMMIT_MCP_DEBUG=1` が設定されていればデバッグログを有効にする

登録される MCP ツールは次のとおりです。

| ツール名 | 役割 |
| --- | --- |
| `git.smart_commit` | 変更ファイルを見てコミット案を作る。必要なら自動ステージも行う |
| `git.status` | ブランチ名、変更ファイル、ステージ済みファイルを返す |
| `git.commit` | 指定 emoji と message で直接コミットする |
| `git.add_files` | ファイルをステージする |
| `git.get_edited_files` | 変更ファイル一覧を返す |
| `git.get_current_branch` | 現在ブランチを返す |
| `alias.list` | エイリアス一覧を返す |
| `config.get` | 一部設定値または設定要約を返す |
| `doctor.check` | 診断結果を返す |

補足:

- `git.smart_commit` はブランチ名から絵文字候補を推定しますが、その規則は `branchMapping` 設定ではなくコード内固定です
- `config.get` で個別取得できるキーは `base.emoji`, `base.filesLength`, `alias.enabled`, `interactive.enabled`, `locale.language` に限られます

## 7. 現行実装と README の差分

現行コードを読む限り、次の点は README の説明より実装が優先されます。

1. `interactive` コマンドは実装されていない
2. `ai` コマンドは実装されていない
3. `branchMapping` 設定は通常 CLI コミットでは使われていない
4. `templates` と `scope` は設定構造と検証はあるが、通常コミット生成では未使用
5. `locale` は保持されるが、通常 CLI の言語切り替えに広く使われていない
6. `XDG_CONFIG_HOME` は保存先決定に使われていない
7. `base.filesLength` は `0` で切り詰めなし。負数はバリデーションで不正

## 8. 実務上の読み方

現時点の `pummit` は、設定スキーマは将来拡張を見据えて広めに用意されていますが、通常 CLI として実際に強く効いているのは主に次の項目です。

- `base.emoji`
- `base.filesLength`
- `alias.enabled`
- `alias.entries`

そのほかの設定は、現段階では保存、表示、検証、移行の対象ではあっても、通常コミット処理そのものには直結していない項目が多いです。
