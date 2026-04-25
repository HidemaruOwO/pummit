package pummit

// #Config is the normative configuration schema for the rewrite.
#Config: close({
    meta:        #MetaConfig
    base:        #BaseConfig
    interactive: #InteractiveConfig
    locale:      #LocaleConfig
    templates:   #TemplatesConfig
    scope:       #ScopeConfig
    alias:       #AliasConfig
    branchMapping: #BranchConfig
})

#MetaConfig: close({
    version: string & !=""
})

#BaseConfig: close({
    emoji:       bool
    filesLength: int & >=0
})

#InteractiveConfig: close({
    enabled:     bool
    defaultMode: "emoji" | "template" | "scope" | "branch"
    showPreview: bool
    fuzzySearch: bool
})

#LocaleConfig: close({
    language:   "ja" | "en"
    autoDetect: bool
})

#TemplatesConfig: close({
    enabled:         bool
    defaultTemplate: string & !=""
    definitions: [string]: #TemplateDefinition
})

#TemplateDefinition: close({
    format:             string & !=""
    description:        string & !=""
    defaultEmoji?:      string & !=""
    scope?:             bool
    requireDescription?: bool
})

#ScopeConfig: close({
    enabled:      bool
    autoDetect:   bool
    suggestions:  [...string]
    fromHistory:  bool
    historyLimit: int & >=0
})

#AliasConfig: close({
    enabled: bool
    entries: [...#AliasEntry]
})

#AliasEntry: close({
    shortcuts: [...string]
    name:      string & !=""
    emoji:     string & !=""
})

#BranchConfig: close({
    enabled:  bool
    fallback: string & !=""
    rules:    [...#BranchRule]
})

#BranchRule: close({
    pattern:     string & !=""
    emoji:       string & !=""
    description: string & !=""
})

// #DefaultConfig is the canonical default configuration.
#DefaultConfig: #Config & {
    meta: {
        version: "3.0"
    }

    base: {
        emoji:       true
        filesLength: 50
    }

    interactive: {
        enabled:     true
        defaultMode: "emoji"
        showPreview: true
        fuzzySearch: true
    }

    locale: {
        language:   "ja"
        autoDetect: true
    }

    templates: {
        enabled:         true
        defaultTemplate: "default"
        definitions: {
            default: {
                format:      "{emoji} {message} ({files})"
                description: "Default template"
            }
            feat: {
                format:             "{emoji} {scope}: {message}\n\n{description}\n\n({files})"
                description:        "Feature addition template"
                defaultEmoji:       "sparkles"
                scope:              true
                requireDescription: true
            }
        }
    }

    scope: {
        enabled:      true
        autoDetect:   true
        suggestions:  ["api", "ui", "core", "auth", "db"]
        fromHistory:  true
        historyLimit: 50
    }

    alias: {
        enabled: true
        entries: [
            {shortcuts: ["s", "feat", "feature"], name: "sparkles", emoji: "✨"},
            {shortcuts: ["c", "wip"], name: "construction", emoji: "🚧"},
            {shortcuts: ["t", "new", "init"], name: "tada", emoji: "🎉"},
            {shortcuts: ["r", "pr", "pull", "merge"], name: "recycle", emoji: "♻️"},
            {shortcuts: ["wb", "rm", "remove", "del", "delete"], name: "wastebasket", emoji: "🗑️"},
            {shortcuts: ["b", "fix", "error"], name: "bug", emoji: "🐛"},
            {shortcuts: ["e", "lint", "format", "refactor"], name: "eyes", emoji: "👀"},
            {shortcuts: ["d", "doc", "docs", "document", "documents"], name: "books", emoji: "📚"},
            {shortcuts: ["a", "ui", "design", "icon", "icons"], name: "art", emoji: "🎨"},
            {shortcuts: ["h", "tune", "tuning", "perf", "perform", "performance"], name: "rocket", emoji: "🚀"},
            {shortcuts: ["w", "change", "tool", "tools", "lib", "library"], name: "wrench", emoji: "🔧"},
            {shortcuts: ["l", "test", "testing"], name: "rotating_light", emoji: "🚨"},
            {shortcuts: ["p", "pack", "mod", "module"], name: "package", emoji: "📦️"},
        ]
    }

    branchMapping: {
        enabled:  true
        fallback: "construction"
        rules: [
            {pattern: "^feature/.*", emoji: "sparkles", description: "Feature branch"},
            {pattern: "^(fix|bugfix|hotfix)/.*", emoji: "bug", description: "Bug fix branch"},
            {pattern: "^docs/.*", emoji: "books", description: "Documentation branch"},
            {pattern: "^refactor/.*", emoji: "eyes", description: "Refactoring branch"},
            {pattern: "^test/.*", emoji: "rotating_light", description: "Test branch"},
        ]
    }
}
