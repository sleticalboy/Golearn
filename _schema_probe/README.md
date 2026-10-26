# Webcast schema probe

Scripts for comparing Douyin live webpack schema chunks against `apps/dy-wss` protos.

Canonical schema copies live under `apps/dy-wss/core/assets/` (see `SCHEMA.md` there).
`dy_live.proto.txt` in that folder is a human-readable field dump for typing nested messages.

```bash
# from repo root
python _schema_probe/dump_maps2.py
python _schema_probe/diff_decode_maps.py
```
