# GoApp

[English Version](./README.md) | **中文說明**

`GoApp` 是一個基於 [Yaegi](https://github.com/traefik/yaegi) 的 Go 腳本解釋器外殼。它允許你直接解析並執行 Go 代碼，並解決了原生 Yaegi 在處理本地專案目錄導入與第三方庫擴展時的痛點。

### 核心特性

  * **內置豐富庫支援**：相比原生 Yaegi，預集成了更多常用的第三方庫（如 Cobra, Gin, Gocui 等）。
  * **動態擴展**：透過簡單的配置文件與內置的 `extract` 腳本，即可輕易添加額外符號（Symbols）。
  * **智能路徑劫持**：創新的 `GOPATH` Hook 機制，讓你無需複雜配置即可在腳本中 `import` 本地專案。
  * **IDE 友善設計**：透過特殊的入口機制，完美解決 VSCode 在同一目錄下處理多個腳本時的語法報錯問題。

-----

## 使用方法

### 基礎命令

```bash
# 執行一個 Go 套件或腳本文件
$ goapp run [flags] -- [Project_Dir | GO_FILE] [args]
```

### 參數說明

| Flag | 說明 | 預設值 |
| :--- | :--- | :--- |
| `-E, --env` | 設定環境變數，格式為 "key=values" | 無 |
| `-F, --func` | 腳本載入後，**明確指定呼叫**的函式名稱（可多次指定） | 無 |
| `-P, --gopath` | 為腳本設定自定義的 GOPATH | `$GOPATH` |
| `-S, --sandboxed` | 開啟沙盒模式，限制 os/exec 等標準庫的執行 | `false` |
| `-T, --tags` | 設定編譯約束（Build Constraints） | 無 |

> **提示**：使用 `--` 符號來區隔解釋器本身的參數與傳遞給腳本的參數。

-----

## 執行入口機制與 IDE 優化

為了讓開發過程更順暢，`GoApp` 實作了特殊的入口偵測邏輯：

### 1\. RunMain 自動執行 (欺騙 IDE 技巧)

在 VSCode 中，若同一個資料夾下的多個 `.go` 檔案都定義了 `func main()`，IDE 會提示重複定義錯誤。
**解決方案**：在 `GoApp` 中，你可以將入口函式命名為 **`RunMain`** 或任何以 **`RunMain`** 開頭的名稱（例如 `RunMainServer`）。

  * **單檔案模式**：執行 `goapp run script.go` 時，解釋器會自動尋找並執行該檔案內所有 `RunMain` 開頭的函式。
  * **優勢**：IDE 不會將其視為標準 `main` 函式，因此你可以在同一個 package 下定義多個不同的入口檔案而不會報錯。

### 2\. 指定呼叫模式 (`-F` 參數)

如果你需要更精確地控制執行流程，可以使用 `-F` 參數呼叫特定函式：

```bash
# 執行載入後，依序呼叫指定的函式
$ goapp run -F InitConfig -F StartServer -- script.go
```

### 3\. 目錄/套件模式

執行 `goapp run ./my-pkg`（指定目錄）時，解釋器會載入整個套件。此時**不會**觸發自動 `RunMain`，請務必搭配 `-F` 參數指定入口。

-----

## 腳本加載機制：Magic Import

`GoApp` 採用了 **Runtime Hook** 技術來簡化本地導入：

  * **原理**：當你執行 `goapp run example/myapp`，程序會向虛擬 `GOPATH` 注入一個隱藏的 `MAGICDIR`。當腳本嘗試加載 `myapp` 時，會自動重定向至實際的實體路徑。
  * **優勢**：你可以在腳本中直接 `import "myapp/utils"`，無需手動配置複雜的環境變數，且能保持與標準 Go 專案一致的邏輯結構。

-----

## VSCode 代碼提示 SOP

若想獲得完美的補全體驗，請遵循以下步驟：

1.  將本項目的 `go.mod` 與 `go.sum` 拷貝到你的 **腳本項目目錄** 下。
2.  **修改 Mod 名稱**：將 `go.mod` 中的 `module` 名稱修改為你的 **腳本資料夾名稱**。
3.  **避免 `main` 衝突**：使用 `RunMain...` 作為入口函式名稱而非 `main()`。
4.  確保安裝了 VSCode Go 插件，此時 IDE 將能識別所有內置庫（包含 `gosdk`）並提供完整的代碼提示。

-----

## 自定義編譯

如果需要增加自定義的第三方庫，請參考以下流程：

1.  編輯 `script/conf.sh`，在 `GOLIB` 陣列中加入套件路徑。
2.  執行 `./build.sh symbols -g`：自動 `go get` 並導出 Yaegi 所需的符號。
3.  執行 `./build.sh go`：編譯最終的 `goapp` 二進位檔案。

-----

## 授權

[MIT License]