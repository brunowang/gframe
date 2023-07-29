package gen

type (
	Option func(options *genOptions)

	genOptions struct {
		components []string
		project    string
		pbGoDir    string
		modPath    string
	}
)

func WithComponents(components []string) Option {
	return func(options *genOptions) {
		for _, com := range components {
			if len(com) > 0 {
				options.components = append(options.components, com)
			}
		}
	}
}

func WithProject(project string) Option {
	return func(options *genOptions) {
		if len(project) > 0 {
			options.project = project
		}
	}
}

func WithPbGoDir(pbGoDir string) Option {
	return func(options *genOptions) {
		if len(pbGoDir) > 0 {
			options.pbGoDir = pbGoDir
		}
	}
}

func WithModPath(modPath string) Option {
	return func(options *genOptions) {
		if len(modPath) > 0 {
			options.modPath = modPath
		}
	}
}
