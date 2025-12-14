# Pummit リファクタリング方針提案

調査の結果、コードベースには既にアーキテクチャ仕様書（`docs/architecture-v3.md`）で詳細なリファクタリング計画が策定されていることを確認しました。現在の実装状況を踏まえ、メンテナンス性向上のための優先度付きリファクタリング方針を提案します。

---

## 1. 優先度: 高 🔥

### 1.1 グローバル状態への依存削減

**問題点:**
- `internal/config/config.go:22-26` および `toml.go:90-93` でグローバル変数を使用
- 各パッケージが `config.CurrentTOMLConfig` に直接アクセス
- テストで並列実行ができない（`Do NOT use t.Parallel()` コメントが必要）

**提案:**
```go
// config/config.go に Config インターフェースを追加
type ConfigService interface {
    GetTOMLConfig() *TOMLConfig
    SaveTOMLConfig() error
}

// 依存性注入パターンへ移行
type GitService struct {
    config ConfigService
}
```

**メリット:**
- テスト容易性の向上
- 並列テスト可能に
- 責務の明確化

---

### 1.2 重複コードの統合

**ブランチ→絵文字推測ロジックの重複:**

| 場所 | 問題 |
|------|------|
| `internal/mcp/tool_git.go:138-160` | ハードコードされたパターン |
| `internal/config/toml.go:165-176` | 設定ファイルで定義 |

**提案:**
新規ファイル `internal/branch/suggestion.go` を作成し、設定ベースの推測機能を統一：

```go
// internal/branch/suggestion.go
package branch

import (
    "regexp"
    "github.com/HidemaruOwO/pummit/internal/config"
)

func SuggestEmoji(branchName string) string {
    for _, rule := range config.CurrentTOMLConfig.Branch.Rules {
        if matched, _ := regexp.MatchString(rule.Pattern, branchName); matched {
            return rule.Emoji
        }
    }
    return config.CurrentTOMLConfig.Branch.Fallback
}
```

**ステージファイル取得の重複:**

| 関数 | 場所 | 出力形式 |
|------|------|---------|
| `GetStagedFiles()` | `git.go:123-137` | `[]string` |
| `GetChangedFiles()` | `git.go:139-156` | カンマ区切り文字列 |

**提案:**
```go
// GetChangedFiles を GetStagedFiles の薄いラッパーに
func GetChangedFiles() (string, error) {
    files, err := GetStagedFiles()
    if err != nil {
        return "", err
    }
    if len(files) == 0 {
        log.Info("Nothing to changed files.")
        os.Exit(0)  // ← これも後述の問題
    }
    return strings.Join(files, ", "), nil
}
```

---

### 1.3 `os.Exit()` をビジネスロジック層から除去

**問題点:**
- `internal/git/git.go:56-57` および `152-153` で `os.Exit(1)` を呼び出し
- ビジネスロジック層でプログラム終了を制御
- テスト困難、エラーハンドリング不可

**提案:**
```go
// Before (git.go:49-57)
func commitWithOfflineModeInternal(...) error {
    changed, err := GetChangedFiles()
    if err != nil {
        log.Error(err.Error())
        os.Exit(1)  // 問題
    }
    ...
}

// After - エラーを返す
var ErrNoStagedFiles = errors.New("no staged files")

func GetChangedFiles() (string, error) {
    ...
    if changed == "" {
        return "", ErrNoStagedFiles  // エラーを返す
    }
    return strings.ReplaceAll(changed, "\n", ", "), nil
}

// CLI層で os.Exit を処理
func runRootCommand() {
    if err := git.Commit(cm); err != nil {
        if errors.Is(err, git.ErrNoStagedFiles) {
            log.Info("Nothing to commit.")
        } else {
            log.Error(err.Error())
        }
        os.Exit(1)
    }
}
```

---

## 2. 優先度: 中 🔷

### 2.1 未使用コードの削除

**削除候補:**

