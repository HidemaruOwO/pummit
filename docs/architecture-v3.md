# Pummit アーキテクチャ仕様書

## 目次

1. [プロジェクト概要](#プロジェクト概要)
2. [システム設計](#システム設計)
3. [モジュール構成](#モジュール構成)
4. [データフロー](#データフロー)
5. [設定システム](#設定システム)
6. [外部依存関係](#外部依存関係)
7. [新機能仕様](#新機能仕様)
8. [MCPサーバー機能](#mcpサーバー機能)
9. [品質保証・テスト戦略](#品質保証・テスト戦略)
10. [セキュリティ考慮事項](#セキュリティ考慮事項)
11. [パフォーマンス設計](#パフォーマンス設計)
12. [リスク分析と対策](#リスク分析と対策)
13. [今後の開発指針](#今後の開発指針)

---

## プロジェクト概要

### 目的

**Pummit**は、Gitコミットメッセージを絵文字付きで美しく作成するためのCLIツールです。開発者の生産性向上と、一貫性のあるコミットメッセージ作成を支援します。

### 主要機能

- **絵文字付きコミット**: 絵文字プレフィックスによる視覚的なコミットメッセージ
- **エイリアス機能**: 短縮形での絵文字指定
- **自動ファイル検出**: 変更されたファイルの自動追加
- **設定管理**: ユーザーカスタマイズ可能な設定
- **Gitmoji連携**: 外部Gitmojiサービスとの統合
- **インタラクティブUI**: 確認ダイアログとプロンプト

### 技術要件

- **言語**: Go 1.20+
- **プラットフォーム**: クロスプラットフォーム対応（Windows, macOS, Linux）
- **依存関係**: 最小限の外部依存
- **パフォーマンス**: 高速な起動とレスポンス

---

## システム設計

### 全体アーキテクチャ

```mermaid
graph TD
    subgraph User Interface Layer
        A[CLI Commands]
        B[Interactive Prompts]
        C[MCP Server (pummit mcp)]
    end

    subgraph Application Layer
        D[Root Command Handler]
        E[Alias Commands]
        F[MCP Command Handler]
    end

    subgraph Business Logic Layer
        G[Git Ops]
        H[Alias Mgmt]
        I[Config Mgmt]
        J[Doctor]
    end

    subgraph Infrastructure Layer
        K[File System]
        L[Network]
        M[Logging]
        N[Stdio JSON-RPC]
    end

    subgraph External Systems
        O[Git Repository]
        P[Gitmoji Service]
        Q[User Config Files]
        R[MCP Client (e.g. Claude Desktop)]
    end

    A & B --> D & E
    C --> F

    D & E --> G & H & I & J
    F --> G & H & I & J

    G & H & I & J --> K & L & M
    F --> N

    K & L & M --> O & P & Q
    N <--> R
```

### レイヤー構成

#### 1. ユーザーインターフェース層

**責務**: ユーザーとの直接的なやり取り
- CLI コマンド解析
- インタラクティブプロンプト
- 出力フォーマッティング

**主要コンポーネント**:
- [`internal/cli/root.go`](internal/cli/root.go:1) - メインコマンド
- [`internal/cli/alias/`](internal/cli/alias/:1) - エイリアス管理コマンド
- [`internal/prompt/prompt.go`](internal/prompt/prompt.go:1) - インタラクティブUI

#### 2. アプリケーション層

**責務**: ビジネスロジックの調整と制御
- コマンド処理の協調
- エラーハンドリング
- ログ管理

#### 3. ビジネスロジック層

**責務**: 核となる業務ロジック
- Git操作の抽象化
- エイリアス管理
- 絵文字変換処理
- 設定管理

#### 4. インフラストラクチャ層

**責務**: 外部システムとの連携
- ファイルシステム操作
- ネットワーク通信
- ログ出力

---

## モジュール構成

### パッケージ構造

```
pummit/
├── main.go                           # エントリーポイント
├── internal/                         # 内部実装（非公開）
│   ├── cli/                         # CLI関連
│   │   ├── root.go
│   │   ├── mcp.go                   # MCPサーバー起動コマンド
│   │   └── ...
│   ├── mcp/                         # MCPサーバー機能
│   │   ├── server.go                # サーバー初期化とツール登録
│   │   ├── tool_git.go              # Git関連ツール
│   │   ├── tool_config.go           # 設定関連ツール
│   │   ├── tool_alias.go            # エイリアス関連ツール
│   │   └── tool_doctor.go           # 診断ツール
│   ├── alias/                       # エイリアス管理ロジック
│   │   └── alias.go
│   ├── config/                      # 設定管理
│   │   └── config.go
│   ├── doctor/                      # 診断機能
│   │   └── checker.go
│   ├── emojis/                      # 絵文字処理
│   │   └── emojis.go
│   ├── git/                         # Git操作
│   │   └── git.go
│   ├── prompt/                      # インタラクティブUI
│   │   └── prompt.go
│   ├── utils/                       # ユーティリティ
│   │   └── slice.go
│   └── variable/                    # 定数・埋め込みデータ
│       ├── consts.go
│       ├── config.json
│       └── discord-emojis.flat.json
├── pkg/                             # 公開パッケージ
│   ├── gitmoji/                     # Gitmoji API連携
│   │   └── gitmoji.go               # API クライアント
│   └── logger/                      # ログ機能
│       └── logger.go                # ロガー実装
└── docs/                            # ドキュメント
    └── ...
```

### 責務分散

#### Core Domain (`internal/`)

**Git Operations** ([`internal/git/git.go`](internal/git/git.go:1))
```go
type CommitMessage struct {
    Emoji   string
    Message string
}

// 主要な責務:
// - Gitリポジトリ操作
// - コミットメッセージ生成
// - 変更ファイル検出
// - ブランチ情報取得
```

**Alias Management** ([`internal/alias/alias.go`](internal/alias/alias.go:1))
```go
type Alias struct {
    Name   string
    Prefix string
    Emoji  string
}

// 主要な責務:
// - エイリアス CRUD 操作
// - 複数エイリアス名の管理
// - エラーハンドリング
```

**Configuration** ([`internal/config/config.go`](internal/config/config.go:1))
```go
type Config struct {
    UseRawEmoji    bool       `json:"writeEmoji"`
    UseAlias       bool       `json:"useAlias"`
    UseFilesLength bool       `json:"useLimitPathesLength"`
    FilesLength    int        `json:"limitPathesLength"`
    Aliases        [][]string `json:"alias"`
}

// 主要な責務:
// - 設定ファイルの読み書き
// - デフォルト設定の管理
// - 設定検証
```

#### Infrastructure (`pkg/`)

**Gitmoji Integration** ([`pkg/gitmoji/gitmoji.go`](pkg/gitmoji/gitmoji.go:1))
```go
type Gitmoji struct {
    Emoji       string `json:"emoji"`
    Entity      string `json:"entity"`
    Code        string `json:"code"`
    Description string `json:"description"`
    Name        string `json:"name"`
}

// 主要な責務:
// - 外部Gitmoji API との連携
// - HTTP クライアント管理
// - レスポンス処理
```

**Logging** ([`pkg/logger/logger.go`](pkg/logger/logger.go:1))
```go
type Logger struct {
    debug bool
}

// 主要な責務:
// - 構造化ログ出力
// - ログレベル管理
// - カラー対応出力
```

---

## データフロー

### 1. メインコマンド実行フロー

```mermaid
sequenceDiagram
    participant User
    participant CLI as CLI Layer
    participant Git as Git Module
    participant Alias as Alias Module
    participant Config as Config Module
    participant Emoji as Emoji Module
    
    User->>CLI: pummit sparkles "新機能追加"
    CLI->>Config: 設定読み込み
    Config-->>CLI: 設定情報
    CLI->>Git: GetChangedFiles()
    Git-->>CLI: ["file1.go", "file2.md"]
    CLI->>Alias: GetEmoji("sparkles")
    Alias->>Config: エイリアス検索
    Config-->>Alias: ["sparkles", "✨"]
    Alias-->>CLI: (found=true, prefix="sparkles", emoji="✨")
    CLI->>Emoji: 絵文字変換処理
    Emoji-->>CLI: "✨"
    CLI->>Git: Commit(emoji="✨", message="新機能追加", files="file1.go, file2.md")
    Git-->>User: コミット成功
```

### 2. エイリアス管理フロー

```mermaid
sequenceDiagram
    participant User
    participant CLI as Alias CLI
    participant Alias as Alias Module
    participant Config as Config Module
    participant Emoji as Emoji Module
    
    User->>CLI: pummit alias:add feature sparkles
    CLI->>Emoji: GetEmojiByName("sparkles")
    Emoji-->>CLI: "✨"
    CLI->>Alias: Add("feature", "sparkles", "✨")
    Alias->>Config: 既存エイリアス検索
    Config-->>Alias: 検索結果
    Alias->>Config: エイリアス追加
    Config->>Config: Save()
    Config-->>Alias: 保存完了
    Alias-->>CLI: 追加成功
    CLI-->>User: "Added alias 'feature' for prefix 'sparkles' with emoji '✨'"
```

### 3. 設定管理フロー

```
アプリケーション起動時の設定読み込み:
┌─────────────┐
│ アプリ起動  │
└─────────────┘
        │
        ▼
┌─────────────┐
│config.Init()│
└─────────────┘
        │
        ▼
┌───────────────────┐
│ 設定ファイル存在? │◄─── ~/.config/pummit/config.json をチェック
└───────────────────┘
        │
        ▼
    ┌─────┴─────┐
    │           │
   Yes          No
    │           │
    ▼           ▼
┌──────────┐  ┌─────────────────┐
│既存設定  │  │デフォルト設定   │
│読み込み  │  │作成             │
└──────────┘  └─────────────────┘
    │           │
    │           ▼
    │       ┌─────────────┐
    │       │設定ファイル │
    │       │作成         │
    │       └─────────────┘
    │           │
    └─────┬─────┘
          ▼
┌─────────────────┐
│CurrentConfigに │
│設定をロード     │
└─────────────────┘
          │
          ▼
┌─────────────┐
│アプリ実行   │
└─────────────┘

設定変更時の保存フロー:
┌─────────────┐
│設定変更操作 │ ◄─── pummit alias:add, alias:delete など
└─────────────┘
        │
        ▼
┌─────────────┐
│config.Save()│
└─────────────┘
        │
        ▼
┌─────────────────┐
│ファイルシステム │
│に保存           │ ◄─── JSON形式で ~/.config/pummit/config.json に保存
└─────────────────┘

設定ファイルの構造:
~/.config/pummit/
├── config.json          (メイン設定ファイル)
└── (将来の拡張用)

権限とセキュリティ:
- config.json: 0644 (読み取り専用、所有者のみ書き込み可)
- ディレクトリ: 0755 (標準的なアクセス権限)
```

### 4. エラーハンドリングフロー

```
操作実行時のエラー処理パターン:

┌─────────────┐
│  操作実行   │
└─────────────┘
        │
        ▼
    ┌─────────┐
    │エラー   │
    │発生?    │
    └─────────┘
        │
   ┌────┴────┐
   │         │
  No        Yes
   │         │
   ▼         ▼
┌─────────┐ ┌─────────────┐
│正常処理 │ │エラー種別   │
│続行     │ │判定         │
└─────────┘ └─────────────┘
                │
        ┌───────┼───────┬─────────┐
        │       │       │         │
        ▼       ▼       ▼         ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│ErrAlias  │ │ErrAlias  │ │FileSystem│ │Network   │
│Exists    │ │NotFound  │ │Error     │ │Error     │
└──────────┘ └──────────┘ └──────────┘ └──────────┘
        │       │           │           │
        └───┬───┘           │           │
            │               │           │
            ▼               ▼           ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │ユーザー     │ │システム     │ │リトライ     │
    │フレンドリー │ │エラーログ   │ │機構         │
    │メッセージ   │ └─────────────┘ └─────────────┘
    └─────────────┘         │           │
            │               │           │
            │               ▼           │
            │       ┌─────────────┐     │
            │       │処理終了     │     │
            │       └─────────────┘     │
            │               │           │
            └───────────────┼───────────┘
                            │
                            ▼
                    ┌─────────────┐
                    │最終処理終了 │
                    └─────────────┘

エラー種別と対応:

[1] ErrAliasExists
    出力: "エイリアス '%s' は既に存在します"
    動作: 処理中断、現在の設定を維持

[2] ErrAliasNotFound
    出力: "エイリアス '%s' が見つかりません"
    動作: 処理中断、何も変更しない

[3] FileSystemError
    出力: "設定ファイルの読み書きに失敗しました: %v"
    動作: 詳細なエラーログを記録、処理中断

[4] NetworkError
    出力: "ネットワーク接続に失敗しました (フォールバック機能を使用)"
    動作: ローカルの絵文字データを使用してリトライ

[5] その他のエラー
    出力: "予期しないエラーが発生しました: %v"
    動作: スタックトレースをログに記録、処理中断
```

---

## 設定システム

### 設定ファイル構造

**場所**: `~/.config/pummit/config.json`

```json
{
  "writeEmoji": true,
  "useAlias": true,
  "useLimitPathesLength": true,
  "limitPathesLength": 50,
  "alias": [
    ["s,feat,feature", "sparkles", "✨"],
    ["c,wip", "construction", "🚧"],
    ["t,new,init", "tada", "🎉"],
    ["r,pr,pull,merge", "recycle", "♻️"],
    ["wb,rm,remove,del,delete", "wastebasket", "🗑️"],
    ["b,fix,error", "bug", "🐛"],
    ["e,lint,format,refactor", "eyes", "👀"],
    ["d,doc,docs,document,documents", "books", "📚"],
    ["a,ui,design,icon,icons", "art", "🎨"],
    ["h,tune,tuning,perform,performance", "horse", "🐎"],
    ["w,change,tool,tools,lib,library", "wrench", "🔧"],
    ["l,test,testing", "rotating_light", "🚨"],
    ["sm,special,important", "snowman", "☃️"],
    ["p,pack,mod,module", "package", "📦️"]
  ]
}
```

### 設定項目詳細

| 項目 | 型 | デフォルト | 説明 |
|------|----|---------:|------|
| `writeEmoji` | `bool` | `true` | 絵文字を直接出力するか（false の場合は `:name:` 形式） |
| `useAlias` | `bool` | `true` | エイリアス機能を使用するか |
| `useLimitPathesLength` | `bool` | `true` | ファイルパス長の制限を使用するか |
| `limitPathesLength` | `int` | `50` | ファイルパス表示の最大文字数 |
| `alias` | `[][]string` | デフォルトエイリアス | エイリアス定義配列 |

### エイリアス配列形式

```go
// 形式1: [エイリアス名(複数可), プレフィックス, 絵文字]
["s,feat,feature", "sparkles", "✨"]

// 形式2: [エイリアス名(複数可), 絵文字のみ]
["custom", "🔥"]
```

### 設定初期化プロセス

1. **起動時**: [`config.Init()`](internal/config/config.go:1008)
2. **ホームディレクトリ取得**: `os.UserHomeDir()`
3. **設定ディレクトリ作成**: `~/.config/pummit/`
4. **デフォルト設定読み込み**: 埋め込み JSON から
5. **既存設定確認**: 設定ファイルの存在チェック
6. **設定マージ**: 既存設定 + デフォルト設定
7. **設定保存**: 最新設定をファイルシステムに保存

---

## 外部依存関係

### Go言語標準ライブラリ

- `encoding/json`: 設定ファイル処理
- `os`: ファイルシステム操作
- `os/exec`: Git コマンド実行
- `path/filepath`: ファイルパス操作
- `net/http`: Gitmoji API 通信

### サードパーティライブラリ

| ライブラリ | バージョン | 用途 | ライセンス |
|-----------|-----------|------|-----------|
| [Cobra](https://github.com/spf13/cobra) | v1.9.1 | CLI フレームワーク | Apache 2.0 |
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | v1.3.4 | TUI フレームワーク | MIT |
| [Bubbles](https://github.com/charmbracelet/bubbles) | v0.20.0 | TUI コンポーネント | MIT |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | v1.0.0 | スタイリング | MIT |
| [Color](https://github.com/fatih/color) | v1.18.0 | カラー出力 | MIT |

### 外部サービス

**Gitmoji API**
- **URL**: `https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json`
- **用途**: 最新の Gitmoji データ取得
- **フォールバック**: 埋め込み絵文字データ

### 依存関係管理戦略

1. **最小依存原則**: 必要最小限の外部依存
2. **バージョン固定**: セマンティックバージョニング遵守
3. **ライセンス互換性**: Apache 2.0 と互換性のあるライセンスのみ
4. **セキュリティ**: 定期的な脆弱性チェック

---

## MCPサーバー機能

### 1. 概要

pummitにModel Context Protocol (MCP) サーバー機能を追加し、外部のMCPクライアント（例: Claude Desktop）からpummitの機能を利用可能にします。これにより、AIアシスタントとの連携を強化し、自然言語によるGit操作や設定管理を実現します。

### 2. アーキテクチャとプロトコル実装

- **統合方法**: `pummit mcp` コマンドを新しいエントリーポイントとして追加します。このコマンドは標準入出力（stdio）を介してJSON-RPC通信を行い、MCPサーバーとして動作します。
- **プロトコル実装**: `github.com/mark3labs/mcp-go` ライブラリを利用してサーバーを構築します。このライブラリはMCPの高レベルな抽象化を提供し、プロトコルの詳細を意識することなくツールやリソースの定義に集中できます。
- **既存ロジックの再利用**: MCPハンドラは、`internal/git`, `internal/config`, `internal/alias`, `internal/doctor` といった既存のビジネスロジック層のモジュールを再利用し、機能の一貫性を保ちます。

### 3. コマンド体系

```bash
# MCPサーバーを標準入出力モードで起動
pummit mcp
```

### 4. モジュール構成

`internal/mcp/` ディレクトリを新設し、MCPサーバー関連のコードを集約します。

```
internal/
├── mcp/
│   ├── server.go        # MCPサーバーの初期化とツール登録
│   ├── tool_git.go      # Git関連ツールの定義とハンドラ
│   ├── tool_config.go   # 設定関連ツールの定義とハンドラ
│   ├── tool_alias.go    # エイリアス関連ツールの定義とハンドラ
│   └── tool_doctor.go   # 診断ツールの定義とハンドラ
└── cli/
    └── mcp.go           # `pummit mcp` コマンドの定義
```

### 5. ツール定義

LLMが自律的なワークフローを実行できるよう、以下のプリミティブなツールを提供します。

#### `git` ツール

ここの説明を実際に実装するときは英語で書いてください。実際の実装上での`git.smart_commit`と`git.commit`の説明にこの文章をそのまま使ってください。（またこれは重要なことなのでllm-context.mdにメモること）

| ツール名 | 説明 | 入力パラメータ | 出力例（成功時） |
| :--- | :--- | :--- | :--- |
| `git.smart_commit` | **PRIMARY TOOL FOR COMMITTING.** Intelligently analyze repository changes and create a commit with auto-generated message and emoji. Use this when the user simply says 'commit', 'make a commit', or gives minimal instructions like 'commit the changes' without specifying exact files or messages. Handles the complete workflow: file analysis, message generation, and execution. | `message_hint: string` (optional), `auto_approve: bool` (default: false) | `{ "status": "pending_approval", "plan": {...} }` or `{ "success": true }` |
| `git.get_edited_files` | `git status -s -u` を実行し、変更があったファイルの一覧をショートフォーマットで取得します。 | `(none)` | `{ "status": " M internal/git/git.go\n?? new_file.txt" }` |
| `git.get_current_branch` | 現在のGitブランチ名を取得します。 | `(none)` | `{ "branch": "main" }` |
| `git.add_files` | 指定されたファイルをステージングします。 | `files: string[]` | `{ "success": true }` |
| `git.commit` | Create a Git commit with specific user-provided emoji, message, and staging preferences. Use this ONLY when the user provides explicit commit details (specific message, emoji, or file selection). For simple 'commit' requests, use git.smart_commit instead. | `emoji: string`, `message: string`, `offline: bool` | `{ "success": true }` |
| `git.move_branch` | 指定されたブランチにチェックアウトします。 | `branch: string` | `{ "success": true }` |
| `git.create_branch` | 新しいブランチを作成してチェックアウトします。 | `branch: string` | `{ "success": true }` |

#### `alias` ツール

| ツール名 | 説明 | 入力パラメータ | 出力例（成功時） |
| :--- | :--- | :--- | :--- |
| `alias.list_aliases` | すべてのエイリアスの一覧を取得します。 | `(none)` | `{ "aliases": { "feat": "✨", "fix": "🐛" } }` |

#### `doctor` ツール

| ツール名 | 説明 | 入力パラメータ | 出力例（成功時） |
| :--- | :--- | :--- | :--- |
| `doctor.run_diagnostics` | システムの健全性を診断します。 | `(none)` | `{ "results": [...] }` |

### 6. LLMの自律的ワークフロー

このツールセットにより、LLMは「コミットして」という単純な指示から、以下のような自律的ワークフローを実行できます。

1.  **現状把握**: `git.get_edited_files` と `git.get_current_branch` を実行し、「どのファイルが変更され、どのブランチにいるか」を把握します。
2.  **計画立案**: 変更されたファイルの内容を分析し、どのようなコミットが適切かを判断します。同時に `alias.list_aliases` を実行して、利用可能な絵文字エイリアスを確認します。
3.  **実行**:
    -   `git.add_files` を呼び出し、コミット対象のファイルをステージングします。
    -   ファイルの内容とエイリアスに基づき、最適なコミットメッセージ（例: `✨ 新機能を追加`）を生成します。
    -   `git.commit` を呼び出して、生成したメッセージでコミットを実行します。

### 7. エラーハンドリング設計

LLMがエラー発生時に自律的に回復を試みたり、ユーザーに的確な報告をしたりできるよう、エラー情報を構造化し、具体的な対応方針を定めます。(実際にはこれらの出力は英語で)

| ツール | エラーシナリオ | LLMの対応策 | LLMからユーザーへの報告メッセージ例 |
| :--- | :--- | :--- | :--- |
| **`git.get_edited_files`** | Gitリポジトリではない | 処理を中断し、ユーザーに `git init` の実行を提案する。 | 「おっと、ここはGitリポジトリではないようです。コミット機能を使うには、まずリポジトリを初期化する必要があります。`git init` を実行して初期化しますか？」 |
| **`git.add_files`** | 指定ファイルが存在しない | `get_edited_files` を再実行してファイル一覧を再確認し、正しいファイル名で再試行する。失敗が続く場合はユーザーに報告する。 | 「ファイル `non_existent_file.go` が見つかりませんでした。もう一度ファイル名を確認してみますね... やはり見つからないようです。ファイル名が正しいか、またはファイルが削除されていないか確認していただけますか？」 |
| **`git.commit`** | ステージングされたファイルがない | `get_edited_files` を実行し、変更があれば `add_files` を呼び出してから `commit` を再試行する。変更がなければユーザーに報告する。 | 「コミットする変更が見つかりませんでした。ですが、いくつか変更されたファイルがあるようです。これらのファイルをステージングしてコミットしますか？」 |
| | Gitの `user.name`/`user.email` が未設定 | `doctor` ツールを実行して診断し、具体的な設定コマンドをユーザーに提示する。 | 「コミットを実行するために、Gitのユーザー設定が必要です。`doctor` ツールで確認したところ、`user.name` と `user.email` が設定されていないようです。以下のコマンドを実行して設定してください。\n```\ngit config --global user.name \"Your Name\"\ngit config --global user.email \"you@example.com\"\n```」 |
| **`git.move_branch`** | 指定ブランチが存在しない | ユーザーにブランチが存在しないことを報告し、`create_branch` を使って新規作成するかどうかを尋ねる。 | 「ブランチ `feature/new-idea` は存在しないようです。同名の新しいブランチを作成して、そちらに移動しますか？」 |
| **`git.create_branch`** | 指定ブランチが既に存在する | ユーザーにブランチが既に存在することを報告し、`move_branch` を使ってそのブランチに移動するかどうかを尋ねる。 | 「ブランチ `feature/existing-work` は既に存在していますね。そのブランチに移動しますか？」 |

### 8. テスト戦略

- **単体テスト**: 各ツールのハンドラ関数が、正常系・異常系の両方で正しく動作することを検証します。既存のビジネスロジックはモック化します。
- **統合テスト**: `mcp-go`のテスト機能を利用し、サーバー全体がMCPリクエストに対して正しく応答できるかを確認します。実際のstdioを模した入力を用いて、リクエストからレスポンスまでの一連の流れをテストします。

---

## 設定ファイルマイグレーション仕様

### JSON → TOML 移行

#### 概要
従来のJSON形式設定ファイルから、より読みやすく管理しやすいTOML形式への移行を行います。

#### 新しいTOML設定構造
```toml
[meta]
version = "3.0"

[base]
emoji = true # writeEmoji から変更
filesLength = 50 # 0 = 無効, -1 = 無制限

[interactive]
enabled = true
defaultMode = "emoji" # "emoji", "template", "scope", "branch"
showPreview = true
fuzzySearch = true

[locale]
language = "ja"  # "en", "ja", "zh", "ko"
autoDetect = true # システム言語から自動検出

[templates]
enabled = true
defaultTemplate = "default"

[templates.definitions.default]
format = "{emoji} {message} ({files})"
description = "標準テンプレート"

[templates.definitions.feat]
format = "{emoji} {scope}: {message}\n\n{description}\n\n({files})"
description = "新機能追加用テンプレート"
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
shortcuts = ["h", "tune", "tuning", "perform", "performance"]
name = "horse"
emoji = "🐎"

[[alias.entries]]
shortcuts = ["w", "change", "tool", "tools", "lib", "library"]
name = "wrench"
emoji = "🔧"

[[alias.entries]]
shortcuts = ["l", "test", "testing"]
name = "rotating_light"
emoji = "🚨"

[[alias.entries]]
shortcuts = ["sm", "special", "important"]
name = "snowman"
emoji = "☃️"

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
description = "新機能ブランチ"

[[branchMapping.rules]]
pattern = "^(fix|bugfix|hotfix)/.*"
emoji = "bug"
description = "バグ修正ブランチ"

[[branchMapping.rules]]
pattern = "^docs/.*"
emoji = "books"
description = "ドキュメントブランチ"

[[branchMapping.rules]]
pattern = "^refactor/.*"
emoji = "eyes"
description = "リファクタリングブランチ"

[[branchMapping.rules]]
pattern = "^test/.*"
emoji = "rotating_light"
description = "テストブランチ"
```

#### マイグレーション仕様

**ファイル検索優先順位**
1. `~/.config/pummit/config.toml` （最優先）
2. `~/.config/pummit/config.json` （レガシー、自動変換対象）
3. デフォルト設定で新規作成

**マイグレーション処理フロー**
```
起動時設定読み込み:
┌─────────────────┐
│config.toml      │
│存在チェック     │
└─────────────────┘
         │
    ┌────┴────┐
   Yes        No
    │          │
    ▼          ▼
┌─────────┐ ┌─────────────────┐
│TOML     │ │config.json      │
│読み込み │ │存在チェック     │
└─────────┘ └─────────────────┘
              │
         ┌────┴────┐
        Yes        No
         │          │
         ▼          ▼
    ┌─────────┐ ┌─────────────┐
    │JSON→TOML │ │デフォルト   │
    │変換実行  │ │TOML作成     │
    └─────────┘ └─────────────┘
         │          │
         ▼          │
    ┌─────────┐     │
    │JSONファイル    │
    │.bakに保存│     │
    └─────────┘     │
         │          │
         └────┬─────┘
              ▼
    ┌─────────────┐
    │設定ロード   │
    │完了         │
    └─────────────┘
```

**JSON→TOML変換ルール**
```go
// 変換マッピング
type MigrationMapping struct {
    // 基本設定
    "writeEmoji" → "base.emoji"
    "useAlias" → "alias.enabled"
    "useLimitPathesLength" → "base.filesLength" (true/false → 50/0)
    "limitPathesLength" → "base.filesLength"
    
    // エイリアス変換
    "alias": [][]string → "alias.entries": []AliasEntry
    // ["s,feat", "sparkles", "✨"] → {shortcuts=["s","feat"], name="sparkles", emoji="✨"}
}
```

**コマンド仕様**
```bash
# 手動マイグレーション
pummit migrate
pummit migrate --force  # 既存TOMLファイルを上書き
pummit migrate --dry-run  # 変換結果のプレビューのみ

# 設定管理コマンド（統合版）
pummit config                                    # 設定ファイルをエディターで開く（デフォルト動作）
pummit config --get base.emoji                  # 設定値の取得
pummit config --set base.emoji=true             # 設定値の変更
pummit config --list                            # 全設定の一覧表示
pummit config --reset                           # 全設定をデフォルトに戻す
pummit config --reset base                      # 特定セクションのみリセット
pummit config --validate                        # 設定ファイルの整合性確認
pummit config --path                            # 設定ファイル場所表示
pummit config --status                          # 設定ファイル状態確認

# 設定バックアップ・復元（シンプル化）
pummit config --backup PATH                     # 設定のバックアップ
pummit config --restore PATH                    # 設定の復元

# 診断機能（早期MVP実装）
pummit doctor                                    # 設定検証 + Git/Terminal 情報出力
```

#### 実装構造
```
internal/
├── config/
│   ├── config.go        # 設定管理（拡張）
│   ├── toml.go         # TOML形式対応（BurntSushi/toml使用予定）
│   ├── migration.go    # マイグレーション処理
│   ├── validator.go    # 設定バリデーション
│   ├── operations.go   # 設定操作（get/set/list/reset統合）
│   └── backup.go       # バックアップ・復元機能
├── cli/
│   ├── migrate.go      # マイグレーションコマンド
│   ├── config.go       # 統合設定コマンド
│   └── doctor.go       # 診断コマンド（MVP）
└── telemetry/          # 可観測性（opt-in）
    ├── metrics.go      # OpenTelemetry（環境変数制御）
    └── tracing.go      # 起動時間・操作トレース
```

---

## 新機能仕様

### 実装優先順位

1. **設定ファイルマイグレーション機能** （最優先）
2. **ブランチ名からのプレフィックス自動提案**
3. **インタラクティブな絵文字・エイリアス選択機能の強化**
4. **スコープ入力支援**
5. **コミットメッセージテンプレート機能**

### 1. ブランチ名からのプレフィックス自動提案

#### 概要
現在のブランチ名から適切な絵文字プレフィックスを自動推測・提案する機能。

#### 実装仕様

**コマンド動作例**
```bash
# feature/user-authブランチの場合
$ git branch --show-current
feature/user-auth

$ pummit "ユーザー認証機能を追加"
🔍 ブランチ名から推測: feature/user-auth → ✨ sparkles
✨ ユーザー認証機能を追加 (src/auth.go, tests/auth_test.go)

# 複数候補がある場合
$ pummit "認証機能修正"
🔍 ブランチ名から複数の候補が見つかりました:
1. ✨ sparkles (新機能)
2. 🐛 bug (バグ修正)
選択してください [1]:
```

**実装構造**
```
internal/
├── branch/
│   ├── analyzer.go      # ブランチ名解析
│   ├── mapping.go       # マッピングルール
│   └── suggestion.go    # 提案エンジン
└── git/
    └── branch.go        # ブランチ情報取得（拡張）
```

### 2. インタラクティブモード (`pummit interactive`)

#### 概要
ファジーファインダー機能とステップ式UIを組み込んだ包括的なインタラクティブコミット作成機能。

#### コマンド仕様
```bash
# インタラクティブモード起動
pummit interactive

# 特定モードで起動
pummit interactive --mode emoji      # 絵文字選択から開始
pummit interactive --mode template   # テンプレート選択から開始
pummit interactive --mode scope      # スコープ入力から開始
pummit interactive --mode branch     # ブランチ提案から開始
```

#### UI仕様（多言語対応）

**1. モード選択画面**
```
┌─────────────────────────────────────────┐
│ 🎯 Pummit Interactive Mode             │
├─────────────────────────────────────────┤
│ How would you like to create commit?    │
├─────────────────────────────────────────┤
│ > 🔍 Select emoji/alias                 │
│   📝 Choose template                    │
│   🎯 Specify scope                      │
│   ⚡ Auto-suggest from branch          │
├─────────────────────────────────────────┤
│ [↑↓] Select  [Enter] Confirm  [q] Quit │
└─────────────────────────────────────────┘
```

**2. 絵文字選択画面（ファジーファインダー）**
```
┌─────────────────────────────────────────────────────────┐
│ 🔍 Search emoji & alias                                │
├─────────────────────────────────────────────────────────┤
│ > feat                                                  │
├─────────────────────────────────────────────────────────┤
│ ✨ sparkles (s,feat,feature) - New feature             │
│ 🎉 tada (t,new,init) - Initial commit                  │
│ 🐛 bug (b,fix,error) - Bug fix                         │
│ 📚 books (d,doc,docs) - Documentation                  │
│ 🔧 wrench (w,change,tool) - Tool changes               │
│                                                         │
│ [↑↓] Select  [Enter] Confirm  [Esc] Back               │
└─────────────────────────────────────────────────────────┘
```

**3. メッセージ入力画面**
```
┌─────────────────────────────────────────────────────────┐
│ ✨ sparkles selected                                    │
├─────────────────────────────────────────────────────────┤
│ Enter commit message:                                   │
│ > Add user authentication feature                       │
├─────────────────────────────────────────────────────────┤
│ Preview:                                                │
│ ✨ Add user authentication feature (src/auth.go, ...)  │
├─────────────────────────────────────────────────────────┤
│ [Enter] Commit  [Ctrl+C] Cancel  [Tab] Options         │
└─────────────────────────────────────────────────────────┘
```

#### 多言語サポート仕様

**対応言語**
- 🇺🇸 English (en)
- 🇯🇵 日本語 (ja)
- 🇨🇳 中文 (zh)
- 🇰🇷 한국어 (ko)

**言語設定**
```toml
[locale]
language = "ja"  # "en", "ja", "zh", "ko"
autoDetect = true # システム言語から自動検出
```

**翻訳例**
```go
// 英語
"Select emoji/alias" → "絵文字・エイリアス選択"
"Enter commit message" → "コミットメッセージを入力"
"New feature" → "新機能"

// 中国語
"Select emoji/alias" → "选择表情符号/别名"
"Enter commit message" → "输入提交消息"
"New feature" → "新功能"

// 韓国語
"Select emoji/alias" → "이모지/별칭 선택"
"Enter commit message" → "커밋 메시지 입력"
"New feature" → "새로운 기능"
```

#### 実装構造
```
internal/
├── interactive/
│   ├── model.go         # Bubble Tea モデル
│   ├── modes/
│   │   ├── emoji.go     # 絵文字選択モード
│   │   ├── template.go  # テンプレート選択モード
│   │   ├── scope.go     # スコープ入力モード
│   │   └── branch.go    # ブランチ提案モード
│   ├── fuzzy.go         # ファジー検索
│   └── i18n.go          # 国際化対応
├── locale/
│   ├── en.json          # 英語翻訳
│   ├── ja.json          # 日本語翻訳
│   ├── zh.json          # 中国語翻訳
│   └── ko.json          # 韓国語翻訳
└── cli/
    └── interactive.go   # interactiveサブコマンド
```

### 3. スコープ入力支援

#### 概要
Conventional Commits形式のスコープ（例: `feat(api): ...`）の入力を支援する機能。

#### 自動検出機能
```
自動検出ルール:
1. 変更ファイルパスからスコープを推測
   src/api/ → api
   frontend/ → ui
   docs/ → docs

2. 過去のコミット履歴からスコープを抽出
   git log --pretty=format:"%s" | grep -E "(feat|fix)\([^)]+\)"

3. ブランチ名からスコープを推測
   feature/api-auth → api
   fix/ui-button → ui
```

### 4. コミットメッセージテンプレート機能

#### 概要
プロジェクトやコミットの種類に応じて、予め定義されたコミットメッセージのテンプレートを適用できる機能。

#### コマンド仕様
```bash
# テンプレート使用
pummit --template feat "ユーザー認証機能を追加"

# インタラクティブモード
pummit --interactive

# テンプレート一覧表示
pummit template:list

# テンプレート追加
pummit template:add custom "カスタムテンプレート" "{emoji} [{scope}] {message}"
```

---
## 品質保証・テスト戦略

### テスト方針

#### 1. テスト構成
```
tests/
├── unit/                    # 単体テスト
│   ├── config_test.go      # 設定管理
│   ├── alias_test.go       # エイリアス機能
│   ├── git_test.go         # Git操作
│   └── migration_test.go   # マイグレーション
├── integration/            # 統合テスト
│   ├── cli_test.go         # CLIコマンド
│   ├── interactive_test.go # インタラクティブモード
│   └── e2e_test.go         # エンドツーエンド
├── compatibility/          # 互換性テスト
│   ├── terminal_test.go    # ターミナル互換性
│   ├── os_test.go          # OS互換性
│   └── git_versions_test.go # Gitバージョン互換性
└── performance/            # パフォーマンステスト
    ├── startup_test.go     # 起動時間
    └── memory_test.go      # メモリ使用量
```

#### 2. テストカバレッジ目標（現実的数値）
- **単体テスト**: 80% 以上
- **統合テスト**: 主要フローの 80-85% 網羅（入力組み合わせ爆発を避ける）
- **E2Eテスト**: クリティカルパスの 85% 網羅（Fail-Fast文化重視）
- **全体目標**: 80% 以上を維持し、品質と保守性のバランスを重視

#### 3. CI/CD テスト戦略
```yaml
# .github/workflows/test.yml の概要
test:
  strategy:
    matrix:
      os: [ubuntu-latest, windows-latest, macos-latest]
      go: ['1.20', '1.21', '1.22']
  steps:
    - name: Unit Tests
      run: go test ./... -race -coverprofile=coverage.out
    
    - name: Integration Tests
      run: go test ./tests/integration/... -v
    
    - name: Vulnerability Scan
      run: govulncheck ./...
    
    - name: SBOM Generation
      run: syft . -o cyclonedx-json=sbom.json
```

#### 4. 特殊テスト項目

**マイグレーションテスト**
```go
func TestJSONToTOMLMigration(t *testing.T) {
    tests := []struct {
        name     string
        jsonData string
        expected TOMLConfig
        wantErr  bool
    }{
        {
            name: "basic migration",
            jsonData: `{"writeEmoji": true, "useAlias": true}`,
            expected: TOMLConfig{Base: BaseConfig{Emoji: true}},
        },
        {
            name: "complex alias migration",
            jsonData: `{"alias": [["s,feat", "sparkles", "✨"]]}`,
            expected: TOMLConfig{/* ... */},
        },
        // エッジケース: エスケープ文字、重複alias
    }
}
```

**ターミナル互換性テスト**
```go
func TestTerminalCompatibility(t *testing.T) {
    terminals := []string{
        "xterm", "xterm-256color", "screen", "tmux", 
        "windows-powershell", "windows-cmd",
    }
    
    for _, term := range terminals {
        t.Run(term, func(t *testing.T) {
            // Bubble Tea UI の描画テスト
            // 文字幅計算の検証
            // プレビュー表示の確認
        })
    }
}
```

**起動時間パフォーマンステスト**
```go
func BenchmarkStartupTime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        start := time.Now()
        
        // アプリケーション起動処理
        cmd := exec.Command("pummit", "--version")
        _ = cmd.Run()
        
        elapsed := time.Since(start)
        if elapsed > 100*time.Millisecond {
            b.Errorf("起動時間が目標値を超過: %v", elapsed)
        }
    }
}
```

### 継続的品質保証

#### 1. 自動品質チェック
```yaml
quality-gates:
  - name: "Test Coverage"
    threshold: 80%
    action: "block-merge"
    
  - name: "Startup Time"
    threshold: 100ms
    action: "warning"
    
  - name: "Memory Usage"
    threshold: 10MB
    action: "warning"
    
  - name: "Vulnerability Scan"
    severity: "high"
    action: "block-merge"
```

#### 2. 依存関係管理
```yaml
dependabot:
  updates:
    - package-ecosystem: "gomod"
      directory: "/"
      schedule:
        interval: "weekly"
      reviewers:
        - "security-team"
      
    - package-ecosystem: "github-actions"
      directory: "/.github/workflows"
      schedule:
        interval: "monthly"
```

#### 3. リリース前チェックリスト
- [ ] 全テストスイートの成功
- [ ] クロスプラットフォーム動作確認
- [ ] マイグレーションテスト（旧バージョンからの移行）
- [ ] パフォーマンス回帰テスト
- [ ] セキュリティスキャン（govulncheck + SBOM）
- [ ] ドキュメント更新確認
- [ ] 翻訳ファイルの整合性確認

---

## セキュリティ考慮事項

### 1. ファイルシステム操作

**脅威**: パストラバーサル攻撃
**対策**: 
- 設定ファイルパスの検証
- ホームディレクトリ外への書き込み防止

```go
// セキュアなパス検証例
func validateConfigPath(path string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    
    configDir := filepath.Join(homeDir, ".config", "pummit")
    cleanPath := filepath.Clean(path)
    
    if !strings.HasPrefix(cleanPath, configDir) {
        return fmt.Errorf("invalid config path: %s", path)
    }
    return nil
}
```

### 2. 外部 API 通信

**脅威**: 中間者攻撃、データ改ざん
**対策**:
- HTTPS通信の強制
- タイムアウト設定
- レスポンス検証

```go
// セキュアな HTTP クライアント
func createSecureClient() *http.Client {
    return &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
            },
        },
    }
}
```

### 3. ユーザー入力検証

**脅威**: インジェクション攻撃
**対策**:
- 入力のサニタイゼーション
- エスケープ処理
- 長さ制限

```go
// 安全な入力検証
func validateCommitMessage(message string) error {
    if len(message) > 72 {
        return fmt.Errorf("commit message too long: %d chars", len(message))
    }
    
    // 危険な文字の検証
    if strings.ContainsAny(message, "\n\r\t") {
        return fmt.Errorf("commit message contains invalid characters")
    }
    
    return nil
}
```

### 4. 権限管理

**原則**: 最小権限の原則
**実装**:
- 実行に必要な最小限の権限
- 設定ファイルの適切なパーミッション (0644)
- 一時ファイルの安全な処理

---

## パフォーマンス設計

### 1. 起動時間最適化

**目標**:
- 基本動作: 100ms 以下
- OpenTelemetry有効時: 110ms 以下（+10ms許容）
- Bubble Tea UI使用時: 120ms 以下（実測ベース調整）

**戦略**:
- 遅延初期化の活用
- TOMLライブラリ選定（BurntSushi/toml vs pelletier/go-toml）
- 設定読み込みの最適化
- OpenTelemetry のopt-in設計（PUMMIT_OTEL_ENDPOINT環境変数制御）

**計測・監視**:
```go
// CI での起動時間ベンチマーク
func BenchmarkStartupTime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        start := time.Now()
        cmd := exec.Command("pummit", "--version")
        _ = cmd.Run()
        elapsed := time.Since(start)
        
        // レッド/イエロー/グリーンの閾値
        if elapsed > 120*time.Millisecond {
            b.Errorf("Booting time red: %v > 120ms", elapsed)
        } else if elapsed > 100*time.Millisecond {
            b.Logf("Booting time yellow: %v", elapsed)
        } else {
            b.Logf("Booting time green %v", elapsed)
        }
    }
}
```

```go
// 遅延初期化の例
var (
    emojiCache     map[string]string
    emojiCacheOnce sync.Once
)

func getEmojiCache() map[string]string {
    emojiCacheOnce.Do(func() {
        emojiCache = loadEmojiCache()
    })
    return emojiCache
}
```

### 2. メモリ使用量最適化

**目標**: 10MB 以下のメモリ使用量
**戦略**:
- 大きなデータ構造の共有
- 不要なオブジェクトの早期解放
- メモリプールの活用

### 3. ネットワーク最適化

**戦略**:
- HTTP/2 対応
- 接続の再利用
- レスポンスキャッシュ

```go
// HTTP クライアントの最適化
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
        DisableCompression: false,
    },
}
```

### 4. ファイル操作最適化

**戦略**:
- バッファードI/O
- 並行処理の活用
- ファイルサイズの最小化

## リスク分析と対策

### 技術的リスク

#### R1. JSON→TOML変換時の損失・不一致
**リスク**: エイリアス正規化、エスケープ文字処理でユーザー設定が破損
**影響度**: 高 | **発生確率**: 中
**対策**:
- 変換前の自動バックアップ（.bak ファイル）
- ドライラン機能での事前検証
- 詳細な変換ログ出力
- ロールバック機能の提供

```go
// 実装例：安全な変換処理
func SafeMigration(jsonPath string) error {
    // 1. バックアップ作成
    if err := createBackup(jsonPath); err != nil {
        return err
    }
    
    // 2. ドライラン実行
    result, err := dryRunConversion(jsonPath)
    if err != nil {
        return err
    }
    
    // 3. ユーザー確認
    if !confirmConversion(result) {
        return ErrUserCancelled
    }
    
    // 4. 実際の変換実行
    return executeConversion(jsonPath, result)
}
```

#### R2. Git操作のクロスOS差異
**リスク**: パス区切り文字、改行コード、gitconfig設定の違い
**影響度**: 中 | **発生確率**: 高
**対策**:
- filepath.Join()によるパス正規化
- Git操作の標準化関数群
- OS別テストスイート
- libgit2バインディングへの段階的移行

#### R3. Bubble Tea Windows PowerShell互換性
**リスク**: TUI描画の崩れ、キー入力の誤認識
**影響度**: 中 | **発生確率**: 中
**対策**:
- サポートターミナル一覧の明文化
- `--no-interactive` フラグでのフォールバック
- Windows用ptyラッパの導入
- PowerShell互換性テストの自動化

#### R4. Gitmoji API依存による初回起動ブロック
**リスク**: ネットワーク障害時の長時間待機
**影響度**: 中 | **発生確率**: 中
**対策**:
- `--offline` フラグの実装
- 埋め込み絵文字データでのフォールバック
- タイムアウト設定の短縮（3秒以内）
- バックグラウンド更新の実装

```go
// 実装例：オフライン対応
func GetEmojiData(offline bool) ([]Emoji, error) {
    if offline {
        return getEmbeddedEmojis(), nil
    }
    
    // 短いタイムアウトでAPI試行
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    if data, err := fetchFromAPI(ctx); err == nil {
        return data, nil
    }
    
    // API失敗時は埋め込みデータを使用
    log.Warn("API unavailable, using embedded emoji data")
    return getEmbeddedEmojis(), nil
}
```

#### R5. マルチバイト絵文字の幅計算差異
**リスク**: UI崩れ、72文字制限違反
**影響度**: 低 | **発生確率**: 高
**対策**:
- runewidth.StringWidth()での正確な幅計算
- Unicode正規化の実装
- プレビュー機能での事前確認
- 設定可能な文字数制限

### 運用リスク

#### R6. 設定フォーマット進化時の後方互換試験コスト
**リスク**: バージョンアップでユーザー設定が動作しない
**影響度**: 高 | **発生確率**: 中
**対策**:
- 設定スキーマのバージョン管理
- 自動バリデーション機能
- CI/CDでの互換性テスト自動化
- 段階的移行計画

#### R7. 自動アップデート機能不在
**リスク**: セキュリティパッチの遅延
**影響度**: 中 | **発生確率**: 高
**対策**:
- `pummit update`コマンドの実装
- GitHub Releasesとの連携
- セキュリティ通知機能
- Homebrew/Scoopでの自動更新推奨

#### R8. ログの分散と集中収集困難
**リスク**: 組織でのトラブルシューティング効率低下
**影響度**: 低 | **発生確率**: 中
**対策**:
- OpenTelemetry対応
- 構造化ログ（JSON形式）
- ログ出力先の設定可能化
#### R9. 設定CLI機能の複雑化
**リスク**: UXが複雑になりユーザーの学習コスト増大
**影響度**: 中 | **発生確率**: 高
**対策**:
- `pummit config` 単一エントリポイントへの統合
- デフォルト動作をエディター起動に設定
- パワーユーザー向けオプションはヘルプに記載
- 使用頻度の高い操作の優先実装

#### R10. TOML・JSON デュアル対応期間のコード分岐
**リスク**: 設定ロジック分岐による保守性悪化・バグ混入
**影響度**: 中 | **発生確率**: 中
**対策**:
- 移行完了後の明確な廃止ロードマップ策定
- v3.0でJSON非推奨警告 → v4.0で完全削除
- 分岐ロジックの単体テスト強化
- 共通インターフェースによる抽象化

#### R11. SBOM・脆弱性スキャン生成物の情報漏洩
**リスク**: CI生成物に機密情報を含めて誤って公開
**影響度**: 高 | **発生確率**: 低
**対策**:
- CI Artifacts の retention policy を60日以内に設定
- 内部限定公開の徹底
- SBOM出力内容の事前レビュー
- 自動削除スクリプトの実装
- 企業向けログ収集ガイド

### リスクマトリックス

```
影響度 ＼ 発生確率    高      中      低
─────────────────────────────────
高                R1,R6           
中                R2     R3,R4,R7  
低                R5              R8
```

### 緊急時対応手順

#### 1. 設定破損時の復旧
```bash
# バックアップからの復元
pummit config:restore ~/.config/pummit/config.json.bak

# デフォルト設定での初期化
pummit config:reset --force

# 手動修復モード
pummit config:repair --interactive
```

#### 2. 起動不能時の診断
```bash
# 詳細診断情報の出力
pummit doctor

# 設定ファイルの検証
pummit config:validate

# 依存関係の確認
pummit deps:check
```

### 監視・可観測性

#### 1. メトリクス収集
```go
// OpenTelemetry メトリクス例
var (
    startupDuration = otel.NewHistogram("pummit_startup_duration_ms")
    apiLatency     = otel.NewHistogram("pummit_api_latency_ms")
    errorCount     = otel.NewCounter("pummit_errors_total")
)
```

#### 2. ヘルスチェック機能
```bash
# システム状態の確認
pummit health
# Output:
# ✅ Git: /usr/bin/git (version 2.39.0)
# ✅ Config: ~/.config/pummit/config.toml (valid)
# ⚠️  Network: API timeout (using offline mode)
# ✅ Terminal: xterm-256color (supported)
```

---
---

## 今後の開発指針

### 短期目標（1〜2週間）- 緊急改善

1. **技術基盤の確立**
   - TOMLライブラリ選定ドキュメント作成（BurntSushi/toml vs pelletier/go-toml）
   - 起動時間ベンチマークをCIで実測・可視化（レッド/イエロー/グリーン）
   - `pummit doctor` MVP実装（設定検証 + Git/Terminal情報出力）

2. **即座のリスク軽減**
   - README にサポート OS/ターミナル表を追加
   - `--offline` フラグ実装（Gitmoji API 不在時対応）
   - `go test ./... -race` を CI に追加
   - Windows PowerShell 互換性の明文化

3. **設定管理の統合**
   - 複数サブコマンドを `pummit config` 単一エントリに統合
   - デフォルト動作をエディター起動に設定

### 中期目標（1〜3ヶ月）- 堅牢性向上

1. **アーキテクチャ改善**
   - UI Adapter層に "Renderer" インターフェースを設けBubble Tea実装を内包
   - 設定サブシステムの統廃合完了
   - Git操作のlibgit2バインディング検討（Windows DLL配布方法検証）

2. **セキュリティ・品質強化**
   - SBOM出力（cyclonedx・syft）の導入
   - `govulncheck` をワークフローに統合
   - CI Artifacts retention policy（60日以内・内部限定）
   - 依存関係の脆弱性スキャン自動化

3. **可観測性の確立（opt-in設計）**
   - OpenTelemetry（PUMMIT_OTEL_ENDPOINT環境変数制御）
   - 構造化ログ（JSON形式）の導入
   - 起動時間・API レイテンシの計測

4. **クロスプラットフォーム対応**
   - Windows 用 pty ラッパ導入
   - PowerShell 互換性テスト自動化
   - ターミナル互換性マトリックスの整備

### 長期目標（戦略的発展）

1. **新機能の段階的実装**
   - 設定ファイルマイグレーション機能（最優先）
   - 統合CLI設定管理機能（`pummit config` 単一エントリ）
   - 診断機能（`pummit doctor` MVP）
   - ブランチ名からのプレフィックス自動提案
   - インタラクティブモード (`pummit interactive`)
   - スコープ入力支援
   - コミットメッセージテンプレート機能

2. **基本的な多言語サポート**
   - 日本語・英語の2言語対応 (l18nのような仕組みを検討する)
   - 設定ファイルでの言語切り替え
   - インタラクティブモードのシンプルな翻訳
   - すでに実装ずみな既存の機能の出力も全て翻訳
   - GitHub上での翻訳ファイル管理

3. **運用性・保守性の向上**
   - 自動更新機能（`pummit update`コマンド）
   - 設定スキーマのバージョン管理
   - 診断機能（`pummit doctor`）の充実
   - エラーレポート・ログ収集機能

4. **企業・チーム向け機能**
   - チーム設定テンプレート
   - 使用統計・分析ダッシュボード
   - 組織向けカスタム絵文字セット
   - CI/CD パイプライン統合ガイド

### 開発プロセス改善

1. **仕様駆動開発**
   - 機能追加前の仕様策定
   - レビュープロセス強化
   - 変更影響範囲の明確化

2. **品質保証**
   - 静的解析ツール導入
   - セキュリティ監査
   - パフォーマンスベンチマーク

3. **リリース管理**
   - セマンティックバージョニング
   - 自動リリース
   - 後方互換性保証

---

## まとめ

このアーキテクチャ仕様書は、Pummit プロジェクトの技術的な設計と実装の指針を定めるものです。プロフェッショナルレビューチームからのフィードバックを反映し、実用性と堅牢性を重視した設計となっています。

### 主要な設計原則

- **堅牢性優先**: マイグレーション安全性とエラー処理を最優先
- **段階的実装**: リスクを最小化した機能追加戦略
- **クロスプラットフォーム対応**: Windows/macOS/Linux での一貫した動作
- **観測可能性**: 問題の早期発見と診断機能の充実
- **ユーザー体験重視**: 直感的なインターフェースと多言語サポート

### 技術的ハイライト

1. **JSON→TOML設定マイグレーション**: 安全な移行戦略とロールバック機能
2. **インタラクティブモード**: `pummit interactive` サブコマンドによる統合UI
3. **多言語サポート**: i18n基盤による日本語・英語対応
4. **リスク管理**: 8つの主要リスクと具体的対策
5. **品質保証**: 包括的テスト戦略と継続的品質監視

### 実装優先順位

**緊急（1〜2週間）**:
- オフライン機能、基本テストスイート、PowerShell互換性対応

**短期（1〜3ヶ月）**:
- SBOM/脆弱性スキャン、UI Adapter パターン、可観測性基盤

**中長期（戦略的）**:
- 新機能段階実装、基本多言語サポート、企業向け機能

### 品質目標（再調整版）

- **起動時間**:
  - 基本動作: 100ms以下
  - OpenTelemetry有効時: 110ms以下
  - Bubble Tea UI使用時: 120ms以下
- **メモリ使用量**:
  - 基本動作: 10MB以下
  - 将来の多言語・絵文字データ拡張時: 15MB以下を目標
- **テストカバレッジ**: 80-85%（現実的な品質・保守性バランス）
- **セキュリティ**: SBOM生成・脆弱性スキャン・60日保持ポリシー
- **互換性**: 主要ターミナル・OS での動作保証（サポート表明記）

この仕様書は開発チームの共通理解を深め、プロダクションレディなCLIツールとしてのpummit発展を支援します。2回のプロフェッショナルレビューを経て、実用性・保守性・拡張性のバランスが取れた設計となっています。

### 次のアクションアイテム

**即座に着手すべき項目**:
1. TOMLライブラリ選定とパフォーマンス比較
2. 起動時間ベンチマークのCI統合
3. `pummit config` コマンドの統合設計
4. `pummit doctor` MVP実装

**重要な意思決定事項**:
- 設定管理機能の複雑化を避け、UXを最優先
- OpenTelemetryはopt-in設計でパフォーマンス影響を最小化
- テストカバレッジは品質と保守性のバランスを重視
- 廃止ロードマップの明確化（JSON設定のサンセット計画）

---

**作成日**: 2025年6月16日
**最終更新日**: 2025年6月17日
**バージョン**: 3.0
**作成者**: HidemaruOwO and Claude 4 Sonnet
**レビュー**: OpenAI o3 (2回完了)
**最終更新**: 再レビューフィードバック全面反映・設定管理統合・リスク11項目対応・現実的品質目標設定