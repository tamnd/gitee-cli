package gitee_test

import (
	"testing"

	"github.com/tamnd/gitee-cli/gitee"
)

func TestParseRepoSlug(t *testing.T) {
	tests := []struct {
		input     string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{"owner/repo", "owner", "repo", false},
		{"foo/bar-baz", "foo", "bar-baz", false},
		{"https://gitee.com/owner/repo", "owner", "repo", false},
		{"http://gitee.com/owner/repo", "owner", "repo", false},
		{"https://gitee.com/owner/repo/", "owner", "repo", false},
		// errors
		{"", "", "", true},
		{"noslash", "", "", true},
		{"/repo", "", "", true},
		{"owner/", "", "", true},
	}
	for _, tc := range tests {
		slug, err := gitee.ParseRepoSlug(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseRepoSlug(%q): want error, got %+v", tc.input, slug)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRepoSlug(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if slug.Owner != tc.wantOwner || slug.Repo != tc.wantRepo {
			t.Errorf("ParseRepoSlug(%q) = {%q,%q}, want {%q,%q}", tc.input, slug.Owner, slug.Repo, tc.wantOwner, tc.wantRepo)
		}
	}
}

func TestRepoSlugString(t *testing.T) {
	s := gitee.RepoSlug{Owner: "foo", Repo: "bar"}
	if got := s.String(); got != "foo/bar" {
		t.Errorf("String() = %q, want %q", got, "foo/bar")
	}
}
