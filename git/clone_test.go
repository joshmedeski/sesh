package git

import (
	"strings"
	"testing"
)

func TestFindRepo(t *testing.T) {
	repos := []string{
		"https://github.com/username/repository.git",
		"git@github.com:username/repository.git",
		"https://github.com/username/repository",
		"git@github.com:username/repository",
		"username/repository",
	}

	for _, repo := range repos {
		result := strings.TrimSpace(findRepo(repo))
		if result != "repository" {
			t.Errorf("Expected repository for URL %s, got %s instead", repo, result)
		}
	}
}
