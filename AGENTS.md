## ❖ Mission

Generate Pull Requests that prioritize **clarity, testability, and minimal diff** so reviewers stay in control.

## ❖ Workflow Directives

1. **Plan First** – Emit a numbered task list (`/plan`) and wait until the human replies “approve plan”.
2. **PR Granularity**
   • One feature or bug-fix per PR.
   • Diff limit: ≤ 400 LOC (hard); warn when exceeding.
3. **Explain Changes**
   • In the PR description, include:
   - `### Why` (business/context)
   - `### How` (technical overview ≤ 10 lines)
   - `### Tests` (link to new/updated tests)
4. **Testing Mandate**
   • For every non-trivial function, add or update unit tests in `__tests__/`.
   • Reject your own PR if tests fail (`/abort`).
5. **Style & Lint**
   • Run project linters/formatters; commit resulting fixes.
6. **Interactive Flags**
   • `/explain <path#Lx-Ly>` – detailed reasoning for a code slice.
   • `/benchmark` – run provided benchmarks before/after, include table in PR.

## ❖ Tone & Comments

- Keep commit messages imperative (“Add”, “Fix”, “Refactor”).
- Use inline GitHub review comments to highlight non-obvious lines.

## ❖ Safeguards

- Never modify CI config or secrets unless plan explicitly includes it.
- If ambiguity in requirements, open an Issue (`/ask`) instead of coding.

## ❖ Review Checklist (auto-appended)

- [ ] Tests added/updated
- [ ] Lint passes
- [ ] Documentation updated
- [ ] Performance <= baseline

## ❖ MCP Tasks

| MCP Name      | Trigger Syntax                                | Typical Use Case                                                                                                             | Decision Criteria                                                                                    | Expected Behavior                                                                                                                                           |
| ------------- | --------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| serena mcp    | `@serena edit <scope> "<instruction>"`        | Large-scale refactoring, API compatibility changes, or bulk pattern replacements within the existing codebase                | • Use serena if 3+ files or 100+ LOC are affected<br>• Prefer serena if no external info needed      | 1. Ask for confirmation of scope.<br>2. Apply edits via Serena MCP.<br>3. Post diff summary + Serena dashboard link.<br>4. Mark checklist item `Serena ✔`. |
| o3 search mcp | `@o3 search "<query>"`                        | Stuck during debugging, need the latest info on library breaking changes, or want up-to-date best practices for dependencies | • Use if external web sources are required<br>• Only trigger if same query hasn't run in last 30 min | 1. Call o3‐search MCP.<br>2. Attach top-5 results table (title, gist, url).<br>3. Cache query for 30 min to avoid spam.                                     |
| context7 mcp  | `@context7 doc "<library>@<version> <topic>"` | Need official documentation or canonical code samples for API specs, configuration, or middleware usage                      | • Library/framework version is clear<br>• Use if Serena cannot resolve by itself                     | 1. Fetch latest docs via Context7 MCP.<br>2. Insert a 7-sentence summary + canonical code snippet.<br>3. Reference doc URL in PR comment.                   |

> **Selection Logic**
>
> 1. If the task involves internal code edits only, use `serena mcp`.
> 2. If external information is required:
>    　　A. For official docs or samples, use `context7 mcp`.
>    　　B. For general search needs, use `o3 search mcp`.
> 3. If in doubt, prefer: `serena` > `context7` > `o3`, and explain your choice in a comment.

## ❖ Sources & Inspiration

Guidelines derived from Anthropic Claude Code best practices and OpenCode OSS docs.
