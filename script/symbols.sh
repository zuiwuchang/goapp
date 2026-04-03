#!/usr/bin/env bash
set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "yaegi extract auto execute helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -g, --get          auto execute go get"
    echo "  -h, --help          help for $Command"
}


ARGS=`getopt -o hg --long help,get -n "$Command" -- "$@"`
eval set -- "${ARGS}"
GET=0
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -g|--get)
            GET=1
            shift
        ;;
        --)
            shift
            break
        ;;
        *)
            echo Error: unknown flag "$1" for "$Command"
            echo "Run '$Command --help' for usage."
            exit 1
        ;;
    esac
done



cd "$BashDir"

if [[ $GET = 1 ]];then
    for name in "${GOLIB[@]}"; do
        echo go get "\"$name\""
        go get "$name"
    done
fi

mkdir ../symbols -p
rm -f "../symbols"/*.go
cd "../symbols"
echo 'package symbols
import "reflect"
var Symbols = make(map[string]map[string]reflect.Value)' > symbols.go
for name in "${GOLIB[@]}"; do
    echo yaegi extract -name symbols "\"$name\""
    yaegi extract -name symbols "$name"
done

echo All success