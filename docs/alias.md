## About the alias function 📎

for example`wastebasket`を入力するのは少し大変ですが、エイリアス機能を使うと`wb`You can enter it easily.

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

In this case`s`Just enter the alias "Emoji prefix" in the commit message.`sparkles`will be able to be substituted.

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

If the alias`s=sparkles`and`t=tada`If it is registered, the following will be output.

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

Even if there are a lot of confusing aliases like this:

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
