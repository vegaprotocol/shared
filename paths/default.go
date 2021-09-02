package paths

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
)

// The default Vega file structure is mapped on the XDG standard. This standard
// defines where the files should be looked for, depending on their purpose,
// through environment variables, prefixed by `$XDG_`. The value of these
// variables matches the standards of the platform the program runs on.
//
// More on XDG at:
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
//
// At each location, Vega files are grouped under a `vega` folder, as follows
// `$XDG_*/vega`, before being sorted in sub-folders. The file structure of
// these sub-folder is described in paths.go.
//
// Default file structure:
//
// $XDG_CACHE_HOME
// └── vega
//
// $XDG_CONFIG_HOME
// └── vega
//
// $XDG_DATA_HOME
// └── vega
//
// $XDG_STATE_HOME
// └── vega

type DefaultPaths struct{}

func (p *DefaultPaths) CachePathFor(relFilePath string) (string, error) {
	return DefaultCachePathFor(relFilePath)
}

func (p *DefaultPaths) CacheDirFor(relDirPath string) (string, error) {
	return DefaultCacheDirFor(relDirPath)
}

func (p *DefaultPaths) ConfigPathFor(relFilePath string) (string, error) {
	return DefaultConfigPathFor(relFilePath)
}

func (p *DefaultPaths) ConfigDirFor(relDirPath string) (string, error) {
	return DefaultConfigDirFor(relDirPath)
}

func (p *DefaultPaths) DataPathFor(relFilePath string) (string, error) {
	return DefaultDataPathFor(relFilePath)
}

func (p *DefaultPaths) DataDirFor(relDirPath string) (string, error) {
	return DefaultDataDirFor(relDirPath)
}

func (p *DefaultPaths) StatePathFor(relFilePath string) (string, error) {
	return DefaultStatePathFor(relFilePath)
}

func (p *DefaultPaths) StateDirFor(relDirPath string) (string, error) {
	return DefaultStateDirFor(relDirPath)
}

// DefaultCachePathFor builds the default path for cache files and creates
// intermediate directories, if needed.
func DefaultCachePathFor(relFilePath string) (string, error) {
	path, err := xdg.CacheFile(filepath.Join(VegaHome, relFilePath))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// DefaultCacheDirFor builds the default path for cache files and creates
// intermediate directories, if needed.
func DefaultCacheDirFor(relDirPath string) (string, error) {
	// We append fake-file to xdg library creates all directory up to fake-file.
	path, err := xdg.CacheFile(filepath.Join(VegaHome, relDirPath, "fake-file"))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relDirPath, err)
	}
	return filepath.Dir(path), nil
}

// DefaultConfigPathFor builds the default path for configuration files and
// creates intermediate directories, if needed.
func DefaultConfigPathFor(relFilePath string) (string, error) {
	path, err := xdg.ConfigFile(filepath.Join(VegaHome, relFilePath))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// DefaultConfigDirFor builds the default path for config files and creates
// intermediate directories, if needed.
func DefaultConfigDirFor(relDirPath string) (string, error) {
	// We append fake-file to xdg library creates all directory up to fake-file.
	path, err := xdg.ConfigFile(filepath.Join(VegaHome, relDirPath, "fake-file"))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relDirPath, err)
	}
	return filepath.Dir(path), nil
}

// DefaultDataPathFor builds the default path for data files and creates
// intermediate directories, if needed.
func DefaultDataPathFor(relFilePath string) (string, error) {
	path, err := xdg.DataFile(filepath.Join(VegaHome, relFilePath))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// DefaultDataDirFor builds the default path for data files and creates
// intermediate directories, if needed.
func DefaultDataDirFor(relDirPath string) (string, error) {
	// We append fake-file to xdg library creates all directory up to fake-file.
	path, err := xdg.DataFile(filepath.Join(VegaHome, relDirPath, "fake-file"))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relDirPath, err)
	}
	return filepath.Dir(path), nil
}

// DefaultStatePathFor builds the default path for state files and creates
// intermediate directories, if needed.
func DefaultStatePathFor(relFilePath string) (string, error) {
	path, err := xdg.StateFile(filepath.Join(VegaHome, relFilePath))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// DefaultStateDirFor builds the default path for state files and creates
// intermediate directories, if needed.
func DefaultStateDirFor(relDirPath string) (string, error) {
	// We append fake-file to xdg library creates all directory up to fake-file.
	path, err := xdg.StateFile(filepath.Join(VegaHome, relDirPath, "fake-file"))
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relDirPath, err)
	}
	return filepath.Dir(path), nil
}