package bootstrap

// Options for the "bootstrap go" generator itself.
type BaseOptions struct {
	ProjectName   string `arg:"long=name,env=BOOTSTRAPPER_PROJECT_NAME"`                   // project name (if needed in templates)
	OutputDir     string `arg:"long=output,env=BOOTSTRAPPER_OUTPUT_DIR"`                   // where the project skeleton is written
	WipeOutputDir bool   `arg:"long=wipe,env=BOOTSTRAPPER_WIPE_OUTPUT_DIR,default=false"`  // overwrite existing files if set
	Force         bool   `arg:"long=force,env=BOOTSTRAPPER_FORCE_OVERWRITE,default=false"` // overwrite existing files if set
	LogQuiet      bool   `arg:"long=quiet,env=BOOTSTRAPPER_LOG_QUIET,default=false"`       // disable all logs (for automated executions, ...)
	LogDebug      bool   `arg:"long=verbose,env=BOOTSTRAPPER_LOG_DEBUG,default=false"`     // activates debug logs
}

type Options interface {
	GetLogDebug() bool
	GetLogQuiet() bool
	GetOutputDir() string
	GetForceOverwrite() bool
	GetWipeOutputDir() bool
	GetProgramName() string
	CheckConsistency() error
}

func (o BaseOptions) GetWipeOutputDir() bool {
	return o.WipeOutputDir
}

func (o BaseOptions) GetLogDebug() bool {
	return o.LogDebug
}

func (o BaseOptions) GetOutputDir() string {
	return o.OutputDir
}

func (o BaseOptions) GetProgramName() string {
	return o.ProjectName
}

func (o BaseOptions) GetLogQuiet() bool {
	return o.LogQuiet
}

func (o BaseOptions) GetForceOverwrite() bool {
	return o.Force
}

func (o BaseOptions) CheckConsistency() error {
	/*
		if o.ModulePath == "" {
			return fmt.Errorf("--module is required (e.g. --module github.com/SR-G/my-new-app)")
		}
	*/
	return nil
}
