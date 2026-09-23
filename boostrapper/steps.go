package bootstrap

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// InitStep is a single, named action executed against a Context while bootstrapping a project.
type InitStep struct {
	Name string
	Run  func(ctx *StepExecutionContext) error
}

// FileToGenerate maps an embedded template name to the relative path it produces in the target project.
type FileToGenerate struct {
	Template string
	Output   string
}

// templateFuncs exposes helpers used by the embedded project templates.
var templateFuncs = template.FuncMap{
	"envPrefix": func(appName string) string {
		return strings.ToUpper(strings.NewReplacer("-", "_", " ", "_").Replace(appName))
	},
}

func StepInitOutputDirectory() InitStep {
	return InitStep{
		Name: "Initialize bootstrap",
		Run: func(ctx *StepExecutionContext) error {
			if err := ForceMkDirAllAndWipeBeforeIfNeeded(ctx.OutputDir, ctx.WipeOutputDir); err != nil {
				return err
			}

			return nil
		},
	}
}

func StepRenderTemplates(filesToGenerate []FileToGenerate, templates *embed.FS, data any) InitStep {
	return InitStep{
		Name: "Generate templates",
		Run: func(ctx *StepExecutionContext) error {

			for _, file := range filesToGenerate {
				targetPath := filepath.Join(ctx.OutputDir, file.Output)

				if !ctx.ForceOverwrite {
					if _, err := os.Stat(targetPath); err == nil {
						return fmt.Errorf("file %q already exists (use --force to overwrite)", targetPath)
					}
				}

				tmpl, err := template.New(filepath.Base(file.Template)).Funcs(templateFuncs).ParseFS(templates, file.Template)
				if err != nil {
					return fmt.Errorf("unable to parse template %q: %w", file.Template, err)
				}

				baseTargetPath := filepath.Dir(targetPath)
				if err := ForceMkDirAll(baseTargetPath); err != nil {
					return err
				}

				target, err := os.Create(targetPath)
				if err != nil {
					return fmt.Errorf("unable to create file %q: %w", targetPath, err)
				}

				if err := tmpl.Execute(target, data); err != nil {
					target.Close()
					return fmt.Errorf("unable to render template %q: %w", file.Template, err)
				}
				target.Close()
			}

			return nil
		},
	}
}
