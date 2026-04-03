Target="goapp"
Package="github.com/zuiwuchang/goapp"
Docker="king011/goapp"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="0.0.1"
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/amd64
)

GOLIB=(
    github.com/zuiwuchang/gosdk
    github.com/spf13/cobra
    github.com/fsnotify/fsnotify
    github.com/jroimartin/gocui
)
