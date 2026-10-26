# Webcast schema sources

Agent workflow: [`../AGENTS.md`](../AGENTS.md) (repo path `apps/dy-wss/AGENTS.md`).

Truth source for IM field numbers (CDN webpack chunks; refresh via `_schema_probe/fetch_latest.py`):

- `live-schema-im.*.js` — IM / chat / gift / member messages
- `live-schema.*.js` (not `-im`) — data helpers
- `transport-schema*.js` — transport

Current checked-in hashes (2026-10):

- `live-schema-im.e322bd8a.js`
- `live-schema.66b4f47c.js`
- `transport-schema.32356986.js` / `transport-schema-im.63ff9a29.js`

Reference:

- `dy_live.proto.txt` — nested types / field typing
- `_schema_probe/` — fetch / dump / diff

Deprecated (do not use for new fields):

- `live-schema.f30b05b0.js` — old combined dump
- `transport-schema.f54f2795.js` — old transport
