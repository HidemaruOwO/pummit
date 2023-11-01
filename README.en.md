<div align="center">

# pummit 🚛

[![Test CLI](https://github.com/HidemaruOwO/pummit/actions/workflows/test.yml/badge.svg)](https://github.com/HidemaruOwO/pummit/actions/workflows/test.yml)![最終コミット](https://img.shields.io/github/last-commit/HidemaruOwO/pummit?style=flat-square)![リポジトリのスター](https://img.shields.io/github/stars/HidemaruOwO/pummit?style=flat-square)![問題](https://img.shields.io/github/issues/HidemaruOwO/pummit?style=flat-square)![オープンな問題](https://img.shields.io/github/issues-raw/HidemaruOwO/pummit?style=flat-square)![バグの問題](https://img.shields.io/github/issues/HidemaruOwO/pummit/bug?style=flat-square)

![image](https://user-images.githubusercontent.com/82384920/225959857-76495875-c426-4669-a8d4-372ebf3acfad.png)

## What is this?

You can easily create beautiful commit messages like this one.

<img width="1441" alt="image" src="https://github.com/HidemaruOwO/pummit/assets/82384920/8461400a-94f6-431d-99d4-32ae74afe7fd">

</div>

-   Select Language

<table>
  <thead>
    <tr>
      <th style="text-align:center"><a href="README.md">🎌日本語</a></th>
      <th style="text-align:center"><a href="README.en.md">🤡English</a></th>
      <th style="text-align:center"><a href="README.zh-CN.md">🐉简体中文</a></th>
      <th style="text-align:center"><a href="README.zh-TW.md">🍜繁体中文</a></th>
      <th style="text-align:center"><a href="README.ko.md">🌸한국어</a></th>
    </tr>
  </thead>
</table>

## How to use 💨

pummit can be used in two ways

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

If you have Go installed, run this.

```bash
go install github.com/HidemaruOwO/pummit/pummit@latest
```

<https://github.com/HidemaruOwO/pummit/releases>

If it is not installed, download the file suitable for your environment from Release and execute the following command.

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

## Dependencies 🪡

To use pummit, please register the following command in your path

-   git

## To use with lazygit 🔍

the following keybindings`.config/lazygit/config.yml`Please set it to

```yaml
customCommands:
  - key: "c"
    prompts:
      - type: "input"
        title: "Commit message"
        initialValue: ""
    command: "pummit '{{index .PromptResponses 0}}'"
    context: "files"
    description: "commit changes(Custom Command)"
```

### Sample emoji prefix 🌟

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

## About the alias function 📎

for example`wastebasket`It's a little difficult to enter, but if you use the alias function, you can`wb`You can enter it easily.

```bash
$ pummit wb モジュールの削除
# Result: :wastebasket: モジュールの削除 (path/to/added/file)
```

The default aliases are as follows.

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

### Add command

This command allows you to add an alias.

```bash
$ pummit alias add 's' 'sparkles'
```

In this case`s`というエイリアスを入力するだけでコミットメッセージのEmoji prefixに`sparkles`will be able to be substituted.

```bash
$ pummit s 新機能の追加
# Run: git commit -m ':sparkles: 新機能の追加 (path/to/added/file)'
```

### Delete command

This command allows you to delete an alias.

```bash
$ pummit alias delete s
```

In this case,`s=spakles`If you run this command assuming that the alias is registered`s`and`sparkles`Since the association of Emoji prefix will be lost, even if you run the following command, the Emoji prefix will not be`s`only is assigned.

```bash
$ pummit s 新機能の追加
# Run: git commmit -m ':s: 新機能の追加 (path/to/added/file)'
```

You can also specify multiple aliases to delete as arguments.

```bash
$ pummit alias delete s sm c h
```

### Delete --all command

This command deletes all registered aliases.

```bash
$ pummit alias delete --all
```

### List command

This command displays all registered aliases.

```bash
$ pummit alias list
```

If the alias`s=sparkles`と`t=tada`If it is registered, the following will be output.

```bash
📎 There is aliases
Alias : Prefix : Emoji
  s : sparkles : ✨
  t : tada : 🎉
```

### Reset command

This command resets the alias.

```bash
$ pummit alias reset
```

If there are so many aliases that it becomes confusing,`config.json`It can be used as a recovery tool in case you cause a bug by directly playing with it.

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

You can restore it to this beautiful state with just one command.

## Special thanks ✨

-   [Qiita - Add emojis to GitHub commit messages to improve development efficiency](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
