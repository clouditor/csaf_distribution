package main

import (
	"os"
	"path/filepath"

	"github.com/Intevation/csaf_trusted/csaf"
)

func doTransaction(
	cfg *config,
	t tlp,
	fn func(string, *csaf.ProviderMetadata) error,
) error {

	wellknown := filepath.Join(cfg.Web, ".well-known", "csaf")

	metadata := filepath.Join(wellknown, "provider-metadata.json")

	pmd, err := func() (*csaf.ProviderMetadata, error) {
		f, err := os.Open(metadata)
		if err != nil {
			if os.IsNotExist(err) {
				return csaf.NewProviderMetadata(
					cfg.Domain + "/.wellknown/csaf/provider-metadata.json"), nil
			}
			return nil, err
		}
		defer f.Close()
		return csaf.LoadProviderMetadata(f)
	}()

	if err != nil {
		return err
	}

	webTLP := filepath.Join(wellknown, string(t))

	oldDir, err := filepath.EvalSymlinks(webTLP)
	if err != nil {
		return err
	}

	folderTLP := filepath.Join(cfg.Folder, string(t))

	newDir, err := mkUniqDir(folderTLP)
	if err != nil {
		return err
	}

	// Copy old content into new.
	if err := deepCopy(newDir, oldDir); err != nil {
		os.RemoveAll(newDir)
		return err
	}

	// Work with new folder.
	if err := fn(newDir, pmd); err != nil {
		os.RemoveAll(newDir)
		return err
	}

	// Write back provider metadata.
	newMetaName, newMetaFile, err := mkUniqFile(metadata)
	if err != nil {
		os.RemoveAll(newDir)
		return err
	}

	if err := pmd.Save(newMetaFile); err != nil {
		newMetaFile.Close()
		os.Remove(newMetaName)
		os.RemoveAll(newDir)
		return err
	}

	if err := newMetaFile.Close(); err != nil {
		os.Remove(newMetaName)
		os.RemoveAll(newDir)
		return err
	}

	if err := os.Rename(newMetaName, metadata); err != nil {
		os.RemoveAll(newDir)
		return err
	}

	// Switch directories.
	symlink := filepath.Join(newDir, string(t))
	if err := os.Symlink(newDir, symlink); err != nil {
		os.RemoveAll(newDir)
		return err
	}
	if err := os.Rename(symlink, webTLP); err != nil {
		os.RemoveAll(newDir)
		return err
	}

	return os.RemoveAll(oldDir)
}