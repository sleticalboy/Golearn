from pathlib import Path

PROBE_DIR = Path(__file__).resolve().parent
REPO_ROOT = PROBE_DIR.parent
ASSETS = REPO_ROOT / "apps" / "dy-wss" / "core" / "assets"
LATEST = PROBE_DIR / "latest"
PROTO = REPO_ROOT / "apps" / "dy-wss" / "core" / "proto" / "messages.proto"
PB_GO_DIR = REPO_ROOT / "apps" / "dy-wss" / "core"


def _newest(paths: list[Path]) -> Path:
    if not paths:
        raise FileNotFoundError("no matching schema files")
    return max(paths, key=lambda p: p.stat().st_mtime)


def _glob_many(pattern: str) -> list[Path]:
    out: list[Path] = []
    for d in (ASSETS, LATEST):
        if d.is_dir():
            out.extend(d.glob(pattern))
    return out


LIVE_SCHEMA_IM = _newest(_glob_many("live-schema-im.*.js"))
LIVE_SCHEMA = _newest(
    [p for p in _glob_many("live-schema.*.js") if "live-schema-im" not in p.name]
)
LIVE_SCHEMA_OLD = ASSETS / "live-schema.f30b05b0.js"
