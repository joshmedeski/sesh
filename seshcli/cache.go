package seshcli

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
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
		Long: `Fetch live session data and rewrite the cache, so the next
'sesh list' reflects sessions that were created or killed outside
'sesh connect' (which refreshes the cache itself). Intended for tmux
session-created/session-closed hooks and picker kill bindings.

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
			deps.CachingLister.RefreshCache(lister.ListOptions{})
			deps.CachingLister.Wait()
			return nil
		},
	}
}
