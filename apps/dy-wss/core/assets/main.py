from collections import defaultdict
import re
from typing import List, Tuple

msg_cls = re.compile(
    r"[a-zA-Z_$][\w$]*\.[A-Z][a-zA-Z0-9]*\s*=\s*function\s*\(\s*\)\s*\{"
)
msg_field = re.compile(r"^[0-9a-fA-F]+: \[")

go_types = {
    "int64String": "uint64",
    "uint64String": "uint64",
    "string": "string",
    "int32": "uint32",
}


class ProtoMsg:

    def __init__(self, first: str):
        self.lines: List[str] = [first]
        self.children: List[ProtoMsg] = []

    def add_child(self, child: "ProtoMsg"):
        self.children.append(child)

    def __repr__(self):
        if len(self.children) == 0:
            return "\n".join(self.lines)
        return (
            "\n".join(self.lines)
            + "\n-----------"
            + "\n".join(map(repr, self.children))
            + "\n----------"
        )


def parse_proto_msg(msg: "ProtoMsg", level: int = 0) -> Tuple[str, str]:
    msg_name = msg.lines[0][2:].split(" ")[0]
    # print(f'parse proto msg: {msg}')
    msg_def = f"{'  ' * level}message {msg_name} " + "{\n"
    # 先解析消息本身
    for line in msg.lines[1:]:
        seq = line.split(":")[0]
        body = line[len(seq) + 3 :].replace(",", "").rstrip("]")
        try:
            seq = int(float(seq))
            # print(f' {int(float(seq))} -> {body}')
        except BaseException as e:
            print(f"!!!!!!!!! {e} !!!!!!!!!")
            full = "\n".join(msg.lines)
            print(f"!!!!!!!!! {full} !!!!!!!!!\n\n")
            break
        segments = body.split(" ")
        field_name = segments[0].strip('"')
        field_type = segments[1][2:]
        repeated = segments[2] == "3"
        if field_type in go_types:
            field_type = go_types[field_type]
        elif field_type.endswith(".decode"):
            field_type = field_type.split(".")[-2]
        # print(f' {seq} -> {field_name} {field_type}')
        if repeated:
            msg_def += (
                f"{'  ' * (level+1)}repeated {field_type} {field_name} = {seq};\n"
            )
        else:
            msg_def += f"{'  ' * (level+1)}{field_type} {field_name} = {seq};\n"
    # 再解析内嵌消息
    for child in msg.children:
        _, child_def = parse_proto_msg(child, level+1)
        msg_def += "\n" + child_def
    msg_def += "  " * level + "}\n"
    return msg_name, msg_def


def parse_one_message(lines: List[str]) -> "ProtoMsg":
    print("parse_one_message ===>\n", ''.join(lines), "\n\n")
    parent = ProtoMsg(lines[0].strip())
    start = 0
    ldepth = 0
    for index, line in enumerate(lines[1:]):
        if start == 0 and msg_field.match(line.strip()):
            parent.lines.append(line.strip())
            continue
        depth = (len(line) - len(line.lstrip())) // 8
        if start > 0 and ldepth > depth:
            continue
        if msg_cls.match(line.strip()):
            if start == 0 or ldepth == 0:
                start = index
                ldepth = depth
                continue
            if start > 0 and ldepth == depth:
                new_lines = [cl[depth*8:] for cl in lines[start+1:index+1]]
                new_lines[0] = new_lines[0].lstrip()
                # print(f'new_lines: ====> \n', ''.join(new_lines))
                parent.add_child(parse_one_message(new_lines))
                start = index
        pass
    if 0 < start < len(lines):
        new_lines = lines[start+1:]
        new_lines[0] = new_lines[0].lstrip()
        # print(f'new_lines: ====> \n', ''.join(new_lines))
        parent.add_child(parse_one_message(new_lines))
        pass
    # print("===>", parent, "\n\n")
    return parent


def do_parse():
    # with open(f"a.txt", "r", encoding="utf-8") as fp:
    with open(f"live-schema.f30b05b0.js.strip.txt", "r", encoding="utf-8") as fp:
        full_lines = fp.readlines()
    messages: List[ProtoMsg] = []
    start = 0
    for index, line in enumerate(full_lines):
        if index == 0 or msg_field.match(line.strip()):
            continue
        if msg_cls.match(line):
            msg = parse_one_message(full_lines[start:index])  # 同级的消息
            messages.append(msg)
            start = index
        pass
    if start < len(full_lines):
        msg = parse_one_message(full_lines[start:])
        messages.append(msg)
    protos = {}
    with open("dy_live.proto", "w", encoding="utf-8") as f:
        f.write('syntax = "proto3";\n\n')
        f.write('option go_package = "../core2";\n\n')
        for msg in messages:
            msg_name, msg_def = parse_proto_msg(msg)
            protos[msg_name] = msg_def
            f.write(msg_def + "\n\n")
        pass
    print(f"messages: {len(messages)}")
    pass


def do_parse_js():
    with open(f"live-schema.f30b05b0.js", "r", encoding="utf-8") as fp:
        full_lines = fp.readlines()
    msg_started, dec_started = False, False
    messages: List[ProtoMsg] = []
    parent: ProtoMsg = None
    child: ProtoMsg = None
    for line in full_lines:
        line = line.strip()
        if line.startswith("//"):
            continue
        if msg_cls.match(line):
            # 判断是否是内嵌的消息
            if msg_started:
                child = ProtoMsg(line)
                parent.add_child(child)
            else:
                msg_started = True
                parent = ProtoMsg(line)
            continue
        if msg_started and line.startswith("e.decode"):
            dec_started = True
            continue
        if dec_started:
            if msg_field.match(line):
                if child:
                    child.lines.append(line)
                else:
                    parent.lines.append(line)
            if msg_started and dec_started and line == "}(),":
                # 如果 child 不为 None，说明一条嵌套消息处理完了
                if child:
                    print(
                        f">>{parent.lines[0]} {len(parent.children)} << \t hit child {child}"
                    )
                    child = None
                else:  # 否则，说明一条消息处理完了
                    msg_started = False
                    dec_started = False
                    messages.append(parent)
                    # print(len(messages), parent, "\n\n")
        pass
    protos = {}
    with open("dy_live.proto", "w", encoding="utf-8") as f:
        f.write('syntax = "proto3";\n\n')
        f.write('option go_package = "../core2";\n\n')
        for msg in messages:
            msg_name, msg_def = parse_proto_msg(msg)
            protos[msg_name] = msg_def
            f.write(msg_def + "\n\n")
        pass
    print(f"messages: {len(protos)}")
    pass


if __name__ == "__main__":
    # do_parse_js()
    do_parse()
    pass
