---
title: Command Line
description: Drive Kano sessions and panes from a shell or a script with the kano CLI.
---

# Command Line

The `kano` command drives sessions from any shell: yours, a script's or an agent's. Install it with **Kano › Install Command Line Tool…**.

## Targets

Most commands take a **target**: a session or a pane. A target can be:

- an id, or any unique prefix of it
- a session's name
- a session, which stands for its focused pane

Inside Kano, `$KANO_SESSION` and `$KANO_PANE` name the session and pane a command runs in, and they're the default target.

## Sessions and panes

```sh
kano ls                          # every session and its panes
kano new -n api -d ~/code/api    # start a session, print its id
kano new -- npm run dev          # start one running a command
kano split --down -- npm test    # a pane beside this one
kano split --url localhost:3000  # a browser pane showing a page
kano rename api backend          # rename a session
kano focus backend               # show it in the app
kano close backend               # end a session (or, given a pane, the pane)
kano restart <pane>              # start an exited pane's shell again
```

`kano new` options:

| Flag                        | Meaning                                         |
| --------------------------- | ----------------------------------------------- |
| `-p`, `--project <name\|id>` | The project to put it in (default: by directory) |
| `-n`, `--name <name>`       | Name the session                                |
| `-d`, `--dir <dir>`         | Start it there (default: the current directory) |
| `-f`, `--focus`             | Show it in the app                              |
| `-w`, `--wait`              | Wait for its command to end, and exit with its status |
| `-u`, `--url <url>`         | Open a browser showing the page instead of a shell |

## Typing and reading

```sh
kano send api "npm run migrate" --enter   # type text, then press Enter
kano send api -e 'y\n'                    # interpret \n \r \t \e \xHH escapes
kano key api ctrl+c                       # press keys
kano read api -n 40                       # the last 40 lines on screen
kano read api -a                          # the screen with its scrollback
```

`kano key` understands `enter`, `tab`, `escape`, `up`, `down`, `left`, `right`, `backspace`, and modifiers such as `ctrl+c` or `ctrl+d`.

## Waiting

```sh
kano wait api                             # until it stops working
kano wait api --for waiting,error --timeout 10m
kano prompt review "summarize the diff"   # send + Enter, wait, print the screen
```

`kano wait` prints the pane's status when it returns. `kano prompt` sends text and Enter, waits until the pane stops working, then prints its screen. It's the quickest way to ask another agent something.

### Waiting for a command

A command such as `make test` doesn't show as *working*, so to wait for one to finish, start it with `--wait`. `kano` then exits with the command's own status:

```sh
if kano split --down --wait -- make test; then
  echo "tests passed"
fi
```

`kano wait --exit` does the same for a pane that is already running a command it started with:

```sh
pane=$(kano split --down -- go test ./...)
# … do other things …
kano wait "$pane" --exit --timeout 10m
echo "go test exited with $?"
```

The pane stays open after the command, with its output, so you can `kano read` it. Waiting for a command doesn't work with `cmd` as the shell, which can't tell when one ends.

## Browsing

`kano browser` drives a browser pane's page the way `kano send` and `kano read` drive a terminal: it clicks, types, reads the page and takes pictures of it. Its target is a browser pane, or a session, for its browser. With no target, it drives the browser of the session it runs in.

```sh
kano split --url localhost:3000            # open the page, print the pane's id
kano browser open localhost:3000/todos     # go to a page, and wait for it to load
kano browser type '#new-todo' "Buy milk" --enter
kano browser wait '.todo-list li'          # until an element shows
kano browser wait --text "1 item left"     # until the text shows
kano browser text                          # the page's text
kano browser text h1                       # an element's
kano browser click 'button[type=submit]'
kano browser screenshot -o page.png
kano browser eval 'document.title'
```

| Command                         | What it does |
| ------------------------------- | ------------ |
| `open <url>`                    | Go to a page, and wait for it to load |
| `back`, `forward`, `reload`     | Go back or forward, or load the page again, and wait for it |
| `text [selector]`               | Print the page's visible text, or the first matching element's |
| `html [selector]`               | Print the page's HTML, or the element's |
| `click <selector>`              | Click the element, as a mouse does |
| `type <selector> <text…>`       | Type text into a field, in place of what it holds; `--enter` presses Enter after it, which submits its form |
| `wait [selector] [--text text]` | Wait until the element shows, or the text does; with neither, until the page loads |
| `eval <js…>`                    | Run JavaScript in the page and print its value: a string as it is, anything else as JSON |
| `screenshot [-o file.png]`      | Save a picture of the page, and print where; `-o -` writes it to stdout |

Selectors are CSS selectors; quote them in the shell. Commands that need an element wait for it to appear, and `--timeout` (default `10s`) sets how long they, and the page loads, may take. `type` and `eval` take `-` for text from stdin.

Clicking and typing send the events a real mouse and keyboard would, so pages built with React and other frameworks see them. A screenshot needs the pane to be on screen: show it with `kano focus` first. The rest works on any browser pane, shown or not.

Like `kano do`, `kano browser` needs Remote Control, in **Settings › Agents**, and the app to be open.

## Status

```sh
kano status api --json                    # status, program, directory
kano report working "Running migrations"  # report your own status
kano notify "Need your input"             # flag this pane as needing you
```

See [Status & Notifications](./status.md) for the states.

## Changes

```sh
kano diff                                 # the pane's uncommitted changes, as git diff prints them
kano diff --turn                          # what its agent changed in its last turn
kano diff review --branch                 # since the default branch
```

See [Changes](./changes.md).

## Configuration and setup

```sh
kano config path                 # where the configuration file is
kano config show                 # every key with its value
kano config show --defaults      # every key as kano comes
kano config edit                 # open it in $EDITOR
kano config check                # report unknown keys and bad values
kano projects                    # list projects
kano skill setup                 # teach coding agents to use kano
kano daemon status|stop|restart  # manage kanod
kano daemon install-service      # run kanod with the machine (systemd, launchd)
kano daemon serve                # run kanod in the foreground, for a service manager
kano clients                     # who is connected to kanod
```

## Other machines

`-H <host>` before the command, or `$KANO_HOST`, points a command that acts on sessions at kanod on another machine: a host you added, or an address such as `me@devbox` or `wsl://Ubuntu`. See [Remote Hosts](./remote.md).

```sh
kano hosts                       # the hosts you added
kano hosts add devbox me@devbox  # add one, which the app shows too
kano hosts check devbox          # connect, and say how it went
kano hosts install devbox        # install Kano there, over SSH
kano -H devbox ls
kano -H devbox new -n review -- claude
kano -H devbox daemon update     # update Kano there to this one's version
```

Inside a pane, `kano` talks to the kanod of the machine the pane runs on, wherever the app is.

## Running app actions

`kano do` runs any of the app's actions, the same ones as its menus. It needs Remote Control, which is on by default in **Settings › Agents**.

```sh
kano actions                     # list them
kano do pane.split.right
kano do pane.browser url=localhost:3000
kano do sidebar.toggle
```

## Exit status

| Code | Meaning                              |
| ---- | ------------------------------------ |
| `0`  | OK                                   |
| `1`  | Failed                               |
| `2`  | Bad usage                            |
| `3`  | No such target, or an ambiguous one  |
| `4`  | kanod is unreachable                 |

Commands that wait for a program, `kano wait --exit` and `--wait`, exit with the program's status instead, as does `kano wait` when the pane's shell exits with an error. A wait that times out exits with `1`.

Add `--json` to `ls`, `status` and `projects` for machine-readable output.
