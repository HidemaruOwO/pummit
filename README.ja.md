<div align="center">

# pummit 🚛
[![Test CLI](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml/badge.svg)](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml)
![最終コミット](https://img.shields.io/github/last-commit/HidemaruOwO/pummit?style=flat-square)
![リポジトリのスター](https://img.shields.io/github/stars/HidemaruOwO/pummit?style=flat-square)
![問題](https://img.shields.io/github/issues/HidemaruOwO/pummit?style=flat-square)
![オープンな問題](https://img.shields.io/github/issues-raw/HidemaruOwO/pummit?style=flat-square)
![バグの問題](https://img.shields.io/github/issues/HidemaruOwO/pummit/bug?style=flat-square)

![image](https://user-images.githubusercontent.com/82384920/225959857-76495875-c426-4669-a8d4-372ebf3acfad.png)

## なんだこれは

このような綺麗な形のコミットメッセージを簡単に作成出来ます
  
![image](https://user-images.githubusercontent.com/82384920/225978215-9ac68cd4-cdb0-44c9-bca3-4d2cff1896cf.png)

</div>

<table>
  <thead>
    <tr>
      <th style="text-align:center"><a href="README.md">🤡English</a></th>
      <th style="text-align:center">🎌日本語</a></th>
    </tr>
  </thead>
</table>

## 使い方 💨

pummitは２つの方法で使用することができます

```bash
pummit emojiprefix 'subject'
# or
pummit 'emojiprefix subject'

# Example
pummit sparkles 'I am unko man'
# Run: git commit -m ':sparkles: I am unko man (path/to/added/file, path/to/added/file)'

pummit 'sparkles I am unko man'
# Run: git commit -m ':sparkles: I am unko man (path/to/added/file, path/to/added/file)'
```

### 絵文字プレフィックスのサンプル 🌟

```
# ==================== Emojis ====================
# 🎉  :tada: 初めてのコミット（Initial Commit）
# ♻️   :recycle: マージ(Merge)
# 🔖  :bookmark: バージョンタグ（Version Tag）
# ✨  :sparkles: 新機能（New Feature）
# 🐛  :bug: バグ修正（Bagfix）
# 👀  :eyes: リファクタリング(Refactoring)
# 📚  :books: ドキュメント（Documentation）
# 🎨  :art: デザインUI/UX(Accessibility)
# 🐎  :horse: パフォーマンス（Performance）
# 🔧  :wrench: ツール（Tooling）
# 🚨  :rotating_light: テスト（Tests）
# 💩  :hankey: 非推奨追加（Deprecation）
# 🗑️  :wastebasket: 削除（Removal）
# 🚧  :construction: WIP(Work In Progress)
# ☃️  :snowman: 仕様変更

# ==================== Format ====================
# :emoji: Subject (Dir/AddedFile Dir/AddedFile)
#
# Commit body...
```

## ビルド 🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
mkdir build && cd build
go build ../pummit/pummit.go
```

## 依存関係 🪡

pummitを使うには以下のコマンドをパスに通してください

- git

## lazygitで使うには 🔍

以下のキーバインドを`.config/lazygit/config.yml`に設定してください

```yaml
customCommands:
- key: 'c'
  prompts:
    - type: 'input'
      title: 'Commit message'
      initialValue: ''
  command: "pummit '{{index .PromptResponses 0}}'"
  context: 'files'
  description: 'commit changes(Custom Command)'
 ```

## スペシャルサンクス ✨
- [Qiita - GitHubのコミットメッセージに絵文字を入れて開発効率をあげる](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
