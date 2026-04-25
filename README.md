<div align="center">

# pummit 🚛

Small Git commit CLI with emoji and AI-automation.

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://go.dev/)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Release](https://img.shields.io/github/v/release/HidemaruOwO/pummit)](https://github.com/HidemaruOwO/pummit/releases)

</div>

## 🚀 Features

- A series of methods of commits
- Alias management and TOML-based configuration
- Supporting MCP for agentic coding
- Cross-platform support for Windows, macOS, and Linux

## 🛠 Installation

```bash
go install github.com/HidemaruOwO/pummit@latest
```

### Build from Source

```bash
git clone https://github.com/HidemaruOwO/pummit.git
cd pummit
go build -o pummit .
```

## 🎯 Usage

### Commit

```bash
# Belief form
pummit sparkles "Add login page"

# Explicit command form
pummit commit --emoji sparkles "Add login page"

# Auto emoji from branchMapping
pummit commit --auto-emoji "Fix auth flow"
```

### Config

```bash
pummit config validate
pummit config list
pummit config get base.emoji
pummit config set base.emoji false
pummit config reset --force
```

### Alias

```bash
pummit alias add featx sparkles --emoji "✨"
pummit alias list
pummit alias delete featx --force
pummit alias reset --force
```

### Doctor / Migrate / MCP

```bash
pummit doctor
pummit migrate
pummit migrate status
pummit migrate rollback --confirm
pummit mcp
```

## 📦 Built-in Aliases

| Alias | Emoji | Meaning |
| --- | --- | --- |
| `s`, `feat`, `feature` | ✨ | feature |
| `b`, `fix`, `error` | 🐛 | bug fix |
| `d`, `doc`, `docs` | 📚 | docs |
| `a`, `ui`, `design` | 🎨 | ui |
| `c`, `wip` | 🚧 | work in progress |
| `r`, `pr`, `merge` | ♻️ | refactor |
| `l`, `test` | 🚨 | test |
| `w`, `tool`, `config` | 🔧 | config/tooling |
| `h`, `perf`, `performance` | 🚀 | performance |

## 🧩 Config

Current config format is TOML.

Default location:
- Windows: `%AppData%/pummit/config.toml`
- macOS/Linux: config dir under `~/.config/pummit`

Legacy `config.json` is migrated to TOML by `pummit migrate` and on compatible legacy paths.


## 📜 License

This project is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).


## 🤝 Reference

This repository was created using the [MicroRepository](https://github.com/HidemaruOwO/MicroRepository) template.

- [HidemaruOwO/MicroRepository](https://github.com/HidemaruOwO/MicroRepository)

---

<div align="center">

**Made with ❤️ by [HidemaruOwO](https://github.com/HidemaruOwO)**

If the projects helps improve your quality of life, please ⭐ this repository!

</div>
