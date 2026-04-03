package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Context struct {
	path, gopath, dir string
	scriptDir         string
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
	strs := strings.Split(gopath, flag)

	abs, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}

	if stat.IsDir() {
		dir = filepath.Dir(abs)

		if len(strs) == 1 && strs[0] == `` {
			strs[0] = magicDir
		} else {
			strs = append(strs, magicDir)
		}
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
	} else {
		abs = filepath.Dir(abs)
	}

	return &Context{
		path:      path,
		gopath:    strings.Join(strs, flag),
		dir:       dir,
		scriptDir: abs,
	}, nil
}

func (c *Context) Open(name string) (fs.File, error) {
	if strings.HasPrefix(name, magicDirSrc) {
		name = filepath.Join(c.dir, name[len(magicDirSrc):])
	}
	return os.Open(name)
}
