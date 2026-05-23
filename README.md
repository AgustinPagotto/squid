```
               .:=:::::::::::@%.
                .@-=#=+=+----::::::.
                 ::--=======%*-@::-::.
                .%@:::-==-*===++---:::@      ..
                 @@@:::-=--==@=*+++-:::      @@:.
                  @*=:::%::=-======++---=.    .::.
                   *@+:::::%--====@==++++-.    @::::::%
                    @%::::-::-#--=====+++-:..       #::.
                     .@#::::::::::-========-:@++..   ::.
                        .#.:::::::::::-::-=+-%++==+ .-:.=-::%:+:.
                           .:@:-@::::::::-*+*++*==-=-@*+::@...:::
                               .@:@=::::@+++++===@-@*%#:@.     .@
                                    .@*:@%@:+===---@=@@.        @
                                        %@@@.#*@---+==---------
                                 .#:.+   ..-.+**=-@-#%*:%:=.+:@#
                               -::::+:::.   @-++@-+--       @:.:
                                ..  .@@::::---==-=:::.       @=-
                                       @*::-.+-:*@.::#.       :.@
                                            +@#:.   =::::.     .@#
                                            @@:       @:::::
                                            .:.           :::
                                             .#:           .@
```

<div align="center">

### squid &nbsp;·&nbsp; project-aware developer context CLI

[![Go](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)

</div>

---

**squid** is a project-aware developer context CLI. It lives inside your project and keeps your shell aliases, notes and todos close to the code they belong to. 

## Why squid

Every project has its own mental model: the commands you run, the things you need to remember, the tasks still pending. squid stores all of that inside a `.squid` folder at the root of your project, so the context travels with the repo.

## Installation

```sh
go install github.com/AgustinPagotto/squid@latest
```

Requires Go 1.21+.

## Getting started

Navigate to your project root and run:

```sh
squid init
```

This creates a `.squid/` folder and walks you through setting up your context. You can choose a **predefined** context (Go, Node.js, Python, React, React Native, Next.js) or  custom, where you define your own aliases from scratch. You can checkout the predefined and modify them if that's your thing. If you do custom a terminal editor will open where you can write your aliases.

### Shell hook

For aliases to actually apply to your shell session, squid needs a small hook in your shell config. During `init` you will be prompted to set it up, or you can run it separately at any time:

```sh
squid init shell
```

This adds a wrapper function to your `.zshrc` or `.bashrc` that intercepts `squid context activate` and pipes it through `eval`.

## Commands

### `squid init`

Initialize squid in the current directory.

```sh
squid init           # interactive setup
squid init shell     # install the shell hook only
squid init -h        # show help
```

---

### `squid context`

Manage your project's shell aliases — short names for the commands you run most.

```sh
squid context add              # open your $EDITOR to add or edit aliases
squid context list             # list all aliases with their IDs
squid context activate         # load aliases into the current shell session
squid context del <id>         # delete an alias by ID
squid context -h               # show help
```

**Alias format**

Each alias is one line in your editor:

```
start="npm run dev"
test="go test ./..."
deploy='./scripts/deploy.sh production'
```

Rules:
- Name: letters, digits, `_`, `-` (max 50 chars)
- Value: non-empty string in double or single quotes (max 250 chars)
- Lines starting with `#` are treated as comments and ignored

**Activating your context**

This is the most important command and why squid was created. This command will execute the defined aliases, making the current session of the terminal custom to the project you are working on, this way you won't have to remember commands or look them up.

```sh
squid context activate
```

This must be wrapped in `eval` to actually apply the aliases to your shell. The shell hook installed by `squid init shell` does this automatically so you can just type the command directly.

---

### `squid notes`

Store free-form notes scoped to the project — architecture decisions, gotchas, links, anything worth remembering.

```sh
squid notes add              # open your $EDITOR to write a new note
squid notes list             # list all notes (title + truncated preview)
squid notes show <id>        # show the full content of a note
squid notes edit <id>        # open a note in your $EDITOR to edit it
squid notes del <id>         # delete a note
squid notes search <query>   # fuzzy-search notes by content
squid notes -h               # show help
```

Notes have a title (first non-comment line) and a body (everything after).

---

### `squid todos`

A lightweight project-local todo list.

```sh
squid todos add "<text>"     # add a new todo item
squid todos list             # list all todos, grouped by status
squid todos show <id>        # show details of a todo
squid todos toggle           # interactively toggle a todo done/pending
squid todos del <id>         # delete a todo
squid todos clear            # remove all completed todos
squid todos -h               # show help
```

Todo text is capped at 280 characters.

---

## Storage

Everything lives inside `.squid/` at your project root:

```
.squid/
  context.json   # shell aliases
  notes.json     # project notes
  todos.json     # todo items
```

You can commit `.squid/` to share context with your team, or add it to `.gitignore` to keep it personal.

## Global flags

```sh
squid -h    # show top-level help
```

