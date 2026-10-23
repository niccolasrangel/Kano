package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"kano/internal/proto"
)

// skillGuide teaches a coding agent kano: the skill `kano skill setup`
// installs, from this binary, so that it follows kano's updates.
const skillGuide = `---
name: kano
description: Drive terminal sessions in Kano — start agents and commands in their own sessions or panes, type into them, wait for them, read their screens, and report your own status. Use when running inside Kano ($KANO_PANE is set) and the task needs another terminal, a long-running process, a parallel agent, or the user's attention.
---

# Kano

You run inside Kano, a terminal whose sessions keep running in the
background. $KANO_SESSION is your session and $KANO_PANE your pane; the
` + "`kano`" + ` command drives them. A target is a session or a pane id (a
prefix works), or a session's name; a session stands for its focused pane.

## Run something beside you

    kano split --down -- npm run dev          # a pane under yours; prints its id
    kano new -n review -- claude              # a session of its own
    kano new -n tests -d ~/code/app -- go test ./...
    kano split --url localhost:3000           # a browser pane, to show the user a page

## Drive it

    kano send <target> "some text" --enter    # type, then press Enter
    kano key <target> ctrl+c                  # enter, escape, up, tab, ctrl+d…
    kano prompt <target> "run the tests"      # send + Enter, wait, print screen
    kano wait <target> [--for idle,waiting,error,exited] [--timeout 10m]
    kano wait <target> --exit                 # until its command ends; exits with its status
    kano read <target> [-n 40] [-a]           # the screen, or with scrollback
    kano status <target> --json               # status, program, directory
    kano ls                                   # every session and pane

A pane's status is working, waiting (it needs the user), idle, error or
exited. ` + "`kano wait`" + ` with no --for waits until it stops working. A
command such as a test run doesn't show as working: wait for it with --exit,
or start it with --wait, which exits with the command's status:

    pane=$(kano split --down -- go test ./...)
    kano wait "$pane" --exit || kano read "$pane" -n 80
    kano split --down --wait -- make build    # the same, in one step

## Use a web page

A browser pane's page is yours to drive, as a user would, to check what
you built: open it with kano split --url, then

    kano browser open localhost:3000/todos     # go to a page, wait for it to load
    kano browser type '#new-todo' "Buy milk"   # fill a field (--enter submits)
    kano browser click 'button[type=submit]'   # click an element (CSS selector)
    kano browser wait --text "Welcome"         # or wait <selector>; --timeout 30s
    kano browser text [selector]               # the page's text, or an element's
    kano browser html [selector]
    kano browser eval 'document.title'         # run JavaScript, print its value
    kano browser screenshot -o /tmp/page.png   # the pane must be on screen
    kano browser back|forward|reload

With no target, kano browser drives your session's browser pane; name
one (kano browser <target> text) when there are several. It needs the app
open and Remote Control on, in Settings › Agents.

## Tell the user

    kano report working "Migrating the database" --progress 40
    kano report blocked "Approve the production deploy?" --kind permission
    kano report done "Deployed to staging"
    kano notify "Need your API key to continue"

Reports are OSC 7501 sequences: Kano shows them in its sidebar, and
notifies the user when you are blocked or done while they look elsewhere.
Report blocked only when you truly cannot go on without them.

## Etiquette

- Close what you started when it is done: ` + "`kano close <target>`" + `.
- Never type into a pane you did not start unless the user asked.
- Exit status 3 means no such target, 4 that Kano is not running.
`

