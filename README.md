<div align="center">

# pummit 🚛

[![Test CLI](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml/badge.svg)](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml)
![最終コミット](https://img.shields.io/github/last-commit/HidemaruOwO/pummit?style=flat-square)
![リポジトリのスター](https://img.shields.io/github/stars/HidemaruOwO/pummit?style=flat-square)
![問題](https://img.shields.io/github/issues/HidemaruOwO/pummit?style=flat-square)
![オープンな問題](https://img.shields.io/github/issues-raw/HidemaruOwO/pummit?style=flat-square)
![バグの問題](https://img.shields.io/github/issues/HidemaruOwO/pummit/bug?style=flat-square)

![image](https://user-images.githubusercontent.com/82384920/225959857-76495875-c426-4669-a8d4-372ebf3acfad.png)

## what is this?

It allows you to lovely create a nicely commit message like this

![image](https://user-images.githubusercontent.com/82384920/225978215-9ac68cd4-cdb0-44c9-bca3-4d2cff1896cf.png)

</div>

<table>
  <thead>
    <tr>
      <th style="text-align:center">🤡English (Translated by ChatGPT)</a></th>
      <th style="text-align:center"><a href="README.ja.md">🎌日本語</a></th>
    </tr>
  </thead>
</table>

## Usage 💨

You can use pummit in two ways.

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

## Install 😊

If Go is installed, please run this.

```bash
go install github.com/HidemaruOwO/pummit/pummit@latest
```

https://github.com/HidemaruOwO/pummit/releases

If Go is not installed, download the appropriate file for your environment from
Release and execute the following command.

```bash
tar xzvf pummit**.tar.gz
sudo mv pummit /usr/local/bin
```

## Build 🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
mkdir build && cd build
go build ../pummit/
```

## dependence 🪡

To use pummit, please register the following command in your path

- git

## How to use on lazygit 🔍

Set the following key bindings in `.config/lazygit/config.yml`

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

### Sample emoji prefixes 🌟

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

## About Alias Function 📎

For example, typing `wastebasket` is a bit of a challenge, but the alias feature
makes it easy to type `wb`.

```bash
$ pummit wb 'Delete module'
# Result: :wastebasket: Delete module (path/to/added/file)
```

The aliases set by default are as follows

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

This command allows aliases to be added.

```bash
$ pummit alias add 's' 'sparkles'
```

In this case, by simply entering the alias `s`, you can assign `sparkles` to the
Emoji prefix of the commit message.

```bash
$ pummit s 'Add new feature'
# Run: git commit -m ':sparkles: Add new feature (path/to/added/file)'
```

### Delete command

This command allows you to delete an alias.

```bash
$ pummit alias delete s
```

In this case, if the `s=spakles` alias is registered and this command is
executed, the association between `s` and `sparkles` will be lost. Therefore, if
you run the following command, only `s` will be assigned to the Emoji prefix.

```bash
$ pummit s 'Add new feature'
# Run: git commmit -m ':s: Add new feature (path/to/added/file)'
```

You can also specify multiple aliases to delete as arguments.

```bash
$ pummit alias delete s sm c h
```

### List command

This command displays all registered aliases.

```bash
$ pummit alias list
```

f aliases such as `s=sparkles` and `t=tada` are registered, the output will be
as follows:

```bash
📎 There is aliases
Alias : Prefix : Emoji
  s : sparkles : ✨ 
  t : tada : 🎉
```

### Reset command

This command resets the aliases.

```bash
$ pummit alias reset
```

It can be used for recovery when there are too many confusing aliases or when
you have directly tampered with `~/.config/pummit/config.json` and caused a bug.

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

Even when there are confusingly many aliases like this one

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

With a single command, it can be returned to this beautiful state.

## Special Thanks ✨

- [Qiita - GitHubのコミットメッセージに絵文字を入れて開発効率をあげる](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
