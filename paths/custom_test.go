package paths_test

import (
	"os"
	"path/filepath"
	"testing"

	vgtest "code.vegaprotocol.io/shared/libs/test"
	"code.vegaprotocol.io/shared/paths"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomPaths(t *testing.T) {
	t.Run("Getting custom path for cache file succeeds", testGettingCustomPathForCacheFileSucceeds)
	t.Run("Getting custom path for config file succeeds", testGettingCustomPathForConfigFileSucceeds)
	t.Run("Getting custom path for data file succeeds", testGettingCustomPathForDataFileSucceeds)
	t.Run("Getting custom path for state file succeeds", testGettingCustomPathForStateFileSucceeds)
	t.Run("Getting custom path from struct for cache file succeeds", testGettingCustomPathFromStructForCacheFileSucceeds)
	t.Run("Getting custom path from struct for config file succeeds", testGettingCustomPathFromStructForConfigFileSucceeds)
	t.Run("Getting custom path from struct for data file succeeds", testGettingCustomPathFromStructForDataFileSucceeds)
	t.Run("Getting custom path from struct for state file succeeds", testGettingCustomPathFromStructForStateFileSucceeds)
}

func testGettingCustomPathForCacheFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	path, err := paths.CustomCachePathFor(home, "fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "cache", "fake-file.empty"), path)
}

func testGettingCustomPathForConfigFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	path, err := paths.CustomConfigPathFor(home, "fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "config", "fake-file.empty"), path)
}

func testGettingCustomPathForDataFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	path, err := paths.CustomDataPathFor(home, "fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "data", "fake-file.empty"), path)
}

func testGettingCustomPathForStateFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	path, err := paths.CustomStatePathFor(home, "fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "state", "fake-file.empty"), path)
}

func testGettingCustomPathFromStructForCacheFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	customPaths := &paths.CustomPaths{CustomHome: home}
	path, err := customPaths.CachePathFor("fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "cache", "fake-file.empty"), path)
}

func testGettingCustomPathFromStructForConfigFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	customPaths := &paths.CustomPaths{CustomHome: home}
	path, err := customPaths.ConfigPathFor("fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "config", "fake-file.empty"), path)
}

func testGettingCustomPathFromStructForDataFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	customPaths := &paths.CustomPaths{CustomHome: home}
	path, err := customPaths.DataPathFor("fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "data", "fake-file.empty"), path)
}

func testGettingCustomPathFromStructForStateFileSucceeds(t *testing.T) {
	home := vgtest.RandomPath()
	defer os.RemoveAll(home)
	customPaths := &paths.CustomPaths{CustomHome: home}
	path, err := customPaths.StatePathFor("fake-file.empty")
	require.NoError(t, err)
	vgtest.AssertDirAccess(t, filepath.Dir(home))
	assert.Equal(t, filepath.Join(home, "state", "fake-file.empty"), path)
}
