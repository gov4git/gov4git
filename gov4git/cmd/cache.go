package cmd

import (
	"context"
	"os"
	"time"

	"github.com/gov4git/lib4git/base"
	lib4git "github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/spf13/cobra"
)

var (
	cacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "Manage the client's local cache",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cacheClearCmd = &cobra.Command{
		Use:   "clear",
		Short: "Clear local client cache",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			must.Assertf(ctx, setup.CacheDir != "", "cache dir not specified in config")
			err := os.RemoveAll(setup.CacheDir)
			must.NoError(ctx, err)
		},
	}

	cacheUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update local cache by prefetching community and user repos",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			must.Assertf(ctx, setup.CacheDir != "", "cache dir not specified in config")
			if cacheUpdateIntervalSeconds > 0 {
				for {
					updateCacheBestEffort(ctx)
					time.Sleep(time.Second * time.Duration(cacheUpdateIntervalSeconds))
				}
			} else {
				updateCacheBestEffort(ctx)
			}
		},
	}
)

func updateCacheBestEffort(ctx context.Context) {
	updateCacheReplicaBestEffort(ctx, lib4git.Address(setup.Gov))
	updateCacheReplicaBestEffort(ctx, lib4git.Address(setup.Organizer.Public))
	updateCacheReplicaBestEffort(ctx, lib4git.Address(setup.Organizer.Private))
	updateCacheReplicaBestEffort(ctx, lib4git.Address(setup.Member.Public))
	updateCacheReplicaBestEffort(ctx, lib4git.Address(setup.Member.Private))
}

func updateCacheReplicaBestEffort(ctx context.Context, addr lib4git.Address) {
	if err := must.Try(func() { lib4git.CloneOne(ctx, addr) }); err != nil {
		base.Infof("best effort cache update for %v failed (%v)", addr, err)
	}
}

var (
	cacheUpdateIntervalSeconds int
)

func init() {
	cacheCmd.AddCommand(cacheClearCmd)

	cacheCmd.AddCommand(cacheUpdateCmd)
	cacheUpdateCmd.Flags().IntVar(&cacheUpdateIntervalSeconds, "seconds", 0, "update cache every N seconds, if set")
}
