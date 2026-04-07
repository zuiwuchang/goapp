# GoApp

[English Version](./README.md) | 中文說明

`GoApp` 是一個基於 [Yaegi](https://github.com/traefik/yaegi) 的 Go 腳本解釋器外殼。它允許你直接解析並執行 Go 代碼，並解決了 Yaegi 在處理本地專案目錄與第三方庫擴展時的痛點。

### 核心特性

  * **內置豐富庫支援**：相比原生 Yaegi，預集成了更多常用的第三方庫。
  * **動態擴展**：透過簡單的配置即可輕易添加額外符號（Symbols）到解釋器中。
  * **智能路徑劫持**：創新的 `GOPATH` Hook 機制，讓你無需複雜配置即可 `import` 本地專案。

-----

## 使用方法

### 基礎命令

```bash
# 執行一個 Go 套件或腳本
$ goapp run [flags] [Project_Dir | GO_FILE] -- [args]
```

### 參數說明

| Flag | 說明 |
| :--- | :--- |
| `-E, --env` | 設定環境變量，格式為 "key=values" |
| `-P, --gopath` | 為腳本設定自定義的 GOPATH |
| `-S, --sandboxed` | 開啟沙盒模式，限制 os/exec 等標準庫的執行 |
| `-T, --tags` | 設定編譯約束（Build Constraints） |

> **提示**：使用 `--` 符號來區隔解釋器參數與傳遞給腳本的參數。

### 腳本加載機制 (Magic Import)

GoApp 簡化了本地開發的導入邏輯：

1.  **自動路徑映射**：當執行 `goapp run example/myapp` 時，程序會自動將該目錄的上層路徑視為 `GOPATH/src` 的一部分。
2.  **Hook 劫持**：採用 Runtime Hook 技術。例如執行時會向虛擬 `GOPATH` 注入一個 `MAGICDIR`，當腳本嘗試加載 `MAGICDIR/src/myapp` 時，會重定向至 `example/myapp` 實際物理路徑。
3.  **優勢**：你可以直接使用「項目文件夾名」作為包名進行導入，無需將代碼手動移入標準的 `GOPATH` 目錄。

-----

## 安裝與自定義編譯

雖然我們提供預編譯版本，但預編譯版僅包含部分常用庫。為了獲得最完整的第三方庫支援，建議自行編譯。

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

編輯 `script/conf.sh` 文件，在 `GOLIB` 數組中加入你需要嵌入的庫路徑：

```bash
GOLIB=(
    github.com/zuiwuchang/gosdk
    github.com/spf13/cobra
    github.com/fsnotify/fsnotify
    github.com/jroimartin/gocui
)
```

### 4\. 生成符號與編譯

執行以下腳本自動生成符號信息、嵌入代碼並編譯：

```bash
# 生成符號信息與嵌入代碼
./build.sh symbols -g

# 執行最終編譯
./build.sh go
```

