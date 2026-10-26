from pathlib import Path

PROBE_DIR = Path(__file__).resolve().parent
REPO_ROOT = PROBE_DIR.parent
ASSETS = REPO_ROOT / "apps" / "dy-wss" / "core" / "assets"

# Preferred (current CDN hashes)
LIVE_SCHEMA_IM = ASSETS / "live-schema-im.e322bd8a.js"
LIVE_SCHEMA = ASSETS / "live-schema.66b4f47c.js"
# Deprecated combined dump
LIVE_SCHEMA_OLD = ASSETS / "live-schema.f30b05b0.js"
