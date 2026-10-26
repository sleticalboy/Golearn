"""Fetch current Douyin live schema chunks into _schema_probe/latest/.

Flow:
  1. GET https://live.douyin.com/ → find runtime~client-entry.<hash>.js
  2. Parse webpack maps: chunkId→name, chunkId→contentHash
  3. Download name.contentHash.js for live/transport schemas
"""
from __future__ import annotations

import argparse
import re
import sys
import urllib.request
from pathlib import Path

from paths import ASSETS, PROBE_DIR

UA = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
    "(KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
)
LIVE_HOME = "https://live.douyin.com/"
CDN_BASE = (
    "https://lf-webcast-platform.bytetos.com/obj/webcast-platform-cdn/webcast/douyin_live"
)
SCHEMA_NAMES = (
    "live-schema",
    "live-schema-im",
    "transport-schema",
    "transport-schema-im",
)


def http_get(url: str, timeout: int = 60) -> bytes:
    req = urllib.request.Request(url, headers={"User-Agent": UA})
    with urllib.request.urlopen(req, timeout=timeout) as resp:
        return resp.read()


def find_runtime_url(html: str) -> str:
    m = re.search(
        r"https://[^\"'\s]+/webcast/douyin_live/runtime~client-entry\.[a-f0-9]+\.js",
        html,
    )
    if m:
        return m.group(0)
    m = re.search(r"runtime~client-entry\.[a-f0-9]+\.js", html)
    if not m:
        raise RuntimeError("runtime~client-entry.*.js not found on live home page")
    return f"{CDN_BASE}/{m.group(0)}"


def parse_chunk_maps(runtime_js: str) -> dict[str, tuple[str, str]]:
    """Return {schema_name: (chunk_id, content_hash)}."""
    id_by_name: dict[str, str] = {}
    for name in SCHEMA_NAMES:
        m = re.search(rf'(\d+):"{re.escape(name)}"', runtime_js)
        if not m:
            raise RuntimeError(f"chunk id for {name!r} not found in runtime")
        id_by_name[name] = m.group(1)

    out: dict[str, tuple[str, str]] = {}
    for name, cid in id_by_name.items():
        m = re.search(rf'{cid}:"([a-f0-9]+)"', runtime_js)
        if not m:
            raise RuntimeError(f"content hash for chunk {cid} ({name}) not found")
        out[name] = (cid, m.group(1))
    return out


def schema_url(name: str, content_hash: str) -> str:
    # IM schemas sit under chunks/; others historically under douyin_live root.
    # Both paths usually work via CDN; prefer chunks/ then root.
    return f"{CDN_BASE}/chunks/{name}.{content_hash}.js"


def download_schema(name: str, content_hash: str, dest: Path) -> Path:
    candidates = [
        f"{CDN_BASE}/chunks/{name}.{content_hash}.js",
        f"{CDN_BASE}/{name}.{content_hash}.js",
    ]
    last_err: Exception | None = None
    for url in candidates:
        try:
            data = http_get(url)
            dest.write_bytes(data)
            print(f"OK  {name}.{content_hash}.js  {len(data)} bytes  <- {url}")
            return dest
        except Exception as e:  # noqa: BLE001 — try next CDN path
            last_err = e
            print(f"..  miss {url}: {e}")
    raise RuntimeError(f"failed to download {name}: {last_err}")


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument(
        "--copy-assets",
        action="store_true",
        help="also copy downloaded files into apps/dy-wss/core/assets/",
    )
    args = ap.parse_args()

    out_dir = PROBE_DIR / "latest"
    out_dir.mkdir(parents=True, exist_ok=True)

    print("GET", LIVE_HOME)
    html = http_get(LIVE_HOME).decode("utf-8", "ignore")
    runtime_url = find_runtime_url(html)
    print("runtime", runtime_url)

    runtime_js = http_get(runtime_url).decode("utf-8", "ignore")
    (PROBE_DIR / Path(runtime_url).name).write_text(runtime_js, encoding="utf-8")
    maps = parse_chunk_maps(runtime_js)
    for name, (cid, h) in maps.items():
        print(f"  {name}: chunk={cid} hash={h}")

    written: list[Path] = []
    for name, (_cid, h) in maps.items():
        dest = out_dir / f"{name}.{h}.js"
        download_schema(name, h, dest)
        written.append(dest)

    if args.copy_assets:
        ASSETS.mkdir(parents=True, exist_ok=True)
        for p in written:
            target = ASSETS / p.name
            target.write_bytes(p.read_bytes())
            print(f"copied → {target.relative_to(ASSETS.parent.parent.parent.parent)}")

    print("done →", out_dir)
    return 0


if __name__ == "__main__":
    sys.exit(main())
