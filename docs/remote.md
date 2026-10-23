---
title: Remote Hosts
description: Run sessions on a Linux box, a cloud VM or another computer, and see them in the same sidebar as your local ones.
---

# Remote Hosts

Kano can show sessions that run on other machines, such as a Linux dev box, a cloud VM, a Mac mini under your desk or a WSL distribution, in the same window as your local ones. They get the same sidebar, statuses, notifications and `kano` command. Each machine runs its own kanod, so closing your laptop or quitting Kano leaves remote agents running, as it does locally.

## Quick start

A host is any machine with Kano on it. It can be another computer where you already use the Kano app, or a server that has no display.

### A computer with the Kano app

The app needs nothing installed beside it: its own kanod serves as the host. Pick how your other machines reach it:

- **Over Tailscale.** In the app on the host, turn on **Settings › Hosts › Let other machines connect**. It shows the address to add elsewhere, such as `http://studio.tail1234.ts.net:7373`. See [Tailscale](#tailscale).
- **Over SSH.** If the computer takes SSH logins (on a Mac, **System Settings › General › Sharing › Remote Login**), add it by the address you'd give `ssh`. Kano finds the app there on its own.

Sessions on the host keep running while its app is closed, as they do on your own computer.

### A server

On a Linux box, a cloud VM or any machine without the app, install the server build of Kano. It needs no root or administrator, no display and no GTK. On Linux or macOS:

```sh
curl -fsSL https://github.com/YOUR_GITHUB_USERNAME/kano/releases/latest/download/install-server.sh | sh
```

On Windows, in PowerShell:

```powershell
irm https://github.com/YOUR_GITHUB_USERNAME/kano/releases/latest/download/install-server.ps1 | iex
```

### Adding the host

On your computer, open **Settings › Hosts**, type the host's address, such as `me@devbox` for SSH or the address the host's app showed, and press **Add**. The host's projects appear in the sidebar with yours, each tagged with the host's name.

The same from the command line:

```sh
kano hosts add devbox me@devbox
kano -H devbox new -n review -- claude
```

::: tip
If Kano isn't installed on a host you add, the sidebar says so and offers **Install Kano…**, which runs the script above on the host over SSH.
:::

## How Kano reaches a host

An address says how to reach a host. Kano takes the same address in Settings, in the configuration file and after `kano -H`.

| Address | How |
| --- | --- |
| `devbox`, `me@devbox`, `me@devbox:2222`, `ssh://me@devbox` | SSH, with your `~/.ssh/config`, keys and agent |
| `wsl://Ubuntu`, `wsl://` | A WSL distribution of this Windows PC, through `wsl.exe` |
| `http://devbox:7373` | kanod's listener on a tailnet address (see [Tailscale](#tailscale)) |
| `https://devbox.tail1234.ts.net` | kanod behind Tailscale Serve |

Kano opens one connection per host and carries everything over it: the sidebar's state, every pane you look at, and the CLI's commands. The connection is compressed, and a heartbeat notices a dead link within seconds. When a link drops, as when a laptop sleeps, the host's sessions stay in the sidebar, faded, and Kano reconnects on its own, waiting longer between tries up to a minute.

### SSH

Over SSH, Kano runs `kano __bridge` on the host, which starts kanod there if it isn't running and passes the connection to it. There is nothing to set up on the host beyond Kano itself, the app or the server build, and nothing listens on the network.

Kano never prompts for anything over SSH, as there's no one to answer in the app:

- **Keys with a passphrase** need your SSH agent.
- **A host you haven't connected to before** needs its key trusted first. Run `ssh me@devbox` once in a terminal.
- **Password logins** aren't supported. Use a key.

Kano looks for `kano` on the host in `~/.local/bin` (on Windows, `%LOCALAPPDATA%\Programs\kano-server`), the desktop app and its command line tool, and your `PATH`. If yours is elsewhere, name it in the address: `me@devbox?kano=/opt/kano/kano`, or on Windows `me@pc?kano=C:/kano/kano.exe`.

### Windows hosts

Kano reaches a Windows machine over SSH, so turn on its OpenSSH Server first (**Settings › System › Optional features**, as an administrator). Kano works with the server's default shell, cmd, and with PowerShell. It finds out which one the first time it connects, which costs that connection a moment.

The install script runs kanod as a task of yours in Task Scheduler, which starts it as you log on to the machine, in your session. While no one is logged on, kanod starts as Kano connects and keeps running after the connection ends. A kanod started that way runs without your Windows password, as programs do over SSH with a key, so programs in its sessions can't reach network shares or the credentials Windows keeps for you. Log on to the machine once, or restart kanod after you do, to run them with your full account.

### WSL

On Windows, `wsl://Ubuntu` reaches the Kano installed in that WSL distribution, as SSH would. `wsl://` is the default distribution. If Kano isn't installed there, Kano offers to install it, as it does over SSH.

WSL may stop a distribution a while after its last program started from Windows ends, and kanod with it. Keep agents running by keeping the distribution up, for example with [systemd enabled](https://learn.microsoft.com/windows/wsl/systemd) and `kano daemon install-service` inside it.

### Tailscale

kanod can also listen on your tailnet, so that it needs no SSH at all. Tailscale tells kanod who connects, and by default only the Tailscale user who owns the machine gets in. On a host with the app, turn on **Settings › Hosts › Let other machines connect**, which shows the address to add elsewhere. On a server, or to set it by hand, add this to the configuration on the host:

```ini
# ~/.config/kano/config on the host
listen = tailnet://:7373
```

Then add the host as `http://devbox:7373`. The connection is plain HTTP, which WireGuard encrypts between tailnet devices.

To serve it over HTTPS with Tailscale's certificate instead, let kanod listen on a socket and point Tailscale Serve at it:

```ini
listen = tailnet-serve://
```

```sh
tailscale serve --bg --https=443 unix:$HOME/.config/kano/serve.sock
```

Then add the host as `https://devbox.tail1234.ts.net`. The socket is `serve.sock` in Kano's data folder: `~/.config/kano` on Linux, `~/Library/Application Support/kano` on macOS.

To let others in besides the owner, list them with `listen-allow`: a tailnet login, `tag:ci` for tagged devices, or `*` for anyone on the tailnet. A tagged machine has no owner, so only `listen-allow` lets anyone in to it. Tailscale Serve doesn't tell kanod who tagged devices are: they connect to a `tailnet://` listener.

```ini
listen-allow = alice@example.com
listen-allow = tag:ci
```

kanod follows `listen` and `listen-allow` as they change, and a listener turned off drops the connections it took. A `tailnet-serve://` listener added while kanod runs waits for `kano daemon restart`. kanod listens only when the configuration asks, and only on Tailscale. It refuses requests from web browsers, so a web page you visit can't reach it. Funnel's requests from the internet never get in. To keep it listening with the machine, run it as a service (below).

## Running kanod as a service

A host with the app starts kanod as the app opens, or as Kano connects over SSH. On a server, run kanod as a service of your own, so that it starts with the machine and keeps your sessions running while no one is logged in. The install scripts do this for you. On a Mac, the service is a launchd agent, which runs while you're logged in to the Mac, so keep a session logged in there; the script installs it when you are. On Windows, it's a Task Scheduler task, which starts kanod as you log on.

```sh
kano daemon install-service      # a systemd user unit, a launchd agent on macOS, a task on Windows
kano daemon status               # how it runs, and where it listens
kano daemon uninstall-service
```

On Linux this first enables lingering (`loginctl enable-linger`), which keeps your services running after you log out. Without it, a service would stop as your last SSH session ends, which a kanod started on demand outlives, so `install-service` refuses to go on without it. If your system asks an administrator for lingering, ask them to run `sudo loginctl enable-linger <you>`, or install the service anyway with `--no-linger`.

`kano daemon serve` runs kanod in the foreground and logs to its standard error, for a service manager of your own. It never exits for being idle. `--listen` takes the place of the configuration's `listen` keys.

## Working on a host

- **The sidebar.** It lists projects, yours first, then each host's in the order of **Settings › Hosts**. A host's projects carry its name after theirs, as `api` `devbox`, so the same repository on two machines sits side by side. A host gets a row of its own only when it needs you: when Kano can't reach it, wants updating, or isn't installed there. A host with no sessions doesn't show; start one there with Open Project. Right-click one of a host's projects for the host's menu: reconnect, restart kanod, update Kano, its color, its settings.
- **Colors.** Each host has a color, picked from its name: its tag in the sidebar and its name in a pane's header wear it, so a remote pane is plain to see. Pick another from the host's menu, **Color**, or with `color=` after its address (`host = devbox=me@devbox color=teal`).
- **New sessions.** <kbd>⌘</kbd> <kbd>T</kbd> and **New Session** start a session on the host of the session you're looking at. A project's **+** starts one in it, on its machine.
- **Open Project.** With hosts added, <kbd>⌘</kbd> <kbd>O</kbd> lists folders in Kano instead of your computer's dialog, starting on the machine of the session you're looking at; the switch at its top lists another's. Type to filter, or type a path such as `~/code`. <kbd>↩</kbd> goes into a folder, and <kbd>⌘</kbd> <kbd>↩</kbd> or **Open** opens it as a project.
- **Paths.** <kbd>⌘</kbd>-clicking a file path in a remote pane opens the file in the host's `$EDITOR`, in a pane below. A folder opens in a new session there. Kano never opens a file of your computer that happens to have the same path. Hovering an image's path with <kbd>⌘</kbd> previews the host's image.
- **Pasting images and files.** Pasting an image or files into a remote pane copies them to a private temporary folder on the host and pastes their paths there, so agents that read image paths, such as Claude Code, keep working. Dragging files onto a remote pane still pastes your computer's paths.
- **Agents.** Inside a remote pane, `kano` talks to that host's kanod, so the agent skill works there unchanged. Run `kano skill setup` on the host.
- **Notifications.** A host's notifications and the window's title carry the host's name, so a remote program can't pass for a local one.

### Remote Control

`kano focus` and `kano do` let a program drive the app. For a host, that's code on another machine driving your computer's window, so it's off by default. Turn it on per host in **Settings › Hosts**, or with `remote-control` after the address:

```ini
host = devbox=me@devbox remote-control
```

### Settings

Appearance, keyboard and notifications are the app's. A host's shell, `env`, scrollback and session restore come from the configuration on that host, since kanod reads them there. Run `kano config edit` on the host to change them.

## Versions and updates

The app and a host can run different versions of Kano. When the host runs an older release than the app, its first project in the sidebar shows **Update**. It downloads the app's version of Kano onto the host, checks it against the key Kano's own updates are signed with, for that version and platform, and then restarts kanod there. Restarting starts the host's sessions again where they were, with their history, and the programs running in them end, so Kano asks first.

- Kano updates a host only when you ask, never on a reconnect, and never back to an older version.
- When the host runs a newer Kano than the app, update the app instead (**Settings › General**).
- Kano updates an install of its install scripts. If you installed Kano on the host another way, update it the way you installed it.
- A host running the Kano desktop app takes updates from your other machines once you let it: on that computer, turn on **Settings › Hosts › Let them update Kano here** (`remote-updates = true`). The app there then updates to the newest release, relaunches, and restarts kanod. Until then, update it on that computer.

From the command line:

```sh
kano -H devbox version           # this kano's version, and the host's
kano -H devbox daemon status
kano -H devbox daemon update     # to this kano's version
kano -H devbox daemon restart
kano daemon update               # on the host, to the latest release
kano daemon update 2.1.0         # on the host, to a version, older ones too
```

## The command line

`-H` (or `--host`), before the command, points a command that acts on sessions at a host's kanod. It takes a host's name, an address, or `local` for this machine. `$KANO_HOST` does the same, except inside a pane, where `kano` acts on the pane's own machine. Commands of this machine alone, such as `kano config` and `kano hook`, ignore `$KANO_HOST`.

```sh
kano hosts                       # the hosts you added (--json for JSON)
kano hosts add build ci@build.lab --remote-control
kano hosts check build           # connect, and say how it went; a name or an address
kano hosts install build         # install Kano there, over SSH or WSL
kano hosts remove build
kano -H build ls
kano -H build prompt review "run the tests and fix them"
kano clients                     # who is connected to this machine's kanod
```

The app and the CLI share one list of hosts, the configuration's `host` keys:

```ini
host = devbox=me@devbox
host = build=ssh://ci@build.lab:2222 remote-control
host = ubuntu=wsl://Ubuntu
```

| Variable | What it does |
| --- | --- |
| `KANO_HOST` | The host commands act on, as `-H` |
| `KANO_SSH` | The `ssh` program Kano runs, in place of the one on your `PATH` |
| `KANO_RELEASES` | A mirror of Kano's releases, for the install scripts and `kano daemon update` |

## Security

- **A connection can do anything its account can,** including starting programs, as a local one can. Over SSH it's your own account, as SSH logs you in. Over Tailscale, only who the machine's owner lets in.
- **The socket is yours alone.** kanod's socket is created readable by its user only, and kanod refuses connections from other users' processes.
- **Updates are signed.** The server build updates itself only from archives signed with the key of Kano's releases.
- `kano clients` lists who is connected to a kanod, with how they came: `uid 501 over unix`, `me from 10.0.0.5 over ssh`, `alice@example.com over tailnet`.

## Limits

- The server build is for Linux, with glibc 2.28 or later, macOS and Windows 10 (1809) or later, on x86-64 and ARM64. A machine with the Kano app needs no server build.
- Dragging files onto a remote pane pastes your computer's paths. Paste them instead.
- One host has one account. Login servers, where several people reach their own accounts on one machine, are planned.
