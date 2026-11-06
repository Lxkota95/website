A [Model Context Protocol](https://modelcontextprotocol.io) server that lets an AI assistant such
as Claude Code help run a Discord server: look up roles and channels, create and manage roles.

- Built by hand from bytes on stdin: JSON lines, JSON-RPC 2.0, and both MCP protocol versions
  (the stateless 2026-07-28 and the older `initialize` handshake)
- Safe by design: its own least-privilege bot, allowlisted servers, a preview before any change,
  and a reason in Discord's audit log for every action
- Worked through as an offline course: milestones with tests, hints and reference solutions
