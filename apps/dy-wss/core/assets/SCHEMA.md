# Webcast schema sources

Truth source for IM field numbers (CDN webpack chunks, 2026-10):

- `live-schema-im.e322bd8a.js` — IM / chat / gift / member messages
- `live-schema.66b4f47c.js` — data helpers
- `transport-schema.32356986.js` / `transport-schema-im.63ff9a29.js` — transport

Reference:

- `dy_live.proto.txt` — human-readable proto dump (nested types / field typing)
- `../../../../_schema_probe/` — diff/dump scripts against these assets

Deprecated (do not use for new fields):

- `live-schema.f30b05b0.js` — old combined dump
- `transport-schema.f54f2795.js` — old transport
