**LANGUAGE**: Japanese (日本語)
**REPOSITORY**: HidemaruOwO/pummit

## Workflow

1. Plan First
   - Create and tell me your plan and wait for “approve plan” I says. And call `Update TODO` built in OpenCode if user says "approve".
2. Implemention
3. Tests
   - Add/update unit tests in `tests/` for every non-trivial function.
4. Quality
   - Run linters/formatters and commit resulting fixes.
   - Run build and tests locally before pushing.

## (IMPORTANT) Coding Best Practices

- Keep logic in one function unless splitting improves reuse or clarity.
- Avoid unnecessary destructuring.
- Prefer early returns; avoid `else` when not needed.
- Use short, descriptive variable names.

## Commands

- Linter: `golanci-lint run`
- Formatter: `gofmt -w .`
- Build: `go build . -o pummit`
- Test: `go test -v ./...`

## Review Checklist

- [ ] Plan approved
- [ ] One feature/bug-fix
- [ ] Tests added/updated
- [ ] Lint/format pass
- [ ] Build/test pass
- [ ] Docs updated if needed
