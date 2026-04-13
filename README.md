# GoApp

[中文說明](./README.zh.md) | **English Version**

`GoApp` is a Go script interpreter shell based on [Yaegi](https://github.com/traefik/yaegi). It allows you to parse and execute Go code directly, resolving common pain points of native Yaegi such as local project directory imports and third-party library extensions.

### Key Features

  * **Rich Built-in Libraries**: Pre-integrated with popular third-party libraries (e.g., Cobra, Gin, Gocui) compared to native Yaegi.
  * **Dynamic Extension**: Easily add extra symbols to the interpreter via a simple configuration file and the built-in `extract` script.
  * **Smart Path Hijacking**: An innovative `GOPATH` Hook mechanism that allows you to `import` local projects in your scripts without complex configuration.
  * **IDE-Friendly Design**: Resolves VSCode syntax errors when handling multiple scripts in the same directory through a specialized entry mechanism.

-----

## Usage

### Basic Command

```bash
# Execute a Go package or script file
$ goapp run [flags] -- [Project_Dir | GO_FILE] [args]
```

### Flags

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-E, --env` | Set environment variables in "key=values" format | None |
| `-F, --func` | **Explicitly call** specified function(s) after loading (can be used multiple times) | None |
| `-P, --gopath` | Set a custom GOPATH for the scripts | `$GOPATH` |
| `-S, --sandboxed` | Enable sandbox mode, restricting standard libraries like `os/exec` | `false` |
| `-T, --tags` | Set build constraints (Build Tags) | None |

> **Tip**: Use the `--` symbol to separate the interpreter flags from the arguments passed to your script.

-----

## Entry Mechanism & IDE Optimization

To streamline the development process, `GoApp` implements a unique entry detection logic:

### 1\. RunMain Auto-Execution (The IDE Trick)

In VSCode, if multiple `.go` files in the same folder define `func main()`, the IDE will report a "duplicate declaration" error.
**Solution**: In `GoApp`, you can name your entry functions **`RunMain`** or anything starting with **`RunMain`** (e.g., `RunMainServer`).

  * **Single File Mode**: When running `goapp run script.go`, the interpreter automatically finds and executes all functions starting with `RunMain` within that file.
  * **Advantage**: The IDE does not treat these as standard `main` functions, allowing you to define multiple entry files within the same package without errors.

### 2\. Explicit Call Mode (`-F` flag)

If you need precise control over the execution flow, use the `-F` flag to call specific functions:

```bash
# After loading, call the specified functions in order
$ goapp run -F InitConfig -F StartServer -- script.go
```

### 3\. Directory/Package Mode

When running `goapp run ./my-pkg` (specifying a directory), the interpreter loads the entire package. In this mode, **automatic `RunMain` detection is disabled**; you must use the `-F` flag to specify the entry point.

-----

## Loading Mechanism: Magic Import

`GoApp` utilizes **Runtime Hook** technology to simplify local imports:

  * **Mechanism**: When you execute `goapp run example/myapp`, the program injects a hidden `MAGICDIR` into the virtual `GOPATH`. When the script attempts to load `myapp`, it is automatically redirected to the actual physical path.
  * **Advantage**: You can `import "myapp/utils"` directly in your script, maintaining a logical structure consistent with standard Go projects without manual environment variable hacking.

-----

## VSCode Code Completion SOP

To achieve a perfect autocompletion experience, follow these steps:

1.  Copy the `go.mod` and `go.sum` files from this project to your **script project directory**.
2.  **Modify Mod Name**: Edit the copied `go.mod` and change the `module` name to match your **script folder name**.
3.  **Avoid `main` Conflict**: Use `RunMain...` as the entry function name instead of `main()`.
4.  Once the VSCode Go extension is active, it will recognize all built-in libraries (including `gosdk`) and provide full IntelliSense.

-----

## Custom Compilation

If you need to add custom third-party libraries, follow this workflow:

1.  Edit `script/conf.sh` and add the package paths to the `GOLIB` array.
2.  Run `./build.sh symbols -g`: This automatically performs `go get` and extracts the symbols required by Yaegi.
3.  Run `./build.sh go`: Compiles the final `goapp` binary.

-----

## License

[MIT License]