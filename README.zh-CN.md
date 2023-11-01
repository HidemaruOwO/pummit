<div align="center">

# 普米特🚛

[![Test CLI](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml/badge.svg)](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml)![最終コミット](https://img.shields.io/github/last-commit/HidemaruOwO/pummit?style=flat-square)![リポジトリのスター](https://img.shields.io/github/stars/HidemaruOwO/pummit?style=flat-square)![問題](https://img.shields.io/github/issues/HidemaruOwO/pummit?style=flat-square)![オープンな問題](https://img.shields.io/github/issues-raw/HidemaruOwO/pummit?style=flat-square)![バグの問題](https://img.shields.io/github/issues/HidemaruOwO/pummit/bug?style=flat-square)

![image](https://user-images.githubusercontent.com/82384920/225959857-76495875-c426-4669-a8d4-372ebf3acfad.png)

## 这是什么？

您可以轻松创建像这样的漂亮提交消息。

<img width="1441" alt="image" src="https://github.com/HidemaruOwO/pummit/assets/82384920/8461400a-94f6-431d-99d4-32ae74afe7fd">

</div>

-   选择语言

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

## 如何使用💨

pummit 有两种使用方式

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

## 安装😊

如果您安装了 Go，请运行它。

```bash
go install github.com/HidemaruOwO/pummit/pummit@latest
```

[HTTPS://GitHub.com/hide MA入O五O/朴MM IT/releases](https://github.com/HidemaruOwO/pummit/releases)

如果未安装，请从Release下载与您的环境相匹配的文件，然后运行以下命令。

```bash
tar xzvf pummit**.tar.gz
sudo mv pummit /usr/local/bin
```

## 构建🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
mkdir build && cd build
go build ../pummit/
```

## 依赖关系🪡

要使用pummit，请在您的路径中注册以下命令

-   git

## 与lazygit一起使用🔍

以下键绑定`.config/lazygit/config.yml`请设置为

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

### 表情符号前缀示例 🌟

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

## 关于别名功能📎

例如`wastebasket`输入起来有点困难，但是如果使用alias功能就可以`wb`您可以轻松输入。

```bash
$ pummit wb モジュールの削除
# Result: :wastebasket: モジュールの削除 (path/to/added/file)
```

默认别名如下。

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

### 添加命令

此命令允许您添加别名。

```bash
$ pummit alias add 's' 'sparkles'
```

在这种情况下`s`只需在提交消息中输入别名“表情符号前缀”即可。`sparkles`您现在可以替换 .

```bash
$ pummit s 新機能の追加
# Run: git commit -m ':sparkles: 新機能の追加 (path/to/added/file)'
```

### 删除命令

此命令允许您删除别名。

```bash
$ pummit alias delete s
```

在这种情况下，`s=spakles`如果您运行此命令并假设别名已注册`s`和`sparkles`由于 Emoji 前缀的关联将会丢失，因此即使运行以下命令，Emoji 前缀也不会被关联`s`仅被分配。

```bash
$ pummit s 新機能の追加
# Run: git commmit -m ':s: 新機能の追加 (path/to/added/file)'
```

您还可以指定多个要删除的别名作为参数。

```bash
$ pummit alias delete s sm c h
```

### 删除 --all 命令

此命令删除所有已注册的别名。

```bash
$ pummit alias delete --all
```

### 列表命令

此命令显示所有已注册的别名。

```bash
$ pummit alias list
```

如果别名`s=sparkles`和`t=tada`如果已注册，将输出以下内容。

```bash
📎 There is aliases
Alias : Prefix : Emoji
  s : sparkles : ✨
  t : tada : 🎉
```

### 复位命令

此命令重置别名。

```bash
$ pummit alias reset
```

如果别名太多，会变得混乱，`config.json`如果您直接弄乱它而导致错误，它可以用作恢复工具。

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

即使有很多像这样令人困惑的别名：

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

您只需一个命令即可将其恢复到这种美丽的状态。

## 特别感谢✨

-   [Qiita - 在 GitHub 提交消息中添加表情符号以提高开发效率](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
