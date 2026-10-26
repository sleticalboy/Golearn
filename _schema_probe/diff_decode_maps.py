"""Diff decode field maps: deprecated live-schema vs current live-schema-im."""
import re
from pathlib import Path

from paths import LIVE_SCHEMA_IM, LIVE_SCHEMA_OLD


def extract_decode_map(src: str, full_name: str) -> dict[str, str]:
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


old = LIVE_SCHEMA_OLD.read_text(encoding="utf-8", errors="ignore")
new = LIVE_SCHEMA_IM.read_text(encoding="utf-8", errors="ignore")

targets = [
    "r.webcast.im.MemberMessage",
    "r.webcast.im.MemberMessage.EffectConfig",
    "r.webcast.im.GiftMessage",
    "r.webcast.im.ChatMessage",
    "r.webcast.data.GiftStruct",
    "r.webcast.im.Common",
    "r.webcast.im.PublicAreaCommon",
]

print("old:", LIVE_SCHEMA_OLD.name)
print("new:", LIVE_SCHEMA_IM.name)
print()

for t in targets:
    mo, mn = extract_decode_map(old, t), extract_decode_map(new, t)
    only_old = sorted(set(mo) - set(mn), key=keynum)
    only_new = sorted(set(mn) - set(mo), key=keynum)
    renamed = [(n, mo[n], mn[n]) for n in sorted(set(mo) & set(mn), key=keynum) if mo[n] != mn[n]]
    print(f"=== {t} old={len(mo)} new={len(mn)} ===")
    print("  only_old:", [(n, mo[n]) for n in only_old])
    print("  only_new:", [(n, mn[n]) for n in only_new])
    if renamed:
        print("  renamed:", renamed)
    print()
