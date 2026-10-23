# Changelog

## 2.1.10

- On Windows, pasting with Ctrl+V into a text field, such as a host's
  address in **Settings › Hosts**, pastes once, not twice.

## 2.1.9

- On a Mac, Command+Backspace clears the line you're typing in the
  shell, as in Ghostty. Thanks @eliottwantz

## 2.1.8

- With the Tailscale app on a Mac, **Settings › Hosts › Let other
  machines connect** no longer says Tailscale isn't running when it is,
  and kanod listens on your tailnet as it should.

## 2.1.7

- Kano can update a host that runs the Kano desktop app, not only the
  server build: turn on **Settings › Hosts › Let them update Kano here**
  on that computer (`remote-updates = true`), and its project in your
  sidebar shows **Update** when it's behind. The app there updates,
  relaunches and restarts kanod. It's off unless you turn it on, and the
  host needs 2.1.7 itself, so update it there once more first.

## 2.1.6

- A session open in Kano on two machines no longer takes the size of a
  machine you've walked away from: Kano sizes it for the one you're at,
  and gives it back to the other as soon as you return there or this one
  quits or sleeps. Hosts get it with their next update.
- Panes and sessions on a Windows host are named by their folder
  (`mygo`), not its whole path (`C:\Users\you\dev\mygo`).

## 2.1.5

- A session open in Kano on two machines takes the size of the one you're
  using, the one you last typed, scrolled or switched to the session in,
  rather than of the last to resize it. Resizing a split on one machine
  no longer cuts off the bottom of the pane on the other. Hosts get it
  with their next update.
- Dragging the divider between panes, or the sidebar's edge, leaves the
  keyboard in the pane you were typing in.
- After an update, Kano offers to restart kanod only when the update
  changed it.

## 2.1.4

- The folder picker and the command palette scroll with the mouse wheel
  and trackpad, rather than snapping back to the highlighted row. Thanks
  @ak4duy.

## 2.1.3

- On Windows, turning on **Let other machines connect** in Settings ›
  Hosts no longer flashes a console window open, nor does Kano on
  another machine connecting to it.
- Settings › Hosts leaves the commands for setting up a host to the guide
  it links to.

## 2.1.2

- Settings › Hosts can share this computer with Kano on your other
  machines: turn on **Let other machines connect**, and copy the address
  it shows to add it there. kanod then listens on your tailnet, and only
  your Tailscale account may connect. kanod now follows `listen` and
  `listen-allow` in the config file as they change, without a restart,
  and a listener turned off drops the connections it took.
- Hosts can be edited in Settings › Hosts: a host renamed keeps its
  connection and sessions, and one given another address connects again.
  The page links to the guide to setting up hosts.
- On Windows, a restored pane's shell starts under its history, rather
  than drawing what you type at the top of the pane.

## 2.1.1

- On Windows, the mouse wheel scrolls full-screen programs that take the
  mouse as they start, such as pi, rather than stepping through their
  prompt history. Kano now runs panes in the newer console Windows
  Terminal uses, which it ships, so programs get the mouse, the wheel and
  what they draw as they would there. Windows hosts get it with their
  next update.

## 2.1.0

- Browser panes show web pages beside your terminals: the dev server an
  agent works on, the docs it reads, the pull request it opened. One
  splits, zooms, moves and closes like any pane, and comes back on the
  same page when Kano restarts. Open one with <kbd>⌘</kbd> <kbd>⇧</kbd>
  <kbd>B</kbd>, by typing an address in the command palette, or with
  `kano split --url localhost:3000`; <kbd>⌘</kbd> <kbd>L</kbd> goes to
  its address bar, which shows the site's icon.
- `kano browser` drives a browser pane's page the way `kano send` and
  `kano read` drive a terminal: it opens pages, clicks, types, waits for
  elements or text, reads the page's text and HTML, runs JavaScript and
  takes screenshots. Like `kano do`, it needs Remote Control.
- Browser panes block ads and trackers, with EasyList and EasyPrivacy to
  begin with. Settings › Browser turns it off or picks other filter
  lists, or `browser-block-ads` and `browser-filter-list` in the config
  file, which takes the URL of any list in Adblock Plus syntax too.
- With `focus-follows-mouse`, the pane the pointer moves into gets the
  focus.
