"""Compare helper: same inputs as Go live_info_compare_test.go (no network)."""
import hashlib

UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0"
JS_PATH = __file__.replace("testdata_compare_sig.py", "assets/webmssdk.js")
JS = open(JS_PATH, encoding="utf-8").read()


def get_signature(x_ms_stub: str, user_agent: str) -> str:
    import jsengine

    js_dom = f"""
document = {{}}
window = {{}}
navigator = {{
'userAgent': '{user_agent}'
}}
""".strip()
    try:
        ctx = jsengine.jsengine()
        ctx.eval(js_dom + JS)
        return ctx.eval(f"get_sign('{x_ms_stub}')")
    except BaseException as e:
        print("PY ERR:", e)
        return "00000000"


ordered = [
    ("live_id", "1"),
    ("aid", "6383"),
    ("version_code", "180800"),
    ("webcast_sdk_version", "1.0.14-beta.0"),
    ("room_id", "7519127586287504191"),
    ("sub_room_id", ""),
    ("sub_channel_id", ""),
    ("did_rule", "3"),
    ("user_unique_id", "7692890738690032830"),
    ("device_platform", "web"),
    ("device_type", ""),
    ("ac", ""),
    ("identity", "audience"),
]
stub = hashlib.md5(",".join([f"{k}={v}" for k, v in ordered]).encode()).hexdigest()
synth = "0123456789abcdef0123456789abcdef"
print("PY UA=", UA)
print("PY stub_from_params=", stub)
print("PY sig_params=", get_signature(stub, UA))
print("PY sig_synth=", get_signature(synth, UA))
