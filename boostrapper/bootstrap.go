// Package bootstrap scaffolds a brand new Go project (Makefile, go.mod, README, .gitignore, ...)
// from an empty folder, using this library as the entry point.
//
// It is meant to be run from an empty project folder, without a local go.mod, via:
//
//	go run github.com/SR-G/sul/cmd/bootstrap@latest
package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	argweave "github.com/SR-G/argweave"
	"github.com/rs/zerolog"
)

const (
	INIT_VALID             = -1
	OS_EXIT_OK             = 0
	OS_EXIT_ERROR_OPTIONS  = 1
	OS_EXIT_ERROR_INIT     = 2
	OS_EXIT_ERROR_GENERATE = 3
)

func newLogger(debug, silent bool) zerolog.Logger {
	if silent {
		return zerolog.Nop()
	}
	writer := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false,
		PartsOrder: []string{"level", "message"},
	} // TimeFormat: "15:04:05",
	l := zerolog.New(writer).With().Logger() // Timestamp().
	if debug {
		l = l.Level(zerolog.DebugLevel)
	} else {
		l = l.Level(zerolog.InfoLevel)
	}
	return l
}

// StepExecutionContext carries the information shared across all init steps.
type StepExecutionContext struct {
	Logger         *zerolog.Logger // Logger to use for all steps
	OutputDir      string
	ForceOverwrite bool
	WipeOutputDir  bool
	ProgramName    string
}

type Bootstrapper[T Options] struct {
	Opts    *T
	Context *StepExecutionContext
}

func (b *Bootstrapper[T]) InitFailOverContext(logger zerolog.Logger) {
	b.Context = &StepExecutionContext{
		Logger: &logger,
	}
}

func (b *Bootstrapper[T]) InitContext(logger zerolog.Logger, options T) error {
	wd, err := os.Getwd()
	if err != nil {
		b.Context = NewContext[T](&logger, options, ".")
		return fmt.Errorf("can't determine current working directory: %w", err)
	} else {
		b.Context = NewContext[T](&logger, options, wd)
	}
	return nil
}

func (b *Bootstrapper[T]) CreateOptions() *T {
	return new(T)
}

func (b *Bootstrapper[T]) Init() (int, error) {
	// Instanciate proper options object
	options := *b.CreateOptions()
	b.Opts = &options

	// Pre-init to be sure to have a valid logger / valid context
	b.InitFailOverContext(newLogger(false, false))

	// Generic parsing of the options
	parser, err := argweave.New(b.Opts, argweave.AppConfig{
		Name:        "bootstrap",
		Description: "Scaffold a new project from a template.",
		Providers:   []argweave.ProviderKind{argweave.ProviderEnv, argweave.ProviderFlags, argweave.ProviderDefault},
	})
	if err != nil {
		return OS_EXIT_ERROR_OPTIONS, err
	}
	if _, err = parser.Parse(os.Args[1:]); err != nil {
		return OS_EXIT_ERROR_OPTIONS, err
	}

	// Real logger / context initialization
	// (now that all options have been parsed without errors)
	err = b.InitContext(newLogger(options.GetLogDebug(), options.GetLogQuiet()), options)
	if err != nil {
		return OS_EXIT_ERROR_INIT, err
	}

	if err := options.CheckConsistency(); err != nil {
		return OS_EXIT_ERROR_OPTIONS, err
	}

	return INIT_VALID, nil
}

func (b *Bootstrapper[T]) Generate(steps []InitStep) int {
	// Execution of all steps, stop after first failure
	for _, step := range steps {
		b.Context.Logger.Info().Msgf("==> %s", step.Name)
		if err := step.Run(b.Context); err != nil {
			b.Context.Logger.Error().Err(err).Msgf("step %q failed", step.Name)
			return OS_EXIT_ERROR_GENERATE
		}
	}
	b.Context.Logger.Info().Msgf("Project %q bootstrapped in %s", b.Context.ProgramName, b.Context.OutputDir)
	return OS_EXIT_OK
}

func ExecuteCommand(outputDir string, commandName string, commandParameter ...string) error {
	cmd := exec.Command(commandName, commandParameter...)
	cmd.Dir = outputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open nemo failed: %w", err)
	}
	return nil

}

// NewContext builds a Context for workDir, detecting the project name from its folder name.
// If workDir is empty, the current working directory is used.
func NewContext[T Options](logger *zerolog.Logger, options T, cwd string) *StepExecutionContext {
	workDir := options.GetOutputDir()
	if workDir == "" {
		workDir = cwd
	}
	name := options.GetProgramName()
	if name == "" {
		name = filepath.Base(workDir)
	}

	return &StepExecutionContext{
		Logger:         logger,
		OutputDir:      workDir,
		ProgramName:    name,
		ForceOverwrite: options.GetForceOverwrite(),
		WipeOutputDir:  options.GetWipeOutputDir(),
	}
}
