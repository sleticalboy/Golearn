"""Scan local JS dumps in this dir for schema chunk name hints."""
import re
from pathlib import Path

d = Path(__file__).resolve().parent
pats = [
    re.compile(rb"live-schema.{0,80}"),
    re.compile(rb"transport-schema.{0,80}"),
    re.compile(rb"schema\.[a-f0-9]{6,12}\.js"),
    re.compile(rb"webcast-platform-cdn.{0,160}schema.{0,80}"),
    re.compile(rb"MemberMessage"),
    re.compile(rb"protobuf|proto\.decode|webcast\.im"),
]
for f in sorted(d.glob("*.js")):
    data = f.read_bytes()
    print("FILE", f.name, "len", len(data))
    for pat in pats:
        ms = pat.findall(data)
        if not ms:
            continue
        print(" ", pat.pattern[:50], "count", len(ms))
        for m in ms[:10]:
            print("  ", m[:140])
