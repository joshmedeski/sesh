package seshcli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewCacheCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage the session list cache",
	}
	cmd.AddCommand(newCacheRefreshCommand(base))
	return cmd
}

func newCacheRefreshCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Rebuild the session list cache from live data",
		Long: `Fetch the slow sources (zoxide and tmuxinator) and rewrite the cache,
so the next 'sesh list' reflects changes made outside 'sesh connect'
(which refreshes the cache itself). tmux sessions are never cached and
need no refresh.

Does nothing when the cache is disabled.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}
			if deps.CachingLister == nil {
				slog.Debug("cache refresh: cache is disabled, nothing to do")
				return nil
			}
			deps.CachingLister.RefreshCache()
			deps.CachingLister.Wait()
			return nil
		},
	}
}
