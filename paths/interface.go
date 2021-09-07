package paths

type Paths interface {
	CachePathFor(string) (string, error)
	CacheDirFor(string) (string, error)
	ConfigPathFor(string) (string, error)
	ConfigDirFor(string) (string, error)
	DataPathFor(string) (string, error)
	DataDirFor(string) (string, error)
	StatePathFor(string) (string, error)
	StateDirFor(string) (string, error)
}

func NewPaths(customHome string) Paths {
	if len(customHome) != 0 {
		return &CustomPaths{
			CustomHome: customHome,
		}
	}

	return &DefaultPaths{}
}
