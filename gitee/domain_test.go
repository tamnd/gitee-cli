package gitee

import (
	"testing"

	"github.com/tamnd/any-cli/kit"
)

// These tests are offline: they cover the kit domain's pure string functions
// and the host wiring, which need no network.

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "gitee" {
		t.Errorf("Scheme = %q, want gitee", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != "gitee.com" {
		t.Errorf("Hosts = %v, want [gitee.com]", info.Hosts)
	}
	if info.Identity.Binary != "gitee" {
		t.Errorf("Identity.Binary = %q, want gitee", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	cases := []struct{ in, typ, id string }{
		{"torvalds", "user", "torvalds"},
		{"tamnd/gitee-cli", "repo", "tamnd/gitee-cli"},
		{"https://gitee.com/tamnd/gitee-cli", "repo", "tamnd/gitee-cli"},
		{"https://gitee.com/tamnd", "user", "tamnd"},
	}
	for _, tc := range cases {
		typ, id, err := Domain{}.Classify(tc.in)
		if err != nil || typ != tc.typ || id != tc.id {
			t.Errorf("Classify(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tc.in, typ, id, err, tc.typ, tc.id)
		}
	}
}

func TestClassifyBad(t *testing.T) {
	_, _, err := Domain{}.Classify("")
	if err == nil {
		t.Error("Classify('') expected error, got nil")
	}
}

func TestLocate(t *testing.T) {
	cases := []struct{ typ, id, want string }{
		{"user", "tamnd", "https://gitee.com/tamnd"},
		{"repo", "tamnd/gitee-cli", "https://gitee.com/tamnd/gitee-cli"},
	}
	for _, tc := range cases {
		got, err := Domain{}.Locate(tc.typ, tc.id)
		if err != nil || got != tc.want {
			t.Errorf("Locate(%q, %q) = (%q, %v), want (%q, nil)", tc.typ, tc.id, got, err, tc.want)
		}
	}
}

func TestLocateBadType(t *testing.T) {
	_, err := Domain{}.Locate("unknown", "tamnd")
	if err == nil {
		t.Error("Locate(unknown) expected error, got nil")
	}
}

func TestParseRepoRef(t *testing.T) {
	cases := []struct{ ref, owner, name string }{
		{"tamnd/gitee-cli", "tamnd", "gitee-cli"},
		{"https://gitee.com/tamnd/gitee-cli", "tamnd", "gitee-cli"},
	}
	for _, tc := range cases {
		o, n, err := parseRepoRef(tc.ref)
		if err != nil || o != tc.owner || n != tc.name {
			t.Errorf("parseRepoRef(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tc.ref, o, n, err, tc.owner, tc.name)
		}
	}
}

func TestResolveOn(t *testing.T) {
	h, err := kit.Open()
	if err != nil {
		t.Fatal(err)
	}
	got, err := h.ResolveOn("gitee", "tamnd")
	if err != nil || got.String() != "gitee://user/tamnd" {
		t.Errorf("ResolveOn = (%q, %v), want gitee://user/tamnd", got.String(), err)
	}
}
