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

綺麗な形のコミットメッセージを簡単に作成出来ます

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
pummit sparkles 'I am unko man'
# Run: git commit -m ':sparkles: I am unko man (path/to/added/file, path/to/added/file)'

pummit 'sparkles I am unko man'
# Run: git commit -m ':sparkles: I am unko man (path/to/added/file, path/to/added/file)'
```

## ビルド 🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
go build src/pummit.go
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
