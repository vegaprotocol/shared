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

func NewPaths(customHome string) Paths {
	if len(customHome) != 0 {
		return &CustomPaths{
			CustomHome: customHome,
		}
	}

	return &DefaultPaths{}
}
