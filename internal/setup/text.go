package setup

const unknownSubcommandTemplate = `
  squid init: unknown subcommand %q

  Valid subcommands:
    squid init          interactive context setup
    squid init shell    install the shell hook

  Run 'squid init -h' for more information.

`

const commandHelp = `Initialize squid in the current directory

Usage:
  squid init           interactive setup — choose predefined or custom context
  squid init shell     install the shell hook into your shell config

Predefined:
  Pick a language or framework (Go, Node.js, Python, React, React Native,
  Next.js) and squid will create a ready-to-use context with common aliases.

Custom:
  Opens your $EDITOR so you can define your own aliases in the format:
    <name>="<command>"
  Lines starting with '#' are ignored.

Shell hook:
  Adds a small function to your shell config (.zshrc, .bashrc, or
  config.fish) so that 'squid context activate' actually applies aliases
  to your current shell session.

Flags:
  -h    show this help message
`

const predefinedInitTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              squid  ·  predefined              │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   Select a predefined context:                 │
  │                                                │
  │   1  ·  Go                                     │
  │   2  ·  Node.js                                │
  │   3  ·  Python                                 │
  │   4  ·  React                                  │
  │   5  ·  React Native                           │
  │   6  ·  Next.js                                │
  │                                                │
  │   7  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`

const initialOptionTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │                 squid  ·  init                 │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   How would you like to set up your context?   │
  │                                                │
  │   1  ·  Predefined  ─  Node, Go, React, ...    │
  │   2  ·  Custom      ─  build your own          │
  │                                                │
  │   3  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`

const successTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              ✓  squid  ·  done  ✓              │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   squid has been initialized successfully.     │
  │                                                │
  │   To activate your context, run:               │
  │                                                │
  │   $ squid context activate                     │
  │                                                │
  └────────────────────────────────────────────────┘

`

const successWithShellTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              ✓  squid  ·  done  ✓              │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   squid has been initialized successfully.     │
  │                                                │
  │   Reload your shell config, then activate:     │
  │                                                │
  │   $ %-43s│
  │   $ squid context activate                     │
  │                                                │
  └────────────────────────────────────────────────┘

`

const overrideWarningTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              squid  ·  warning                 │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   A .squid folder already exists here.         │
  │   This will override your current context.     │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   1  ·  Yes, override it                       │
  │   2  ·  No, keep it                            │
  │                                                │
  └────────────────────────────────────────────────┘

`

const zshEvalConfig = `
squid() {
  if [[ "$1 $2" == "context activate" ]]; then
    eval "$(command squid "$@")"
  else
    command squid "$@"
  fi
}
`

const bashEvalConfig = `
squid() {
  if [ "$1 $2" = "context activate" ]; then
    eval "$(command squid "$@")"
  else
    command squid "$@"
  fi
}
`

const fishEvalConfig = `
function squid
  if test "$argv[1] $argv[2]" = "context activate"
    eval (command squid $argv)
  else
    command squid $argv
  end
end
`

const selectShellTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              squid  ·  shell selection         │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   Select a shell:                              │
  │                                                │
  │   1  ·  bash                                   │
  │   2  ·  zsh                                    │
  │   3  ·  fish                                   │
  │                                                │
  │   4  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`

const shellNextDialogTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              squid  ·  shell setup             │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   For your aliases to work, squid needs to     │
  │   add a small hook to your shell config        │
  │   (e.g. .zshrc, .bashrc).                      │
  │                                                │
  │   This lets squid activate your context        │
  │   automatically when you run:                  │
  │                                                │
  │   $ squid context activate                     │
  │                                                │
  │   Without it, aliases will be printed but      │
  │   not applied to your shell session.           │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   Set up shell hook now?                       │
  │                                                │
  │   1  ·  Yes                                    │
  │   2  ·  No, I'll do it later                   │
  │                                                │
  └────────────────────────────────────────────────┘

`

const addAliasesTemplate = `
# Add your custom aliases below, one per line.
# Each alias will be loaded into your shell when you run:
#
#   squid context activate
#
# Format:
#   <name>="<command>"
#
# Examples:
#   start="npm run dev"
#   test="go test ./..."
#   deploy="./scripts/deploy.sh production"
#   lint="golangci-lint run ./..."
#
# Lines starting with '#' are ignored.
# Save and close the editor when you are done.
# An empty file aborts the setup.
#
`
