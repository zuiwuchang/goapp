# GoApp

[English Version](./README.md) | **中文說明**

`GoApp` 是一個基於 [Yaegi](https://github.com/traefik/yaegi) 的 Go 腳本解釋器外殼。它允許你直接解析並執行 Go 代碼，並解決了原生 Yaegi 在處理本地專案目錄導入與第三方庫擴展時的痛點。

### 核心特性

  * **內置豐富庫支援**：相比原生 Yaegi，預集成了更多常用的第三方庫。具體清單請參閱 [third-party.md](https://www.google.com/search?q=./third-party.md)。
  * **動態擴展**：透過簡單的配置文件即可輕易添加額外符號（Symbols）到解釋器中。
  * **智能路徑劫持**：創新的 `GOPATH` Hook 機制，讓你無需複雜配置即可在腳本中 `import` 本地專案。

-----

## 使用方法

### 基礎命令

```bash
# 執行一個 Go 套件或腳本文件
$ goapp run [flags] [Project_Dir | GO_FILE] -- [args]
```

### 參數說明

| Flag | 說明 |
| :--- | :--- |
| `-E, --env` | 設定環境變量，格式為 "key=values" |
| `-P, --gopath` | 為腳本設定自定義的 GOPATH |
| `-S, --sandboxed` | 開啟沙盒模式，限制 os/exec 等標準庫的執行 |
| `-T, --tags` | 設定編譯約束（Build Constraints） |

> **提示**：使用 `--` 符號來區隔解釋器本身的參數與傳遞給腳本的參數。

-----

## 腳本加載機制與 IDE 支持

### 1\. Magic Import (Hook 劫持)

GoApp 簡化了本地開發的導入邏輯：

  * 當執行 `goapp run example/myapp` 後，程序會自動將 `example` 目錄添加到虛擬 `GOPATH/src` 下。
  * **原理**：採用 Runtime Hook 技術。例如執行時會向虛擬 `GOPATH` 注入一個 `MAGICDIR`，當腳本嘗試加載 `MAGICDIR/src/myapp` 時，會自動重定向至 `example/myapp` 的實際內容。
  * **優勢**：你可以直接使用「項目文件夾名」作為包名進行導入，無需手動配置複雜的路徑。

### 2\. VSCode 代碼提示技巧

若想在編寫腳本時獲得完美的補全體驗：

1.  將本項目的 `go.mod` 與 `go.sum` 拷貝到你的 **腳本項目目錄** 下。
2.  **重要步驟**：編輯拷貝後的 `go.mod`，將其中的 `module` 名稱修改為你的 **腳本文件夾名稱**（否則 VSCode 會提示找不到專案）。
3.  確保安裝了 VSCode Go 插件及 Go 環境，此時 IDE 將能識別解釋器內置的所有第三方庫並提供完整的代碼提示。

-----

## 安裝與自定義編譯

你可以直接下載預編譯版本，但預編譯版僅打包了一部分常用庫。如果需要完整的第三方庫支持，建議自行編譯。

### 1\. 準備環境

確保已安裝最新版的 Yaegi 符號工具：

```bash
go install github.com/traefik/yaegi/cmd/yaegi@latest
```

### 2\. 獲取源碼

```bash
git clone https://github.com/zuiwuchang/goapp.git
cd goapp
```

### 3\. 配置所需的第三方庫

編輯 `script/conf.sh` 文件，在 `GOLIB` 數組中加入你需要的庫：

```bash
GOLIB=(
    github.com/zuiwuchang/gosdk
    github.com/spf13/cobra
    github.com/fsnotify/fsnotify
    github.com/jroimartin/gocui
)
```

### 4\. 生成符號與編譯

執行以下腳本自動生成符號信息、嵌入代碼並完成編譯：

```bash
# 生成符號信息與嵌入代碼
./build.sh symbols -g

# 執行最終編譯
./build.sh go
```

-----

## 授權

[MIT License]