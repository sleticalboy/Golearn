import re


msg_cls = re.compile(
    r"[a-zA-Z_$][\w$]*\.[A-Z][a-zA-Z0-9]*\s*=\s*function\s*\(\s*\)\s*\{"
)
msg_field = re.compile(r"^[0-9a-fA-F]+: \[")


with open("live-schema.f30b05b0.js.txt") as f:
    lines = f.readlines()

with open("live-schema.f30b05b0.js.strip.txt", "w") as f:
    for i, line in enumerate(lines):
        temp = line.strip()
        if temp == "e" or temp == ",":
            continue
        if (
            False
            # or temp.startswith("e.decode")
            # or temp.startswith("for (var")
            or msg_cls.match(temp)
            or msg_field.match(temp)
        ):
            f.write(line)
