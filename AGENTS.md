# AGENTS.md

## Build / Lint / Test / Run Commands

- Build : `go build -o pummit .`
- Lint: `golangci-lint run`
- Test: `go test`

## Code Style Guidelines

- Indentation: 2 spaces, max line length 80 characters
- Imports: external modules first, then internal paths; alphabetize each group
- **Import Paths (TypeScript)**: ALWAYS use @/ alias for all internal imports; NO relative paths (../or ./) allowed
- Naming: camelCase for variables/functions, PascalCase for types/classes, UPPER_SNAKE_CASE for constants
- Types: annotate all public interfaces and function signatures; avoid `any`
- Error Handling: handle errors immediately; wrap external errors with context
- Comments: see Comment Writing Rules
- Cursor rules: none
- Copilot instructions: none

### TypeScript Import Path Rules

**MANDATORY**: All TypeScript files must use @/ alias for internal imports

```typescript
// ✅ Correct
import { AppConfig } from "@/core/config.schema";
import { ILogger } from "@/services/interfaces";
import { TOKENS } from "@/core/tokens";

// ❌ Forbidden
import { AppConfig } from "../core/config.schema";
import { ILogger } from "./interfaces";
import { TOKENS } from "../core/tokens";
```

- **Never use relative paths** (../ or ./) for internal imports
- **Always use @/ alias** pointing to src/ directory
- When creating new files: use @/ alias from the start
- When modifying existing files: convert any relative paths to @/ alias
- tsconfig.json is configured with baseUrl: "./src" and paths: {"@/_": ["_"]}

---

## Language & Communication

- All responses must be in Japanese

## Workflow Efficiency

- After receiving tool output, assess quality and plan next steps before acting
- For independent tasks, invoke tools concurrently rather than sequentially
- If you encounter an error during implementation that you cannot resolve, use o3-search mcp to investigate.

## Branch Management and Task Tracking

### Branch Naming Convention

Before starting any work, create an appropriately named branch:

- `feature/working-name` - for new features
- `bug/working-name` - for bug fixes
- `fix/working-name` - for general fixes
- `refactor/working-name` - for refactoring work

### Task Checklist Management

Before beginning work, locate and update the appropriate task checklist:

1. Navigate to `docs/` directory
2. Find the appropriate version specifications
3. Open `tasks.md` in the relevant version folder
4. Check off completed items in the task checklist as work progresses

## MCP Guidelines

### Serena MCP

- When starting a new project and wanting to understand the code structure
- When you want AI to plan complex refactoring or design
- When you want to speed up bug fixes for websites or applications
- When you want to understand the overall architecture of the codebase
- When you want to quickly grasp the purpose of a file or functions
- When you want to get an overview of the codebase and its components
- When you want to read the code
- When you want to generate or edit code directly
- When you want to check logs and manage processes in a monitoring dashboard

### Context7 MCP

- When you want to reference the latest API documentation for libraries or frameworks in use
- When you want to prevent errors from outdated training data suggesting "non-existent APIs" or "deprecated methods"
- When you want to get the latest code examples immediately by instructing "use context7" in natural language in your prompt
- When dealing with rapidly evolving libraries such as Next.js, Tailwind CSS, React Query, etc.

## Comment Writing Rules

### Format

- Write all comments in English

### Content

#### Required

- Why the code was written (background, rationale)
- Business logic and rule explanations
- External dependencies and constraints
- Important notes for future developers
- Performance and security considerations
- Intent and purpose of complex algorithms

#### Prohibited

- Comments that only describe what the code does
- Details that mirror implementation
- Unnecessary or outdated explanations
- Information obvious from names

### Quality Standards

- Provide information understandable without reading the code
- Focus on “why” rather than “what”
- Include only details that remain valid when implementation changes
- Be concise and specific

### Layout

- Surround comments with blank lines
- Match indentation to code level
- Use consistent multi-line notation

## Response Structure (CoT + Answer)

1. Include reasoning in every response
2. Structure as CoT (Chain of Thoughts) followed by the final answer
3. Present CoT in a plain-text code block and close it; do not wrap the final answer
4. To preserve Markdown rendering, use `"""` for internal code snippets in CoT
5. Follow CoT with a conversational final answer
6. Do not use emojis or `**` for emphasis in the final answer
7. If corrected by the user, analyze the cause, reflect on appropriateness, and consider alternatives
