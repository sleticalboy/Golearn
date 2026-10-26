# Webcast schema probe

Tooling for CDN schema refresh. **Agent playbook:** [`apps/dy-wss/AGENTS.md`](../apps/dy-wss/AGENTS.md).

```bash
# from repo root
python _schema_probe/fetch_latest.py              # → latest/
python _schema_probe/fetch_latest.py --copy-assets
python _schema_probe/dump_maps2.py
python _schema_probe/diff_decode_maps.py
```

`fetch_latest.py`: live home → `runtime~client-entry.*.js` → chunk maps → download `live-schema*` / `transport-schema*`.

`paths.py` auto-picks newest matching files under `assets/` and `latest/`.
