package daemon

import (
	"log"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/sync"
)

func Sync(files ...string) error {
	syncConfig := config.SyncConfig()

	// stop if synchronization is not enabled (no config/bad config)
	if syncConfig.Host == "" {
		log.Println("Synchronization not enabled.")
		return nil
	}

	s, err := sync.NewClient(sync.SyncConfig(syncConfig))
	if err != nil {
		return err
	}

	if err := s.Update(files...); err != nil {
		return err
	}

	log.Println("Done syncing.")

	return nil
}
