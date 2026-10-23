---
title: Quick Start
description: Install Kano, open a project and start your first agent session.
---

# Quick Start

Install Kano, open a project, and run your first agent in a session that keeps going after you quit.

## Install

Download the latest Kano from the [releases page](https://github.com/YOUR_GITHUB_USERNAME/kano/releases). On macOS, drag **Kano.app** into your Applications folder; on Windows, run the **Kano Setup** installer, which needs no administrator rights.

On Windows, sessions run in PowerShell 7 when it is installed, else Windows PowerShell; choose another shell, such as `cmd` or Git Bash's `bash`, in **Settings › Shell**. Kano needs Windows 10 1809 or later, for its pseudo consoles.

### Install the command line tool

Choose **Kano › Install Command Line Tool…** to put `kano` on your `PATH`: a link in `/usr/local/bin` on macOS, in `~/.local/bin` on Linux, and on Windows Kano's `bin` folder added to your user `PATH`, for terminals you open from then on. Run it in any terminal to check:

```sh
kano help
```

Typing `kano` alone in a terminal opens the app.

## Open a project

Press <kbd>⌘</kbd> <kbd>O</kbd> and choose a folder. Kano adds it to the sidebar as a project and opens a session there.

You can also start from the command line. A new session goes into the project that matches its directory:

```sh
cd ~/code/my-app
kano new --focus
```

## Start an agent

Run your agent in the session the way you would in any terminal:

```sh
claude
```

Or give it a session of its own with a name you'll recognize in the sidebar:

```sh
kano new -n fix-login -- claude
```

While the agent works, its session shows as **working**. When it stops to ask you something, it turns to **waiting** and Kano notifies you. Press <kbd>⌘</kbd> <kbd>J</kbd> to jump to it.

## Quit without losing anything

Press <kbd>⌘</kbd> <kbd>Q</kbd>. Your sessions keep running in kanod. Open Kano again and they're all there, with their scrollback.

## Next steps

- [Panes](./panes.md): split a session to run a dev server beside your agent
- [Coding Agents](./agents.md): run `kano skill setup` so agents learn to use Kano
