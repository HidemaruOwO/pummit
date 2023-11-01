<div align="center">

# 푸읭 🚛

[![Test CLI](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml/badge.svg)](https://github.com/HidemaruOwO/pummit/actions/workflows/build-test.yml)![最終コミット](https://img.shields.io/github/last-commit/HidemaruOwO/pummit?style=flat-square)![リポジトリのスター](https://img.shields.io/github/stars/HidemaruOwO/pummit?style=flat-square)![問題](https://img.shields.io/github/issues/HidemaruOwO/pummit?style=flat-square)![オープンな問題](https://img.shields.io/github/issues-raw/HidemaruOwO/pummit?style=flat-square)![バグの問題](https://img.shields.io/github/issues/HidemaruOwO/pummit/bug?style=flat-square)

![image](https://user-images.githubusercontent.com/82384920/225959857-76495875-c426-4669-a8d4-372ebf3acfad.png)

## 이게 뭐야

이러한 아름다운 형태의 커밋 메시지를 쉽게 만들 수 있습니다.

<img width="1441" alt="image" src="https://github.com/HidemaruOwO/pummit/assets/82384920/8461400a-94f6-431d-99d4-32ae74afe7fd">

</div>

-   언어 선택

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

## 사용법 💨

pummit은 두 가지 방법으로 사용할 수 있습니다.

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

## 설치 😊

Go가 설치된 경우 여기를 실행합니다.

```bash
go install github.com/HidemaruOwO/pummit/pummit@latest
```

<https://github.com/HidemaruOwO/pummit/releases>

설치되어 있지 않은 경우 Release에서 환경에 있던 파일을 다운로드하고 다음 명령을 실행하십시오.

```bash
tar xzvf pummit**.tar.gz
sudo mv pummit /usr/local/bin
```

## 빌드 🔨

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
mkdir build && cd build
go build ../pummit/
```

## 종속성 🪡

pummit을 사용하려면 다음 명령을 경로에 등록하십시오.

-   자식

## lazygit에서 사용하려면 🔍

다음 키 바인딩`.config/lazygit/config.yml`로 설정

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

### 이모티콘 접두사 샘플 🌟

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

## 별칭 기능 정보 📎

예를 들면`wastebasket`입력하는 것은 조금 어렵지만 별칭 기능을 사용하면`wb`에서 쉽게 입력할 수 있습니다.

```bash
$ pummit wb モジュールの削除
# Result: :wastebasket: モジュールの削除 (path/to/added/file)
```

기본적으로 설정된 별칭은 다음과 같습니다.

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

### 명령 추가

이 명령은 별칭을 추가할 수 있습니다.

```bash
$ pummit alias add 's' 'sparkles'
```

이 경우에는`s`별칭을 입력하면 커밋 메시지의 Emoji prefix에`sparkles`을 대입할 수 있게 됩니다.

```bash
$ pummit s 新機能の追加
# Run: git commit -m ':sparkles: 新機能の追加 (path/to/added/file)'
```

### 삭제 명령

이 명령은 별칭을 삭제할 수 있습니다.

```bash
$ pummit alias delete s
```

이 경우에는`s=spakles`별칭이 등록되어 있는 전제로 이 명령을 실행한 경우`s`그리고`sparkles`연결이 없기 때문에 다음 명령을 실행해도 Emoji prefix`s`しか代入されません。

```bash
$ pummit s 新機能の追加
# Run: git commmit -m ':s: 新機能の追加 (path/to/added/file)'
```

또한 인수에 삭제하고 싶은 별칭을 복수 지정할 수 있습니다.

```bash
$ pummit alias delete s sm c h
```

### --all 명령 삭제

이 명령은 등록된 별칭을 모두 삭제합니다.

```bash
$ pummit alias delete --all
```

### 목록 명령

이 명령은 등록된 별칭을 모두 표시합니다.

```bash
$ pummit alias list
```

만약 별칭으로`s=sparkles`그리고`t=tada`가 등록되어 있는 경우는 다음과 같이 출력됩니다.

```bash
📎 There is aliases
Alias : Prefix : Emoji
  s : sparkles : ✨
  t : tada : 🎉
```

### 재설정 명령

이 명령은 별칭을 재설정합니다.

```bash
$ pummit alias reset
```

만약 별칭이 이렇게 많아서 혼란스러울 정도나`config.json`를 직접 Fuck하고 버그 시켜 버렸을 때의 복구로서 사용할 수 있습니다.

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

이와 같이 별칭이 혼란스러울 정도라도

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

커맨드 하나로 이러한 깨끗한 상태로 되돌릴 수 있습니다.

## 스페셜 선크스 ✨

-   [Qiita - GitHub 커밋 메시지에 이모티콘을 넣어 개발 효율성을 높입니다.](https://qiita.com/Jung0/items/0a9a7a97a2c17f92d3c5)
