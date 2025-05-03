```
pummit/
├── cmd/                      # エントリーポイント
│   └── pummit/               # メインCLIアプリケーション
│       └── main.go           # エントリーポイント
├── internal/                 # 内部実装（外部からインポートされない）
│   ├── config/               # 設定関連
│   │   └── config.go
│   ├── alias/                # エイリアス機能
│   │   ├── alias.go
│   │   ├── handler.go
│   │   └── storage.go
│   ├── emoji/                # 絵文字関連
│   │   └── emoji.go
│   ├── git/                  # Git操作関連
│   │   └── git.go
│   ├── utils/                # 共通ユーティリティ
│   │   ├── error.go
│   │   └── slice.go
│   └── cli/                  # CLIコマンド実装
│       ├── root.go
│       ├── version.go
│       └── alias/
│           ├── add.go
│           ├── list.go
│           ├── delete.go
│           └── reset.go
├── pkg/                      # 外部からも使用可能な機能
│   ├── gitmoji/              # Gitmoji関連
│   │   └── gitmoji.go
│   └── logger/               # ロガー
│       └── logger.go
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```
