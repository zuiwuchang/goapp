# GoApp

[中文說明](./zh.md) | **English Version**

`GoApp` is a powerful wrapper for [Yaegi](https://github.com/traefik/yaegi), a Go engine that enables you to interpret and execute Go code as scripts. It addresses common limitations in the standard Yaegi, such as local project imports and the effort required to extend third-party library support.

### Key Features

  * **Expanded Built-in Libraries**: Pre-packaged with a wider range of popular third-party libraries. See [third-party.md](https://www.google.com/search?q=./third-party.md) for the full list.
  * **Easy Extension**: Effortlessly embed additional Go symbols and packages into the interpreter through simple configuration.
  * **Smart Path Hijacking**: Features a custom `GOPATH` hook mechanism that allows you to `import` local projects without complex environment setup.

-----

## Usage

### Basic Command

```bash
# Execute a Go package or script file
$ goapp run [flags] [Project_Dir | GO_FILE] -- [args]
```

### Flags

| Flag | Description |
| :--- | :--- |
| `-E, --env` | Set environment variables in "key=values" format |
| `-P, --gopath` | Set a custom GOPATH for the script execution |
| `-S, --sandboxed` | Enable sandbox mode (restricts sensitive stdlib symbols like `os/exec`) |
| `-T, --tags` | Set build constraints (build tags) for the scripts |

> **Tip**: Use the `--` separator to distinguish between `goapp` flags and the arguments you wish to pass to your script.

-----

## Script Loading & IDE Support

### 1\. Magic Import (Hook Hijacking)

GoApp simplifies the import logic for local development:

  * When you run `goapp run example/myapp`, the program automatically maps the parent directory to a virtual `GOPATH/src`.
  * **How it works**: It uses a runtime hook to inject a `MAGICDIR`. When the script attempts to load `MAGICDIR/src/myapp`, GoApp transparently redirects it to the actual physical path of `example/myapp`.
  * **The Benefit**: You can use the "Project Folder Name" as the package name for imports directly, eliminating the need to move code into a standard Go workspace.

### 2\. VSCode IntelliSense Trick

To get a full-featured IDE experience with code completion and navigation:

1.  Copy the `go.mod` and `go.sum` files from this repository into your **script project directory**.
2.  **Crucial Step**: Edit the copied `go.mod` and change the `module` name to match your **script folder name**.
3.  Ensure you have the VSCode Go extension and Go environment installed. The IDE will now recognize all built-in third-party libraries and provide full IntelliSense.

-----

## Installation & Custom Build

While pre-compiled binaries are available, they only include a subset of supported libraries. For full library support, we recommend building from source.

### 1\. Prerequisites

Install the latest Yaegi symbol extraction tool:

```bash
go install github.com/traefik/yaegi/cmd/yaegi@latest
```

### 2\. Clone the Repository

```bash
git clone https://github.com/zuiwuchang/goapp.git
cd goapp
```

### 3\. Configure Dependencies

Edit the `script/conf.sh` file and add your required library paths to the `GOLIB` array:

```bash
GOLIB=(
    github.com/zuiwuchang/gosdk
    github.com/spf13/cobra
    github.com/fsnotify/fsnotify
    github.com/jroimartin/gocui
)
```

### 4\. Generate Symbols and Build

Run the build scripts to automate symbol generation, code embedding, and compilation:

```bash
# Generate symbol information and glue code
./build.sh symbols -g

# Perform final compilation
./build.sh go
```

-----

## License

[MIT License]