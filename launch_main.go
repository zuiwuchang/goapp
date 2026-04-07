package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"testing"

	_version "github.com/zuiwuchang/goapp/version"

	"github.com/spf13/cobra"
)

const App = `goapp`

func main() {
	var (
		deps, version bool
	)
	var cmd = &cobra.Command{
		Use:   App,
		Short: "A high-performance Go-native scripting engine powered by Yaegi.",
		Long:  `A high-performance Go-native scripting engine powered by Yaegi.`,
		Run: func(cmd *cobra.Command, args []string) {

			print := false
			if version {
				yaegi := `unknow`
				buildInfo, ok := debug.ReadBuildInfo()
				if ok {
					for _, dep := range buildInfo.Deps {
						if dep.Path == "github.com/traefik/yaegi" {
							yaegi = dep.Version
						}
					}
				}
				fmt.Printf(`Platform:
  - %s %s %s
  - yaegi %s
  - `+App+` %s
  - %s
  - %s
`, runtime.GOOS, runtime.GOARCH, runtime.Version(),
					yaegi,
					_version.Version,
					_version.Commit,
					_version.Date,
				)
				print = true
			}
			if deps {
				if print {
					fmt.Println()
				}
				buildInfo, ok := debug.ReadBuildInfo()
				fmt.Println(`Third-party:`)
				if ok {
					for _, dep := range buildInfo.Deps {
						fmt.Println(`  -`, dep.Path, dep.Version)
					}
				}
				print = true
			}
			if !print {
				fmt.Println(`A high-performance Go-native scripting engine powered by Yaegi.

Usage: ` + App + " run [SCRIPT_DIR] [ARGUMENTS...]")
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&deps,
		"deps",
		"d",
		false, "print third-party library version")
	flags.BoolVarP(&version,
		"version",
		"v",
		false, "print platform version")
	cmd.AddCommand(
		createRun(),
		createTest(),
	)
	cmd.Execute()
}
func createRun() *cobra.Command {
	var (
		gopath    = os.Getenv(`GOPATH`)
		tags      []string
		env       []string
		sandboxed bool
	)

	var cmd = &cobra.Command{
		Use:          `run`,
		Short:        `Execute a Go package or script`,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, e := newContext(args[0], gopath)
			if e != nil {
				panic(e)
			}
			i, e := ctx.Create(CreateOptions{
				Args:      args,
				BuildTags: tags,
				Env:       env,
				Sandboxed: sandboxed,
			})
			if e != nil {
				panic(e)
			}
			_, e = i.EvalPath(ctx.path)
			if e != nil {
				panic(e)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&gopath, `gopath`, `P`, gopath, `sets GOPATH for the scripts`)
	flags.StringSliceVarP(&tags, `tags`, `T`, nil, `sets build constraints for the scripts`)
	flags.BoolVarP(&sandboxed, `sandboxed`, `S`, false, `run sandboxed stdlib symbols such as os/exec and environment`)
	flags.StringSliceVarP(&env, `env`, `E`, nil, `environment in the form "key=values"`)
	return cmd
}
func createTest() *cobra.Command {
	var (
		gopath    = os.Getenv(`GOPATH`)
		tags      []string
		env       []string
		sandboxed bool
	)

	var cmd = &cobra.Command{
		Use:          `test`,
		Short:        `Run unit tests (*_test.go) within the script context`,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, e := newContext(args[0], gopath)
			if e != nil {
				panic(e)
			}
			i, e := ctx.Create(CreateOptions{
				BuildTags: tags,
				Env:       env,
				Sandboxed: sandboxed,
			})
			if e != nil {
				panic(e)
			}
			_, e = i.EvalPath(ctx.path)
			if e != nil {
				panic(e)
			}

			var (
				tests      []testing.InternalTest
				benchmarks []testing.InternalBenchmark
				pkgs       = i.Symbols(ctx.path)
			)
			for _, syms := range pkgs {
				for name, sym := range syms {
					switch fun := sym.Interface().(type) {
					case func(*testing.B):
						benchmarks = append(benchmarks, testing.InternalBenchmark{name, fun})
					case func(*testing.T):
						tests = append(tests, testing.InternalTest{name, fun})
					}
				}
			}
			testing.Main(func(pat, str string) (bool, error) {
				return regexp.MatchString(pat, str)
			}, tests, benchmarks, nil)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&gopath, `gopath`, `P`, gopath, `sets GOPATH for the scripts`)
	flags.StringSliceVarP(&tags, `tags`, `T`, nil, `sets build constraints for the scripts`)
	flags.BoolVarP(&sandboxed, `sandboxed`, `S`, false, `run sandboxed stdlib symbols such as os/exec and environment`)
	flags.StringSliceVarP(&env, `env`, `E`, nil, `environment in the form "key=values"`)
	return cmd
}
