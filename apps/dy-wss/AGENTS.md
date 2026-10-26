# dy-wss — agent handoff

Douyin live WSS client (Go + goja `webmssdk`). Room URL → page scrape → `get_sign` → WSS → gzip + Protobuf.

## Layout

| Path | Role |
|------|------|
| `main.go` | entry; default room URL arg |
| `core/live_info.go` | `ParseLiveRoomURL` + goja sign |
| `core/dy_live.go` | WSS loop / unpack by `method` |
| `core/proto/messages.proto` | IM message defs (edit here) |
| `core/messages.pb.go` | generated — do not hand-edit |
| `core/assets/` | schema JS + `webmssdk.js` + `dy_live.proto.txt` |
| `core/assets/SCHEMA.md` | which schema files are current |
| `_schema_probe/` | fetch / dump / diff tooling |

## Refresh schema + fill proto (manual codegen)

Proto field **numbers must match** webpack decode maps. Prefer `live-schema-im.*.js` over deprecated `live-schema.f30b05b0.js`.

```bash
# 1) pull CDN schemas → _schema_probe/latest/ (and optionally assets)
python _schema_probe/fetch_latest.py
python _schema_probe/fetch_latest.py --copy-assets

# 2) inspect / diff vs old dump
python _schema_probe/dump_maps2.py
python _schema_probe/diff_decode_maps.py

# 3) patch apps/dy-wss/core/proto/messages.proto
#    - fix wrong field numbers (wire-breaking)
#    - add only_new fields; nested types: see assets/dy_live.proto.txt
#    - update SCHEMA.md filenames if hashes changed

# 4) regenerate pb.go (from proto/)
cd apps/dy-wss/core/proto
protoc --go_out=.. --go_opt=paths=source_relative messages.proto

# 5) smoke
cd <repo>
go build ./apps/dy-wss/...
go run ./apps/dy-wss https://live.douyin.com/<room>
# expect: wss_url + signature, WebSocket connected, WebcastMemberMessage unpack OK
```

**Not automated yet:** step 3 (proto patch). Diff/dump only guide the edit.

## Conventions

- Git: commit when asked; **push only if user says push**.
- Do not commit `_schema_probe/latest/`, page/entry scrapes (gitignored), or `main.exe`.
- Signing: use embedded `webmssdk.js` via goja; do not RE alternative signers unless asked.
- Android: not used here. Frontend scripts: pnpm if any.

## Smoke room

Last known working: `https://live.douyin.com/80017709309` (may go offline; swap URL if needed).
