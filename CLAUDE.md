# CLAUDE.md

See `README.md` for setup, configuration, architecture, and migration commands. This file covers only what the README
doesn't.

## What we're building

Telegnome is a Telegram bot that guards group chats against spam. Admins add it to a group, grant it admin rights, and
it keeps unwanted messages and accounts out. The repo started from a generic Go bot template, so some of the plumbing
(WebApp server, T&C gate, Redis conversation storage) predates the product. Reuse it where it fits rather than building
a parallel mechanism, and prefer extending an existing pattern over introducing a new one.

## Voice

Telegnome is a small, friendly garden gnome guarding the group. User-facing texts are warm, short, and a little
whimsical (an occasional 🍄, 🌱 or 🪓 is fine). Moderation notices stay clear and unambiguous first, cute second:
whoever reads one must understand what happened and why. Never mock the person being moderated.

## Languages

The pilot supports English (`en`) and Russian (`ru`). Every user-facing string goes in **both** `locale/en.yaml` and
`locale/ru.yaml` (and both `*_commands.yaml` for command descriptions) in the same change. Never hard-code user-facing
text in Go. Group-wide messages use the group's configured language, not the sender's `language_code`.

## Verify before reporting done

```shell
go build ./... && go vet ./... && gofmt -l .   # gofmt must print nothing
```

Write `_test.go` tests for non-trivial new logic and run `go test ./...`. Spam-detection rules in particular need
table-driven tests covering both what they catch and what they must let through.

## Keep the README in sync

Update the relevant README section in the same change when adding or altering:

- a config key → Configuration table (also `config.dist.yaml`, a getter in `configuration/config.go`, and
  `requiredConfig` in `main.go` if required)
- a handler → handler-execution-order table (also register it in `main.go` at the right group)
- a dependency, migration, or locale file

## Handler groups

Registration order in `main.go` doesn't matter; the group number does. -2/-1 are context enrichment, 0 is the T&C gate,
2+ are commands.

The T&C gate only covers **private** chats. Handlers that act on group messages (moderation, anti-spam) must not assume
the sender has accepted T&C, or has ever talked to the bot at all.

## Moderation safety

- Never act against chat administrators or the group owner.
- Check the bot's own admin rights before restricting, banning or deleting; when a right is missing, degrade gracefully
  and tell the admins once instead of erroring on every message.
- Prefer reversible actions (delete, restrict, mute) over permanent bans unless the group's settings say otherwise.
- Group chats share a 20 requests/minute outbound limit (see Rate Limiting in the README). Batch or skip notices
  instead of announcing every deleted message.
