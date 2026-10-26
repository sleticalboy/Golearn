"""Dump full field maps from current live-schema-im (and GiftStruct)."""
import re

from paths import LIVE_SCHEMA_IM


def extract(src: str, full_name: str) -> dict[str, str]:
    pat = (
        rf"new\s+{re.escape(full_name)}\s*\(\s*{re.escape(full_name)}\.prototype\s*\)"
        rf"\s*,\s*[a-z]\s*=\s*\{{"
    )
    m = re.search(pat, src)
    if not m:
        m = re.search(rf"new\s+{re.escape(full_name)}\s*\([^)]*\)\s*,\s*[a-z]\s*=\s*\{{", src)
    if not m:
        return {}
    body_start = m.end()
    depth = 1
    j = body_start
    while j < len(src) and depth:
        c = src[j]
        if c == "{":
            depth += 1
        elif c == "}":
            depth -= 1
        j += 1
    body = src[body_start : j - 1]
    fields = re.findall(r"""(\d+|5e3):\s*\[\s*["']([^"']+)["']""", body)
    return {num: name for num, name in fields}


def keynum(x: str) -> int:
    return 5000 if x == "5e3" else int(x)


src = LIVE_SCHEMA_IM.read_text(encoding="utf-8", errors="ignore")
print("source:", LIVE_SCHEMA_IM.name)

for t in [
    "r.webcast.im.GiftMessage",
    "r.webcast.im.ChatMessage",
    "r.webcast.im.MemberMessage",
    "r.webcast.im.MemberMessage.EffectConfig",
    "r.webcast.im.PublicAreaCommon",
    "r.webcast.im.Common",
    "r.webcast.data.GiftStruct",
]:
    m = extract(src, t)
    print(f"=== {t} count={len(m)} ===")
    for n in sorted(m, key=keynum):
        print(f"  {n}: {m[n]}")
    print()
