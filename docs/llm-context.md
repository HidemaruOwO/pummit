# Pummit プロジェクトコンテキスト

## 1. プロジェクト概要

### 1.1. 目的と機能
Pummitは、Gitコミットメッセージを絵文字付きで美しく、一貫性のあるものにするためのCLIツールです。主な目的は、開発者の生産性向上と、Conventional CommitsやGitmojiの標準に準拠した意味のあるコミットメッセージ作成を支援することです。

### 1.2. 主要な特徴
- **絵文字付きコミット**: 絵文字プレフィックスによる視覚的なコミットメッセージ。
- **スマートな自動化**: ブランチ名に基づく絵文字の自動提案 (`feature/` -> ✨)。
- **インタラクティブモード**: 対話形式でのコミット作成支援。
- **強力な設定**: TOMLベースの設定ファイルによる柔軟なカスタマイズ。
- **エイリアスシステム**: 絵文字やプレフィックスの短縮形を定義可能。
- **多言語対応**: UIは英語と日本語をサポート。
- **クロスプラットフォーム**: Windows, macOS, Linuxで動作。

### 1.3. バージョン情報とライセンス
- **Goバージョン**: 1.23.0 以上
- **ライセンス**: Apache License 2.0, SUSHI-WARE

---

## 2. アーキテクチャ分析

### 2.1. ディレクトリ構造
プロジェクトは標準的なGoのレイアウトに従っています。
```
pummit/
├── main.go               # エントリーポイント
├── internal/             # 内部実装（非公開）
│   ├── cli/              # CLIコマンド (Cobra)
│   ├── config/           # 設定管理
│   ├── git/              # Git操作
│   ├── alias/            # エイリアス管理
│   └── ...
├── pkg/                  # 公開可能なライブラリ
│   ├── gitmoji/          # Gitmoji API連携
│   └── logger/           # ロガー
└── docs/                 # ドキュメント
```

### 2.2. レイヤー化アーキテクチャ
`docs/architecture-v3.md`で定義されている通り、4層のレイヤー化アーキテクチャを採用しています。
1.  **User Interface Layer**: CLIコマンドの解析とインタラクティブプロンプトを担当 (`internal/cli`, `internal/prompt`)。
2.  **Application Layer**: ビジネスロジックの調整と制御。
3.  **Business Logic Layer**: Git操作、エイリアス管理、設定管理などのコアロジック (`internal/git`, `internal/alias`, `internal/config`)。
4.  **Infrastructure Layer**: ファイルシステム、ネットワーク通信、ロギングなど外部システムとの連携 (`os`, `net/http`, `pkg/logger`)。

---

## 3. コードベース分析

### 3.1. 主要な型とインターフェース
- `config.Config`: アプリケーションの設定を保持する構造体。JSON/TOMLからのデシリアライズを想定。
- `git.CommitMessage`: コミット情報を保持する構造体（絵文字とメッセージ）。
- `alias.Alias`: エイリアスの定義を保持する構造体。

### 3.2. コア機能の実装パターン
コミット作成処理 (`git.CommitWithOfflineMode`) は以下の手順で実行されます。
1.  `git diff --name-only --cached`でステージングされたファイルを取得。
2.  入力された絵文字/エイリアスを`alias.GetEmoji()`で解決。
3.  `config.CurrentConfig`に基づき、生の絵文字 (`✨`) か名前 (`:sparkles:`) かを決定。
4.  `emojis.GetEmojiByNameOffline()`で最終的な絵文字を取得（オフライン対応）。
5.  `fmt.Sprintf`でコミットメッセージを整形。
6.  `exec.Command("git", "commit", "-m", ...)`でGitコマンドを実行。

### 3.3. 設定管理システム
- `internal/config/config.go`が中心的な役割を担います。
- `Init()`: 起動時に呼び出され、デフォルト設定の読み込み、設定パスの解決、`AutoMigrate()`の実行を行います。
- `Load()` / `Save()`: `config.json`の読み書きを行います。
- `internal/config/migration.go`: 古いJSON形式から新しいTOML形式への自動マイグレーションロジックを実装しています。

### 3.4. エラーハンドリング戦略
- `os/exec`やファイルI/Oでエラーが発生した場合、`log.Error()`でエラーメッセージを出力し、`os.Exit(1)`でプロセスを終了するのが基本パターンです。
- `docs/architecture-v3.md`には詳細なエラー分類が記載されていますが、現在の実装はよりシンプルです。

---

## 4. 実装済み機能

- **CLIコマンド**: `spf13/cobra`を利用。
    - `pummit [emoji] [message]`: 基本的なコミットコマンド。
    - `pummit alias [add|list|delete|reset]`: エイリアス管理。
    - `pummit migrate`: 設定ファイルのマイグレーション。
    - `pummit --version`, `pummit --offline`
- **設定システム**: 当初はJSON、現在はTOMLへ移行中。`~/.config/pummit/` (Unix) または `%APPDATA%\pummit` (Windows) に保存。
- **エイリアス管理**: 複数の名前（例: `s`, `feat`）を一つの絵文字にマッピング可能。
- **Gitmoji連携**: オンライン時はAPIを叩き、オフライン時はローカルのキャッシュ (`discord-emojis.flat.json`) を利用するフォールバック機能があります。

---

## 5. 技術スタック

- **言語**: Go 1.23.0
- **主要ライブラリ**:
    - `github.com/spf13/cobra`: CLIフレームワーク
    - `github.com/charmbracelet/bubbletea`: TUIフレームワーク（インタラクティブモード用）
    - `github.com/fatih/color`: コンソール出力のカラーリング
    - `github.com/BurntSushi/toml`: TOMLファイルの解析
- **ビルド・配布**: `go build`による手動ビルド、または`.goreleaser.yaml`を用いたリリース自動化。

---

## 6. 開発パターンと慣例

- **コーディング規約**: `gofmt`による標準フォーマットが適用されています。
- **ファイル命名規則**: `snake_case.go`形式で統一されています。
- **コメント**: 日本語でのコメントが多く、コードの意図を補足しています。
- **テスト戦略**: `README.md`には`go test ./...`の記述がありますが、現状のコードベースにはテストファイル (`*_test.go`) がほとんど存在せず、テストカバレッジは非常に低いです。

---

## 7. 今後の開発に重要な知見

- **技術的負債**:
    1.  **テストコードの欠如**: 最大の技術的負債。リファクタリングや機能追加を安全に行うための回帰テストが存在しません。
    2.  **設定マイグレーションの複雑性**: `config/compatibility.go`や`config/migration.go`にJSON/TOMLの過渡期対応コードがあり、将来的に削除する必要があります。
- **拡張時の注意点**:
    - `internal`パッケージに主要ロジックが密結合しているため、機能追加の際は既存モジュールへの影響範囲を慎重に調査する必要があります。
    - `os/exec`で直接Gitコマンドを叩いているため、複雑なGit操作を追加する際は、コマンドのインジェクションやOS間の挙動差異に注意が必要です。
- **パフォーマンス考慮事項**:
    - 現状の`os/exec`の利用はシンプルですが、多数のGit操作を行う機能（例: 履歴の分析）を追加する場合、パフォーマンスが問題になる可能性があります。その際は`go-git`のようなライブラリの導入を検討する必要があります。
- **クロスプラットフォーム対応**:
    - `config.GetConfigDir()`で示されているように、ファイルパスや外部コマンドの扱いはOS間の差異を意識する必要があります。