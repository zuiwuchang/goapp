package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
	"github.com/zuiwuchang/goapp/symbols"
	_version "github.com/zuiwuchang/goapp/version"
)

type Context struct {
	path, gopath string
	// Open magic dir
	magicDir string

	// Script directory
	scriptDir string
}

const magicDir = `"48f8dc68666c4283a8887e0c25303c0b_5182346e341446cb8f4943f254226fbc`

var magicDirSrc = filepath.Join(magicDir, `src`)

func newContext(path, gopath string) (*Context, error) {
	var dir string
	stat, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	flag := `:`
	if runtime.GOOS == `windows` {
		flag = `;`
	}
	var strs []string
	if gopath == `` {
		strs = make([]string, 1)
	} else {
		strs = strings.Split(flag+gopath, flag)
	}

	abs, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	// get magic dir
	var scriptDir string
	if stat.IsDir() {
		scriptDir = abs

		dir = filepath.Join(abs, `..`)
	} else {
		scriptDir = filepath.Dir(abs)
		dir = filepath.Join(abs, `..`, `..`)
	}

	// magic dir append to GOPATH
	strs[0] = magicDir

	cwd, e := os.Getwd()
	if e != nil {
		return nil, e
	}

	path, e = filepath.Rel(cwd, abs)
	if e != nil {
		return nil, e
	}
	if runtime.GOOS == `windows` {
		path = `.\` + path
	} else {
		path = `./` + path
	}
	return &Context{
		path:      path,
		gopath:    strings.Join(strs, flag),
		magicDir:  dir,
		scriptDir: scriptDir,
	}, nil
}

// magic fs.Open
func (c *Context) Open(name string) (fs.File, error) {
	if strings.HasPrefix(name, magicDirSrc) {
		name = filepath.Join(c.magicDir, name[len(magicDirSrc):])
	}
	return os.Open(name)
}

type CreateOptions struct {
	Args      []string
	Env       []string
	BuildTags []string
	Sandboxed bool
}

// create Yaegi
func (c *Context) Create(opts CreateOptions) (*interp.Interpreter, error) {
	i := interp.New(interp.Options{
		Args:                 opts.Args,
		GoPath:               c.gopath,
		SourcecodeFilesystem: c,
		Env:                  opts.Env,
		BuildTags:            opts.BuildTags,
		Unrestricted:         !opts.Sandboxed,
	})
	e := i.Use(stdlib.Symbols)
	if e != nil {
		return nil, e
	}
	if !opts.Sandboxed {
		if e = i.Use(syscall.Symbols); e != nil {
			return nil, e
		}
		if e = os.Setenv("YAEGI_SYSCALL", "1"); e != nil {
			return nil, e
		}

		if e = i.Use(unsafe.Symbols); e != nil {
			return nil, e
		}
		if e = os.Setenv("YAEGI_UNSAFE", "1"); e != nil {
			return nil, e
		}

		if e = i.Use(unrestricted.Symbols); e != nil {
			return nil, e
		}
		if e = os.Setenv("YAEGI_UNRESTRICTED", "1"); e != nil {
			return nil, e
		}
	}

	keys := symbols.Symbols[`github.com/zuiwuchang/gosdk/gosdk`]
	var (
		dir   = c.scriptDir
		yaegi = `unknow`
	)
	keys["AppCommit"] = reflect.ValueOf(&_version.Commit).Elem()
	keys["AppDate"] = reflect.ValueOf(&_version.Date).Elem()
	keys["AppVersion"] = reflect.ValueOf(&_version.Version).Elem()
	osArgs := os.Args
	keys["Args"] = reflect.ValueOf(&osArgs).Elem()
	keys["Dir"] = reflect.ValueOf(&dir).Elem()

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range buildInfo.Deps {
			if dep.Path == "github.com/traefik/yaegi" {
				yaegi = dep.Version
			}
		}
	}
	keys["Yaegi"] = reflect.ValueOf(&yaegi).Elem()
	e = i.Use(symbols.Symbols)
	if e != nil {
		return nil, e

	}

	return i, nil
}
