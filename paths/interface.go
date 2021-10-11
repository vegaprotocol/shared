package paths

type Paths interface {
	CachePathFor(CachePath) (string, error)
	CacheDirFor(CachePath) (string, error)
	ConfigPathFor(ConfigPath) (string, error)
	ConfigDirFor(ConfigPath) (string, error)
	DataPathFor(DataPath) (string, error)
	DataDirFor(DataPath) (string, error)
	StatePathFor(StatePath) (string, error)
	StateDirFor(StatePath) (string, error)
}

// New instantiates the specific implementation of the Paths interface based on
// the value of the customHome. If a customHome is specified the custom
// implementation CustomPaths is returned, the standard DefaultPaths otherwise.
func New(customHome string) Paths {
	if len(customHome) != 0 {
		return &CustomPaths{
			CustomHome: customHome,
		}
	}

	return &DefaultPaths{}
}
