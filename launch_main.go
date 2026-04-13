package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/traefik/yaegi/interp"
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

type Caller struct {
	i    *interp.Interpreter
	path string

	keys map[string]reflect.Value
}

func (c *Caller) Call(name string) error {
	if c.keys == nil {
		pkgs := c.i.Symbols(c.path)
		if pkgs == nil {
			return fmt.Errorf(`func () not found: %s`, name)
		}
		keys := pkgs[c.path]
		if keys == nil {
			return fmt.Errorf(`func () not found: %s`, name)
		}
		c.keys = keys
	}
	if i, ok := c.keys[name]; ok {
		if f, ok := i.Interface().(func()); ok {
			f()
			return nil
		} else {
			return fmt.Errorf(`not a func (): %s`, name)
		}
	}
	return fmt.Errorf(`func () not found: %s`, name)
}

func createRun() *cobra.Command {
	var (
		gopath    = os.Getenv(`GOPATH`)
		tags      []string
		env       []string
		funcNames []string
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
			if len(funcNames) != 0 {
				caller := &Caller{
					i: i,
				}
				if ctx.idDir {
					caller.path = ctx.path
				} else {
					caller.path = `main`
				}
				for _, name := range funcNames {
					e = caller.Call(name)
					if e != nil {
						panic(e)
					}
				}
			} else if !ctx.idDir {
				// auto RunMain
				keys := i.Symbols(`main`)[`main`]
				for k, v := range keys {
					if strings.HasPrefix(k, `RunMain`) {
						if f, ok := v.Interface().(func()); ok {
							f()
						}
					}
				}
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&gopath, `gopath`, `P`, gopath, `sets GOPATH for the scripts`)
	flags.StringSliceVarP(&tags, `tags`, `T`, nil, `sets build constraints for the scripts`)
	flags.BoolVarP(&sandboxed, `sandboxed`, `S`, false, `run sandboxed stdlib symbols such as os/exec and environment`)
	flags.StringSliceVarP(&env, `env`, `E`, nil, `environment in the form "key=values"`)
	flags.StringSliceVarP(&funcNames, `func`, `F`, nil, `to display the name of the function called`)
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
