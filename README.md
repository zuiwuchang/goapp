# GoApp

English Version | [中文說明](./zh.md)


`GoApp` is a powerful wrapper for [Yaegi](https://github.com/traefik/yaegi), a Go engine that enables you to interpret and execute Go code as scripts. Compared to the standard Yaegi, `GoApp` comes with more built-in third-party libraries and provides a streamlined workflow for adding your own dependencies.

## Features

  * **Expanded Standard Library**: Includes a wider range of pre-compiled third-party symbols.
  * **Easy Extension**: Simplifies the process of embedding extra Go libraries into the interpreter.
  * **Smart Module Loading**: Features a custom `GOPATH` hook that allows you to import local projects effortlessly.

-----

## Usage

### Command Syntax

```bash
$ goapp run [flags] [Project_Dir | GO_FILE] -- [options]
```

*Use `--` to separate `goapp` flags from the arguments passed to your script.*

### Flags

  * `-E, --env strings`: Set environment variables in "key=values" format.
  * `-P, --gopath string`: Set a custom `GOPATH` for the script execution.
  * `-S, --sandboxed`: Run with sandboxed standard library symbols (e.g., restricted `os/exec`).
  * `-T, --tags strings`: Define build constraints for the scripts.

### Script Loading Mechanism

GoApp utilizes a **GOPATH Hook Hijacking** strategy to handle local imports:

  - When you pass a `Project_Dir` or `GO_FILE`, GoApp automatically calculates the base directory and maps it to a virtual `GOPATH/src`.
  - **The Magic**: It injects a `MAGICDIR` into the runtime. For example, running `goapp run example/myapp` triggers a hook that redirects imports from `MAGICDIR/src/myapp` to the actual physical path of `example/myapp`.
  - This allows scripts to `import` their own project folders directly without needing to be located inside a real `GOPATH` directory.

-----

## Installation & Custom Build

Pre-compiled binaries are available for quick use, but they only include a subset of supported libraries. For full compatibility with your specific stack, building from source is recommended.

### 1\. Prerequisites

Install the Yaegi symbol extraction tool:

```bash
go install github.com/traefik/yaegi/cmd/yaegi@latest
```

### 2\. Clone the Repository

```bash
git clone https://github.com/zuiwuchang/goapp.git
cd goapp
```

### 3\. Configure Dependencies

Edit `script/conf.sh` and add the import paths of the third-party libraries you need into the `GOLIB` array:

```bash
GOLIB=(
    github.com/zuiwuchang/gosdk
    github.com/spf13/cobra
    github.com/fsnotify/fsnotify
    github.com/jroimartin/gocui
)
```

### 4\. Generate Symbols and Build

Run the build scripts to generate the necessary glue code and compile the binary:

```bash
# Generate symbol information and embedded code
./build.sh symbols -g

# Compile the GoApp binary
./build.sh go
```