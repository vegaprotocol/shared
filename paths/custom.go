package paths

import (
	"fmt"
	"path/filepath"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
)

// When opting for a custom Vega home, the all files are located under the
// specified folder. They are sorted, by purpose, in sub-folders. The structure
// of these sub-folder is described in paths.go.
//
// File structure for custom home:
//
// VEGA_HOME
// 	├── cache/
// 	├── config/
// 	├── data/
// 	└── state/

type CustomPaths struct {
	CustomHome string
}

func (p *CustomPaths) CacheDirFor(relDirPath string) (string, error) {
	return CustomCacheDirFor(p.CustomHome, relDirPath)
}

func (p *CustomPaths) CachePathFor(relFilePath string) (string, error) {
	return CustomCachePathFor(p.CustomHome, relFilePath)
}

func (p *CustomPaths) ConfigDirFor(relDirPath string) (string, error) {
	return CustomConfigDirFor(p.CustomHome, relDirPath)
}

func (p *CustomPaths) ConfigPathFor(relFilePath string) (string, error) {
	return CustomConfigPathFor(p.CustomHome, relFilePath)
}

func (p *CustomPaths) DataDirFor(relDirPath string) (string, error) {
	return CustomDataDirFor(p.CustomHome, relDirPath)
}

func (p *CustomPaths) DataPathFor(relFilePath string) (string, error) {
	return CustomDataPathFor(p.CustomHome, relFilePath)
}

func (p *CustomPaths) StateDirFor(relDirPath string) (string, error) {
	return CustomStateDirFor(p.CustomHome, relDirPath)
}

func (p *CustomPaths) StatePathFor(relFilePath string) (string, error) {
	return CustomStatePathFor(p.CustomHome, relFilePath)
}

// CustomCachePathFor builds the path for cache files at a given root path and
// creates intermediate directories. It scoped the files under a "cache" folder,
// and follow the default structure.
func CustomCachePathFor(customHome, relFilePath string) (string, error) {
	fullPath := filepath.Join(customHome, "cache", relFilePath)
	dir := filepath.Dir(fullPath)
	if err := vgfs.EnsureDir(dir); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", dir, err)
	}
	return fullPath, nil
}

// CustomCacheDirFor builds the path for cache directories at a given root path
// and creates intermediate directories. It scoped the files under a "data"
// folder, and follow the default structure.
func CustomCacheDirFor(customHome, relDirPath string) (string, error) {
	path := filepath.Join(customHome, "cache", relDirPath)
	if err := vgfs.EnsureDir(path); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}

// CustomConfigPathFor builds the path for configuration files at a given root
// path and creates intermediate directories. It scoped the files under a
// "config" folder, and follow the default structure.
func CustomConfigPathFor(customHome, relFilePath string) (string, error) {
	fullPath := filepath.Join(customHome, "config", relFilePath)
	dir := filepath.Dir(fullPath)
	if err := vgfs.EnsureDir(dir); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", dir, err)
	}
	return fullPath, nil
}

// CustomConfigDirFor builds the path for config directories at a given root path
// and creates intermediate directories. It scoped the files under a "data"
// folder, and follow the default structure.
func CustomConfigDirFor(customHome, relDirPath string) (string, error) {
	path := filepath.Join(customHome, "config", relDirPath)
	if err := vgfs.EnsureDir(path); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}

// CustomDataPathFor builds the path for data files at a given root path and
// creates intermediate directories. It scoped the files under a "data" folder,
// and follow the default structure.
func CustomDataPathFor(customHome, relFilePath string) (string, error) {
	fullPath := filepath.Join(customHome, "data", relFilePath)
	dir := filepath.Dir(fullPath)
	if err := vgfs.EnsureDir(dir); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", dir, err)
	}
	return fullPath, nil
}

// CustomDataDirFor builds the path for data directories at a given root path
// and creates intermediate directories. It scoped the files under a "data"
// folder, and follow the default structure.
func CustomDataDirFor(customHome, relDirPath string) (string, error) {
	path := filepath.Join(customHome, "data", relDirPath)
	if err := vgfs.EnsureDir(path); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}

// CustomStatePathFor builds the path for cache files at a given root path and
// creates intermediate directories. It scoped the files under a "cache" folder,
// and follow the default structure.
func CustomStatePathFor(customHome, relFilePath string) (string, error) {
	fullPath := filepath.Join(customHome, "state", relFilePath)
	dir := filepath.Dir(fullPath)
	if err := vgfs.EnsureDir(dir); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", dir, err)
	}
	return fullPath, nil
}

// CustomStateDirFor builds the path for data directories at a given root path
// and creates intermediate directories. It scoped the files under a "data"
// folder, and follow the default structure.
func CustomStateDirFor(customHome, relDirPath string) (string, error) {
	path := filepath.Join(customHome, "state", relDirPath)
	if err := vgfs.EnsureDir(path); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}