| ファイル | 関数/変数 | 理由 |
|---------|----------|------|
| `internal/config/compatibility.go:30-69` | `GetAlias()` | `alias.GetEmoji()` と重複、実質未使用 |
| `internal/config/config.go:23-26` | `DefaultConfig`, `CurrentConfig` | TOML移行完了で不要 |

### 2.2 MCPツールのビジネスロジック分離

**問題点:**
`internal/mcp/tool_git.go` に以下のビジネスロジックが混在：
- `suggestEmojiFromBranch()` (138-160行)
- `generateSmartCommitMessage()` (162-196行)

**提案:**
```
internal/
├── git/
│   ├── git.go           # 既存
│   ├── branch.go        # 新規: ブランチ関連
│   └── message.go       # 新規: メッセージ生成
├── mcp/
│   └── tool_git.go      # ハンドラのみ
```

### 2.3 設定パッケージの整理

**現状:**
```
config/
├── config.go           # JSON設定（レガシー）
├── toml.go            # TOML設定
├── migration.go       # マイグレーション
└── compatibility.go   # 互換性
```

**提案:**
```
config/
├── config.go          # 共通インターフェース + パス解決
├── toml.go            # TOML固有処理
├── migration.go       # マイグレーション（v4.0で削除予定）
└── json_legacy.go     # レガシーJSON（v4.0で削除予定）
```

---

## 3. 優先度: 低 🔵

### 3.1 テストカバレッジの拡充

**テストが不足しているパッケージ:**

| パッケージ | 優先度 | 理由 |
|-----------|--------|------|
| `internal/alias/` | 高 | エイリアスCRUDのコアロジック |
| `internal/emojis/` | 高 | オフライン/オンライン切り替え |
| `internal/mcp/` | 中 | MCPハンドラー |
| `internal/prompt/` | 低 | UI依存 |
| `internal/utils/` | 低 | 単純なユーティリティ |

### 3.2 パッケージ構造の最終整理（v4.0向け）

**将来の理想的な構造:**
```
internal/
├── cli/              # CLI層（Cobra）
│   ├── root.go
│   ├── alias/
│   ├── config/
│   └── migrate/
├── core/             # ビジネスロジック層
│   ├── git/          # Git操作
│   ├── alias/        # エイリアス管理
│   ├── emoji/        # 絵文字変換
│   └── branch/       # ブランチ解析
├── config/           # 設定管理
├── mcp/              # MCPサーバー（外部統合）
└── ui/               # TUI（Bubble Tea）
```

---

## 4. 実装優先順位のまとめ

| 順序 | タスク | 影響範囲 | 推定工数 |
|------|-------|---------|---------|
| 1 | `os.Exit()` の除去 | `git/git.go` | 1日 |
| 2 | ブランチ推測ロジックの統合 | `mcp/`, `branch/` 新規 | 1日 |
| 3 | ステージファイル取得の統合 | `git/git.go` | 0.5日 |
| 4 | 未使用コード削除 | `config/compatibility.go` | 0.5日 |
| 5 | MCPビジネスロジック分離 | `mcp/`, `git/` | 2日 |
| 6 | テスト追加（alias, emojis） | `*_test.go` | 3-5日 |

---

## 5. 質問・確認事項

1. **依存性注入パターンの導入について**: グローバル状態を完全に排除するのは大きな変更になります。段階的に進めるか、一括で実施するかどちらを希望しますか？

2. **v4.0でのJSON設定削除**: `docs/architecture-v3.md` によると v4.0 でJSON設定を完全削除予定ですが、この計画に沿って進めてよろしいですか？

3. **テスト優先度**: 現在のテストカバレッジが低い状況で、リファクタリングを先に進めるか、テスト基盤を先に整備するかどちらを優先しますか？

---

## 参考資料

- [アーキテクチャ仕様書](./architecture-v3.md)
- [TODO管理](./TODO-v3.md)

---

**作成日**: 2025年12月5日
**バージョン**: 1.0.0
