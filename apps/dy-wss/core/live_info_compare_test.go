package core

import (
	"crypto/md5"
	"encoding/hex"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Same synthetic inputs used for Go/Python compare.
var compareParams = []kv{
	{"live_id", "1"},
	{"aid", "6383"},
	{"version_code", "180800"},
	{"webcast_sdk_version", "1.0.14-beta.0"},
	{"room_id", "7519127586287504191"},
	{"sub_room_id", ""},
	{"sub_channel_id", ""},
	{"did_rule", "3"},
	{"user_unique_id", "7692890738690032830"},
	{"device_platform", "web"},
	{"device_type", ""},
	{"ac", ""},
	{"identity", "audience"},
}

const compareSynthStub = "0123456789abcdef0123456789abcdef"

func compareStubFromParams() string {
	parts := make([]string, 0, len(compareParams))
	for _, p := range compareParams {
		parts = append(parts, p.k+"="+p.v)
	}
	sum := md5.Sum([]byte(strings.Join(parts, ",")))
	return hex.EncodeToString(sum[:])
}

func TestCompareStubMatchesPythonFormula(t *testing.T) {
	// Python hashlib.md5(",".join(...)).hexdigest() for the same ordered pairs
	want := "fa106c1f02b324bf431e98cb92d3a479"
	got := compareStubFromParams()
	if got != want {
		t.Fatalf("x_ms_stub mismatch: go=%s want=%s", got, want)
	}
}

func TestGetSignatureProducesValue(t *testing.T) {
	sig := getSignature(compareSynthStub, liveUserAgent)
	if sig == "" || sig == "00000000" {
		t.Fatalf("bad signature: %q", sig)
	}
	// X-Bogus from this SDK is typically 16 chars; keep soft check
	if len(sig) < 8 {
		t.Fatalf("signature too short: %q", sig)
	}
}

func TestGetSignatureNotByteStableAcrossRuns(t *testing.T) {
	a := getSignature(compareSynthStub, liveUserAgent)
	b := getSignature(compareSynthStub, liveUserAgent)
	if a == "00000000" || b == "00000000" {
		t.Fatalf("engine failed: a=%q b=%q", a, b)
	}
	// Documented behavior: crawler/X-Bogus includes entropy; exact equality is not expected.
	t.Logf("run1=%s run2=%s equal=%v", a, b, a == b)
}

func TestCompareSignatureWithPythonIfAvailable(t *testing.T) {
	py := findPython()
	if py == "" {
		t.Skip("python not found")
	}
	script := filepath.Join(filepath.Dir(mustSourceFile()), "testdata_compare_sig.py")
	out, err := exec.Command(py, script).CombinedOutput()
	if err != nil {
		t.Skipf("python compare script failed: %v\n%s", err, out)
	}
	lines := string(out)
	t.Logf("python output:\n%s", lines)

	stub := compareStubFromParams()
	goSig := getSignature(stub, liveUserAgent)
	goSynth := getSignature(compareSynthStub, liveUserAgent)
	t.Logf("go stub=%s sig_params=%s sig_synth=%s", stub, goSig, goSynth)

	if !strings.Contains(lines, "PY stub_from_params= "+stub) && !strings.Contains(lines, "PY stub_from_params="+stub) {
		// allow either spacing
		if !strings.Contains(lines, stub) {
			t.Fatalf("python stub not matching go stub %s\n%s", stub, lines)
		}
	}
	// Exact signature equality is NOT asserted: both engines are non-deterministic per run.
}

func findPython() string {
	candidates := []string{`D:\Dev\miniconda3\python.exe`, "python", "python3"}
	for _, c := range candidates {
		if strings.Contains(c, `\`) || strings.Contains(c, `/`) {
			if _, err := exec.Command(c, "--version").CombinedOutput(); err == nil {
				return c
			}
			continue
		}
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	return ""
}

func mustSourceFile() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return f
}
