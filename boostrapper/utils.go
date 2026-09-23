package bootstrap

import (
	"fmt"
	"os"
)

func ForceMkDirAllAndWipeBeforeIfNeeded(directory string, wipe bool) error {
	// Delete previous dir
	// - if --wipe is activated
	// - and if previously found on disk
	if wipe {
		if _, err := os.Stat(directory); err == nil {
			if err = os.RemoveAll(directory); err != nil {
				return fmt.Errorf("can't delete directory %q", directory)
			}
		}
	}

	// In all other situations, recreate the target dir
	// - directory was found on disk, with wipe activated, and has been deleted
	// - directory was found or not
	return ForceMkDirAll(directory)
}

func ForceMkDirAll(directory string) error {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("unable to create output directory %q: %w", directory, err)
	}
	return nil
}