func skillCmd(args []string) error {
	f, err := parse(args, map[string]string{"agent": "agent"}, map[string]string{"dry-run": "dry", "n": "dry"})
	if err != nil {
		return err
	}
	sub := "status"
	if len(f.pos) > 0 {
		sub = f.pos[0]
	}
	agent := strings.ToLower(f.vals["agent"])
	dirs, hooks := skillDirs(agent), hooksFor(agent)
	switch sub {
	case "show":
		fmt.Print(skillGuide)
		return nil
	case "status":
		for _, name := range sortedKeys(dirs) {
			fmt.Printf("%-8s skill  %-9s %s\n", name, skillState(dirs[name]), proto.ShortDir(dirs[name]))
		}
		for _, h := range hooks {
			fmt.Printf("%-8s hook   %-9s %s\n", h.name(), h.state(), proto.ShortDir(h.path()))
		}
		return nil
	case "setup", "remove":
		for _, name := range sortedKeys(dirs) {
			dir := dirs[name]
			switch {
			case f.set["dry"]:
				fmt.Printf("would %s %s (%s)\n", sub, proto.ShortDir(dir), name)
			case sub == "remove":
				if err := os.RemoveAll(dir); err != nil {
					return err
				}
				fmt.Printf("removed %s\n", proto.ShortDir(dir))
			default:
				if err := writeSkill(dir); err != nil {
					return err
				}
				fmt.Printf("installed %s (%s)\n", proto.ShortDir(dir), name)
			}
		}
		for _, h := range hooks {
			switch {
			case f.set["dry"]:
				fmt.Printf("would %s the hook in %s (%s)\n", sub, proto.ShortDir(h.path()), h.name())
			case sub == "remove":
				if err := h.remove(); err != nil {
					return err
				}
				fmt.Printf("removed the hook from %s\n", proto.ShortDir(h.path()))
			default:
				if err := h.install(); err != nil {
					return err
				}
				fmt.Printf("installed the hook in %s (%s)\n", proto.ShortDir(h.path()), h.name())
			}
		}
		if sub == "setup" && !f.set["dry"] && len(hooks) > 0 {
			fmt.Println("\nAgents started from now on in Kano learn which session they're in.")
			for _, h := range hooks {
				if note := map[string]string{
					"claude":   "Claude Code runs hooks only in folders you trust.",
					"codex":    "Codex asks you to approve a new hook: run /hooks in it.",
					"opencode": "OpenCode's background service loads plugins as it starts: run `opencode service restart` when no session is busy.",
				}[h.name()]; note != "" {
					fmt.Println(note)
				}
			}
		}
		return nil
	}
	return usageErr("kano skill setup|show|status|remove [--agent claude|codex|pi|opencode|agents] [--dry-run]")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func skillState(dir string) string {
	b, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	switch {
	case err != nil:
		return "missing"
	case string(b) != skillGuide:
		return "outdated"
	}
	return "current"
}

func writeSkill(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillGuide), 0o644)
}

// hooksFor returns the hooks that go with an --agent: Codex's, Pi's and
// OpenCode's with the skill they share.
func hooksFor(agent string) []integration {
	switch agent {
	case "", "claude", "codex", "pi", "opencode", "agents":
		return agentHooks(agent)
	}
	return nil
}

// skillDirs returns where agents read skills: Claude Code's, and the
// shared one of Codex, OpenCode and others.
func skillDirs(agent string) map[string]string {
	home := proto.Home()
	claude := os.Getenv("CLAUDE_CONFIG_DIR")
	if claude == "" {
		claude = filepath.Join(home, ".claude")
	}
	all := map[string]string{
		"claude": filepath.Join(claude, "skills", "kano"),
		"agents": filepath.Join(home, ".agents", "skills", "kano"),
	}
	if agent = strings.ToLower(agent); agent != "" {
		if d, ok := all[agent]; ok {
			return map[string]string{agent: d}
		}
		if agent == "codex" || agent == "opencode" || agent == "pi" {
			return map[string]string{"agents": all["agents"]}
		}
	}
	return all
}

// SkillState is where a skill for agents is, and whether it is current,
// with the hooks that go with it.
type SkillState struct {
	Agent string // "claude", or "agents" (Codex, OpenCode, Pi…)
	Dir   string
	State string // "current", "outdated" (or partly installed) or "missing"
}

// SkillStates returns the state of kano's skill and hooks for each kind
// of agent.
func SkillStates() []SkillState {
	var out []SkillState
	for _, agent := range []string{"claude", "agents"} {
		dir := skillDirs(agent)[agent]
		states := []string{skillState(dir)}
		for _, h := range hooksFor(agent) {
			states = append(states, h.state())
		}
		st := "outdated"
		switch {
		case !slices.ContainsFunc(states, func(s string) bool { return s != "current" }):
			st = "current"
		case !slices.ContainsFunc(states, func(s string) bool { return s != "missing" }):
			st = "missing"
		}
		out = append(out, SkillState{Agent: agent, Dir: dir, State: st})
	}
	return out
}

// InstallSkill installs, or updates, kano's skill and hooks for a kind of
// agent.
func InstallSkill(agent string) error {
	dir := skillDirs(agent)[agent]
	if dir == "" {
		return fmt.Errorf("no agent %q", agent)
	}
	if err := writeSkill(dir); err != nil {
		return err
	}
	for _, h := range hooksFor(agent) {
		if err := h.install(); err != nil {
			return err
		}
	}
	return nil
}

// RemoveSkill removes kano's skill and hooks for a kind of agent.
func RemoveSkill(agent string) error {
	dir := skillDirs(agent)[agent]
	if dir == "" {
		return fmt.Errorf("no agent %q", agent)
	}
	for _, h := range hooksFor(agent) {
		if err := h.remove(); err != nil {
			return err
		}
	}
	return os.RemoveAll(dir)
}
