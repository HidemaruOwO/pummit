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

## インストール 😊
Goがインストールされている場合はこちらを実行してください。
```bash
go install github.com/HidemaruOwO/pummit/pummit@latest
```

https://github.com/HidemaruOwO/pummit/releases
 
インストールされていない場合はReleaseから環境にあったファイルをダウンロードして、以下のコマンドを実行してください。
```bash
tar xzvf pummit**.tar.gz
sudo mv pummit /usr/local/bin
```

## ビルド 🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
mkdir build && cd build
go build ../pummit/pummit.go
```

## 依存関係 🪡

pummitを使うには以下のコマンドをパスに登録してください

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

## エイリアス機能について 📎
例えば`wastebasket`を入力するのは少し大変ですが、エイリアス機能を使うと`wb`で簡単に入力できるようになります。
```bash
$ pummit wb モジュールの削除
# Result: :wastebasket: モジュールの削除 (path/to/added/file)
```
デフォルトで設定されているエイリアスは以下の通りです。
```
 📎 There is aliases
Alias : Prefix : Emoji
----------------------
  sm : snowman : ☃️
  h : horse : 🐎
  w : wrench : 🔧
  l : rotating_light : 🚨
  p : hankey : 💩
  wb : wastebasket : 🗑️
  c : construction : 🚧
  r : recycle : ♻️
  s : sparkles : ✨
  t : tada : 🎉
  e : eyes : 👀
  b : bug : 🐛
  d : books : 📚
  a : art : 🎨
```
### Add command
このコマンドはエイリアスを追加することが出来ます。
```bash
$ pummit alias add 's' 'sparkles'
```
この場合では`s`というエイリアスを入力するだけでコミットメッセージのEmoji prefixに`sparkles`を代入できるようになります。
```bash
$ pummit s 新機能の追加
# Run: git commit -m ':sparkles: 新機能の追加 (path/to/added/file)'
```
### Delete command
このコマンドはエイリアスを削除することが出来ます。
```bash
$ pummit alias delete s
```
この場合では、`s=spakles`というエイリアスが登録されている前提でこのコマンドを実行した場合`s`と`sparkles`の関連付けがなくなるため、以下のコマンドを実行してもEmoji prefixには`s`しか代入されません。
```bash
$ pummit s 新機能の追加
# Run: git commmit -m ':s: 新機能の追加 (path/to/added/file)'
```
また引数に削除したいエイリアスを複数指定することが出来ます。
```bash
$ pummit alias delete s sm c h
```

### Delete --all command
このコマンドは登録されているエイリアスを全て削除します。
```bash
$ pummit alias delete --all
```
### List command
このコマンドは登録されているエイリアスを全て表示します。
```bash
$ pummit alias list
```
もし、エイリアスに`s=sparkles`と`t=tada`が登録されている場合は以下のように出力されます。
```bash
📎 There is aliases
Alias : Prefix : Emoji
  s : sparkles : ✨ 
  t : tada : 🎉 
```
### Reset command
このコマンドをエイリアスをリセットします。
```bash
$ pummit alias reset
```
もし、エイリアスがこのように沢山あって混乱するほどあったり、`config.json`を直接弄ってバグらせてしまったときのリカバリとして使うことが出来ます。
```bash
$ pummit alias list
 📎 There is aliases
Alias : Prefix : Emoji
----------------------
  hjjciiiisadsadasda : sparkles : ✨
  w : wrench : 🔧
  s : sparkles : ✨
  l : rotating_light : 🚨
  p : hankey : 💩
  wb : wastebasket : 🗑️
  c : construction : 🚧
  sm : snowman : ☃️
  hj : sparkles : ✨
  hjjjksda : sparkles : ✨
  hjjca : sparkles : ✨
  hjjciiiia : sparkles : ✨
  a : art : 🎨
  h : horse : 🐎
  r : recycle : ♻️
  t : tada : 🎉
  b : bug : 🐛
  e : eyes : 👀
  d : books : 📚
 ```
  このようにエイリアスが混乱するほどある場合でも
```bash
$ pummit alias reset
> May I reset the aliases? :(Y/n) y
[INFO] Alias reseted

 📎 There is aliases
Alias : Prefix : Emoji
----------------------
  sm : snowman : ☃️
  h : horse : 🐎
  w : wrench : 🔧
  l : rotating_light : 🚨
  p : hankey : 💩
  wb : wastebasket : 🗑️
  c : construction : 🚧
  r : recycle : ♻️
  s : sparkles : ✨
  t : tada : 🎉
  e : eyes : 👀
  b : bug : 🐛
  d : books : 📚
  a : art : 🎨
```
コマンド一つでこのような綺麗な状態に戻せます。


## スペシャルサンクス ✨
- [Qiita - GitHubのコミットメッセージに絵文字を入れて開発効率をあげる](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
