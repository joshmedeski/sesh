package model

import "testing"

func TestGithubConfigEffectiveTTL(t *testing.T) {
	thirty := 30
	zero := 0
	cases := []struct {
		name string
		ttl  *int
		want int
	}{
		{"nil defaults to 60", nil, 60},
		{"explicit zero disables", &zero, 0},
		{"explicit value", &thirty, 30},
	}
	for _, c := range cases {
		got := GithubConfig{IssueTTL: c.ttl}.EffectiveTTL()
		if got != c.want {
			t.Errorf("%s: EffectiveTTL() = %d, want %d", c.name, got, c.want)
		}
	}
}
