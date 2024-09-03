package lister

func srcs(opts ListOptions) []string {
	var srcs []string
	count := 0
	if opts.Tmux {
		count++
	}
	if opts.Tmuxinator {
		count++
	}
	if opts.Config {
		count++
	}
	if opts.Zoxide {
		count++
	}
	if opts.Tmuxinator {
		count++
	}
	if count == 0 {
		return []string{"tmux", "config", "zoxide", "tmuxinator"}
	}
	srcs = make([]string, count)
	i := 0
	if opts.Tmux {
		srcs[i] = "tmux"
		i++
	}
	if opts.Tmuxinator {
		srcs[i] = "tmuxinator"
		i++
	}
	if opts.Config {
		srcs[i] = "config"
		i++
	}
	if opts.Zoxide {
		srcs[i] = "zoxide"
		i++
	}
	if opts.Tmuxinator {
		srcs[i] = "tmuxinator"
		i++
	}
	return srcs
}
