<div align="center">

# 普米特🚛

[![Test CLI](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml/badge.svg)](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml)![最終コミット](https://img.shields.io/github/last-commit/HidemaruOwO/pummit?style=flat-square)![リポジトリのスター](https://img.shields.io/github/stars/HidemaruOwO/pummit?style=flat-square)![問題](https://img.shields.io/github/issues/HidemaruOwO/pummit?style=flat-square)![オープンな問題](https://img.shields.io/github/issues-raw/HidemaruOwO/pummit?style=flat-square)![バグの問題](https://img.shields.io/github/issues/HidemaruOwO/pummit/bug?style=flat-square)

![image](https://user-images.githubusercontent.com/82384920/225959857-76495875-c426-4669-a8d4-372ebf3acfad.png)

## 這是什麼？

您可以輕鬆創建像這樣的漂亮提交訊息。

<img width="1441" alt="image" src="https://github.com/HidemaruOwO/pummit/assets/82384920/8461400a-94f6-431d-99d4-32ae74afe7fd">

</div>

-   選擇語言

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

pummit 有兩種使用方式

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

## 安裝😊

如果您安裝了 Go，請運行它。

```bash
go install github.com/HidemaruOwO/pummit/pummit@latest
```

<https://github.com/HidemaruOwO/pummit/releases>

如果未安裝，請從Release下載與您的環境相符的文件，然後執行下列命令。

```bash
tar xzvf pummit**.tar.gz
sudo mv pummit /usr/local/bin
```

## 建構🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
mkdir build && cd build
go build ../pummit/
```

## 依賴關係🪡

若要使用pummit，請在您的路徑中註冊以下命令

-   git

## 與lazygit一起使用🔍

以下鍵綁定`.config/lazygit/config.yml`請設定為

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

### 表情符號前綴範例 🌟

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

## 關於別名功能📎

例如`wastebasket`輸入起來有點困難，但是如果使用alias功能就可以`wb`您將能夠輕鬆輸入它。

```bash
$ pummit wb モジュールの削除
# Result: :wastebasket: モジュールの削除 (path/to/added/file)
```

預設別名如下。

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

### 新增命令

此命令允許您新增別名。

```bash
$ pummit alias add 's' 'sparkles'
```

在這種情況下`s`只需在提交訊息中輸入別名“表情符號前綴”即可。`sparkles`將能夠被替代。

```bash
$ pummit s 新機能の追加
# Run: git commit -m ':sparkles: 新機能の追加 (path/to/added/file)'
```

### 刪除命令

此命令允許您刪除別名。

```bash
$ pummit alias delete s
```

在這種情況下，`s=spakles`如果您執行此命令並假設別名已註冊`s`和`sparkles`由於 Emoji 前綴的關聯將會遺失，因此即使執行以下命令，Emoji 前綴也不會被關聯`s`しか代入されません。

```bash
$ pummit s 新機能の追加
# Run: git commmit -m ':s: 新機能の追加 (path/to/added/file)'
```

您也可以指定多個要刪除的別名作為參數。

```bash
$ pummit alias delete s sm c h
```

### 刪除 --all 命令

此命令刪除所有已註冊的別名。

```bash
$ pummit alias delete --all
```

### 列表命令

此指令顯示所有已註冊的別名。

```bash
$ pummit alias list
```

如果別名`s=sparkles`和`t=tada`如果已註冊，將輸出以下內容。

```bash
📎 There is aliases
Alias : Prefix : Emoji
  s : sparkles : ✨
  t : tada : 🎉
```

### 重設指令

此命令重置別名。

```bash
$ pummit alias reset
```

如果別名太多，會變得混亂，`config.json`如果您直接使用它而導致錯誤，它可以用作恢復工具。

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

即使有很多像這樣令人困惑的別名：

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

您只需一個命令即可將其恢復到這種美麗的狀態。

## 特別感謝✨

-   [Qiita - 在 GitHub 提交訊息中加入表情符號以提高開發效率](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
