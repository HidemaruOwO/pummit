**LANGUAGE**: Japanese (日本語)

## Workflow

1. Plan First
   - Create and tell me your plan and wait for “approve plan” I says. And call `plan tools mcp` built in OpenCode or codex or claude code if user says "approve".
2. Branching
   - From `main`: `feature/<short>` or `fix/<short>` ..etc (Refer `Branch Naming` section).
3. Tests
   - Add/update unit tests in `tests/` for every non-trivial function.
   - If tests fail, reject your PR (`/abort`).
4. Quality
   - Run linters/formatters and commit resulting fixes.
   - Run build and tests locally before pushing.

## (IMPORTANT) Coding Best Practices

- Keep logic in one function unless splitting improves reuse or clarity.
- Avoid unnecessary destructuring.
- Prefer early returns; avoid `else` when not needed.
- Avoid `try/catch` where possible; if needed, keep scope narrow.
- Avoid `any` type.
- Use short, descriptive variable names.
- Prefer Bun APIs when applicable (e.g., `Bun.file()`).

## Commands

- Linter: `bun x eslint . --max-warnings=0`
- Formatter: `bun x prettier --write .`
- Build: `bun build`
- Test: `bun test`

## Branch Naming

- feature: new feature
- fix: bug fix
- docs: documentation only
- refactor: code change without feature change
- test: adding/updating tests
- chore: build process or auxiliary tools

## Review Checklist

- [ ] Plan approved
- [ ] New branch from `main`
- [ ] One feature/bug-fix
- [ ] Tests added/updated
- [ ] Lint/format pass
- [ ] Build/test pass
- [ ] Docs updated if needed

## MCP (Optional)

- `@context7 doc "<lib>@<ver> <topic>"` — fetch latest docs. In answers, include
  a 7-sentence summary, a minimal canonical snippet, and the official URL. See
  “System Prompt Defaults” for when to use.
  Always use context7 when to generation codes, setup or configuration steps, or
  library/API documentation. This means you should automatically use the Context7 MCP
  tools to resolve library id and get library docs without me having to explicitly ask.
- `@fetch url` - supply `url` (with optional `headers`);
  When ask have urls, always use `@fetch` MCP tool to get the content of the url and use it to answer.
