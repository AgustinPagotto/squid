package setup

const commandHelp = `Initialize squid in the current directory

Usage:
  squid init

Creates a .squid/ directory in the current directory. Run this once
per project. All squid data (notes, todos, context) will be stored there.
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
  │   7  ·  ← Back                                 │
  │   8  ·  Exit                                   │
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
  │            ✓  squid  ·  done  ✓                │
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

const zshEvalConfig = `
squid() {
  if [[ "$1 $2" == "context activate" ]]; then
    eval "$(command squid "$@")"
  else
    command squid "$@"
  fi
}
`

const selectShellTemplate = `
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              squid  ·  predefined              │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   Select a predefined context:                 │
  │                                                │
  │   1  ·  bash				   │
  │   2  ·  zsh                                    │
  │   3  ·  fish				   │
  │                                                │
  │   4  ·  ← Back                                 │
  │   5  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`