- Dim unfocused terminals, in Settings › Terminal, or
  `dim-unfocused-terminals = false`, keeps the terminals you're not in
  from fading.
- The command palette finds commands and Settings pages by keywords, and
  by their English names in any language.
- Kano shows sessions that run on other machines, such as a Linux dev box,
  a cloud VM or a WSL distribution, in the same sidebar, tagged with each
  host's name, with their statuses and notifications. Add a host in
  Settings › Hosts by the address you'd give ssh, `me@devbox`: Kano
  reaches it over
  SSH with your keys and agent, or through kanod's listener on your
  tailnet (`listen = tailnet://:7373`), which Tailscale tells who
  connects to. One compressed connection carries all of a host's panes,
  and a host whose link drops shows its sessions faded until Kano
  reconnects on its own.
- On a host, New Session, Open Project (which lists the host's folders),
  Command-clicking a path (which opens it in the host's `$EDITOR`) and
  pasting an image or files (which copies them to the host) all act on
  the host. Its programs drive the app with `kano focus` and `kano do`
  only once you turn on its Remote Control.
- The sidebar stays flat with hosts: their projects follow yours, each
  tagged with its host's name, so the same repository on two machines
  sits side by side. A host has a row of its own only when it needs you,
  as when it can't be reached; one with no sessions doesn't show. Its
  menu is in its projects'. Each host wears a color, picked from its name or set
  with `color=`, on its tag and its name in a pane's header. With hosts,
  Open Project lists any machine's folders in Kano, with a switch
  between them.
- The server build of Kano runs kanod and the CLI on Linux, macOS and
  Windows with no app, display or GTK: `curl -fsSL
  https://github.com/YOUR_GITHUB_USERNAME/kano/releases/latest/download/install-server.sh
  | sh`, or on Windows `irm .../install-server.ps1 | iex`. The app offers
  to install it on a host that has no Kano, and to update it to the app's
  version, signed as the app's updates are. Windows hosts are reached
  over OpenSSH, whose shell may be cmd or PowerShell.
- `kano -H <host>` and `$KANO_HOST` run commands on a host's sessions; `kano
  hosts` adds, checks and installs hosts; `kano clients` lists who is
  connected to kanod; `kano daemon install-service` runs kanod as a
  systemd user service, a launchd agent or a Task Scheduler task, and
  `kano daemon serve` in the foreground.
- kanod's socket is its user's from the moment it is made, and kanod
  refuses connections of other users' processes. The login shell is
  looked up through NSS, so LDAP and SSSD accounts get theirs.

## 2.0.11

- The changes panel no longer jumps to the bottom of the list the first
  time a long file's diff scrolls into view, and the header of the file at
  the top of the list stays there as a file is folded, rather than
  vanishing for a moment. Only the diffs in view are built, so a long list
  of changes scrolls as quickly as a short one.
- On Windows 11 24H2 and later, clearing a pane clears the console's screen
  too, so PowerShell's inline predictions no longer draw on the line the
  prompt was on before.

## 2.0.10

- Each pane's title shows the branch it works on and the lines its changes
  add and remove since the last commit, `fix-login +40 −2`, and a merge or
  rebase in progress; click it to open the changes panel. A session on a
  branch other than the repository's default shows that branch in the
  sidebar, and its tooltip how far it is ahead of or behind its upstream
  and how many files changed.
- The changes panel (Option-Command-G, Ctrl+Shift+G on Windows and Linux)
  lists the files a session changed since its agent's last turn, since the
  last commit, or since the branch left the default one, each with its diff
  below it, which follows the file as it changes. Fold a file, resize the
  panel, switch between unified and side-by-side diffs, and in Last Turn
  revert a change or mark it kept.
- `kano diff` prints a pane's changes as `git diff` does, untracked files
  included; `--turn` and `--branch` choose what to compare against. `kano
  status --json` and `kano ls --json` include each pane's branch and
  changes.
- Sessions in a worktree outside its repository join the repository's
  project instead of starting one of their own.
- `git-status = false` stops Kano running git in your repositories, and
  `diff-layout` keeps the changes panel's diffs unified or side by side.

## 2.0.9

- Only agents that show their work with a progress bar, as Kiro and
  Qwen, turn their pane working by it: Neovim's language-server progress
  no longer marks its pane as working.

## 2.0.8

- Agents no longer rename the session they start in: an unnamed session
  keeps showing its agent's own title, as "Duplicate spinner in pane
  header", not a name in words joined by dashes.
- Files dropped on a pane paste their paths quoted for the program's
  shell, newlines and PowerShell's included.
- Find in a pane highlights a match that wraps in from above the view.
- A pane's scrollback stays whole as it reaches its limit, and a restored
  pane keeps the blank rows at the bottom of its screen.

## 2.0.7

- A split pane's title leads with its program's logo: Claude Code's,
  Codex's and other agents', Node.js's for a script node runs, as vite,
  and those of bun, deno, python, vim, git and docker; a shell's is a
  bolder terminal icon.
- A split pane shows its status once, by its title, not again in the
  chip beside it.
- An unnamed session takes its name from the program in the pane you're
  in, else in the one that most needs you; `session-name` chooses
  between these and the first pane's folder.
- Clear Scrollback in a pane's context menu clears as Cmd+K does, the
  screen too, and works on a full-screen program's.

## 2.0.6

- Images programs show appear in their panes, as grok's preview of a
  pasted image, through the Kitty graphics protocol, as in Ghostty.
- Command-hovering an image's path in a pane previews the image.
- Files dropped on a pane paste their paths, escaped for the shell, as in
  Ghostty: an agent takes a dropped image as a file it attaches.
- The mouse lands where it points in programs that take it in pixels:
  Claude Code's "Jump to bottom" lights up and takes a click, and dragging
  selects its text again.
- Programs asking the terminal's colors get the theme's: Copilot CLI's
  input box is light on a light theme. Colors a program sets for the
  terminal show, and programs that follow the appearance, as Claude Code
  and Copilot CLI, change with it when it goes from light to dark.
- Box drawing and blocks keep their colors under `minimum-contrast`, and
  only rows a program fills extend into a pane's padding, as in Ghostty:
  Copilot CLI's input box and tabs draw as they should.
- A name an agent gave its session with `kano rename` goes when a new
  conversation starts in the session's first pane.
- After updating, Kano asks to restart kanod, whose sessions restore as
  after a restart of the Mac.

## 2.0.5

- Find in a pane's scrollback with Cmd+F (Ctrl+F on Windows).
- Text stays readable on a light background: colors meant for a dark one,
  as Claude Code's near-white grays, get enough contrast with it. The new
  `minimum-contrast` setting sets how much, as in Ghostty.
- Command-hovering a link in a terminal underlines it.
- A pane shows where its agent works: the worktree Claude Code went into,
  rather than the folder its shell started in. Links, Copy Directory and
  Reveal follow it.
- Codex sessions are told to the pane that runs them, for resuming after
  kanod restarts, rather than to the pane that first started Codex's
  background service.

## 2.0.4

- Terminals show the cursor chosen in Settings, its shape and whether it
  blinks, rather than always a blinking block.

## 2.0.3

- An agent sitting at its prompt no longer shows as needing you, as Claude
  Code did a while after answering. A bell or a program's desktop
  notification leaves its status alone; programs' notifications show as
  they are.

## 2.0.2

- Installing the command line tool on macOS asks for an administrator's
  password instead of failing.

## 2.0.1

- Kano speaks Simplified Chinese; pick the language in Settings.
- Rename a project from its menu in the sidebar.
- `kano skill setup` hooks Claude Code, Codex, Pi and OpenCode into Kano as
  they start: an agent learns its session and pane, names an unnamed session
  after its task, and resumes when kanod restarts.
- `kano new|split --wait` and `kano wait --exit` exit with the command's
  status.
- When kanod runs another version of Kano, a Restart Service chip offers to
  restart it. Settings shows updates on its General page.
- A home directory pane shows `~` once in its header.

## 2.0.0

Kano, rewritten from the ground up: sessions grouped by project in a sidebar
that shows which are working, which need you and which are done. Sessions run
in kanod, a daemon of their own, so quitting Kano leaves them running.

- Runs on macOS, Windows and Linux.
- The `kano` CLI drives sessions from a shell or an agent.
- Updates itself from Kano's releases.
