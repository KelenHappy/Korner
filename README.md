# 軟體需求規格書（SRS）
## 專案名稱：**Korner — 跨平台全域截圖提問助理**
> 一鍵截圖 / AI 提問 → 即時回應
> 支援 Windows / macOS / Linux（GNOME & KDE via Flatpak）

### 1. 項目簡介
Korner 是一款輕量級 Agentic AI 桌面助理，讓使用者無需切換瀏覽器或多個應用，直接透過全域快捷鍵（預設 `Ctrl+Alt+Q` / `Cmd+Option+Q`）截圖並向 **AMD Instinct MI300X GPU 上的 GPT OSS 120B 模型** 提問（例如：「解釋這段程式碼」、「找出圖表錯誤」）或點開懸浮的icon進行提問或截圖。
本專案完整支援三大平台：
- **Windows 10/11**
- **macOS 12+ (Apple Silicon & Intel)**
- **Linux — 僅限 GNOME 42+ 與 KDE Plasma 5.24+，以 Flatpak 發佈**

所有請求透過 OpenAI API 相容格式發送至 AMD 提供的 LLM 端點，展現 **「感知 → 思考 → 行動」** 的 Agentic 能力。

---

### 2. 核心功能需求

| 功能 | 說明 |
|------|------|
| **全域快捷鍵觸發** | 跨平台註冊系統級快捷鍵，應用在背景時仍可喚起截圖 |
| **選區截圖** | 啟動後顯示全螢幕半透明遮罩，拖曳選取區域 |
| **AI 提問（含圖片）** | 截圖轉 Base64 + 使用者文字，發送至 GPT OSS 120B Vision API |
| **浮動回應視窗** | 顯示 LLM 回答，支援複製、固定、關閉 |
| **多平台一致體驗** | UI/UX 在三大平台保持一致，行為符合各系統慣例 |

---

### 3. 各平台技術實現細節

#### 3.1 Windows
- **截圖**：使用 `github.com/kbinani/screenshot` 或 Win32 `BitBlt` API 實作選區截圖
- **全域熱鍵**：透過 Win32 `RegisterHotKey()` + Wails Go 後端註冊
- **打包**：`wails build` 產出 `.exe`，內建系統匣圖示（tray icon）作為備援觸發方式

#### 3.2 macOS
- **截圖**：呼叫原生 `screencapture -i -x -t png /tmp/Korner.png`（無提示音、選區模式）
- **全域熱鍵**：使用 Carbon Event Manager（透過 CGO）註冊 `Cmd+Option+Q`
  - 需在 `Info.plist` 聲明 `Accessibility` 權限（首次啟動引導使用者授權）
- **打包**：`wails build` 產出 `.app`，支援 Apple Silicon 與 Intel 雙架構（Universal Binary）

#### 3.3 Linux（GNOME / KDE + Flatpak）
- **截圖**：透過 `org.freedesktop.portal.Screenshot`（D-Bus + xdg-desktop-portal）
  - 自動適配 Wayland/X11，由桌面環境提供原生 UI
- **全域熱鍵**：
  - 因 Flatpak 沙盒限制，**不自動註冊熱鍵**
  - **首次啟動時顯示教學**：引導使用者至 GNOME Settings 或 KDE System Settings 手動新增快捷鍵，命令為：
    ```bash
    flatpak run com.korner.Korner --screenshot
    ```
  - 同時提供 **system tray icon** 點擊觸發（相容 GNOME Shell Extensions / KDE）
- **打包**：
  - 使用 `flatpak-builder` 打包
  - runtime: `org.freedesktop.Platform//23.08`
  - 宣告權限：
    ```ini
    --socket=wayland
    --socket=fallback-x11
    --device=dri
    --talk-name=org.freedesktop.portal.Desktop
    ```

---

### 4. 架構與開發技術

| 層級 | 技術 |
|------|------|
| 框架 | **Wails v2**（Go 1.22 + Vue 3） |
| 前端 | Vue 3 + Tailwind CSS（輕量、響應式） |
| 後端 | Go（處理截圖、熱鍵、LLM 通訊） |
| LLM 連接 | 發送 OpenAI API 格式至 AMD 提供的 GPT OSS 120B Vision 端點 |
| 部署 |
  - Windows: `.exe`
  - macOS: `.app` (Universal)
  - Linux: Flatpak（支援 Flathub 提交） |

---

### 5. 非功能需求

| 項目 | 說明 |
|------|------|
| **隱私安全** | 截圖僅存於記憶體；不寫入硬碟；未經同意不儲存使用者資料 |
| **開源與文件** | GitHub 開源 + **英文 README**，包含：
  - 各平台安裝與快捷鍵設定指引
  - 如何連接 AMD LLM 端點
  - Flatpak 權限與建置步驟 |
| **體積與效能** | 安裝包 ≤ 60MB；截圖到回應 ≤ 4 秒（取決於 AMD 端點） |
| **符合競賽規範** | 使用 AMD Instinct MI300X + GPT OSS 120B，展現 Agentic AI 能力 |

---

### 6. 對應 AMD Hackathon 評分標準

| 評分項目 | 本專案實現方式 |
|---------|----------------|
| **技術創新與完整度 (30%)** | 跨平台截圖 + Vision LLM 提問，完整 Agentic 循環 |
| **實用性與應用潛力 (30%)** | 解決「多應用切換提問」痛點，可推廣為通用桌面助理 |
| **創意與問題定義 (25%)** | 聚焦「圖像+文字」混合提問場景（如程式碼、圖表、錯誤畫面） |
| **展示與呈現 (15%)** | 提供三平台 Demo 影片 + Live 操作展示 |

---

### 7. 交付內容（12/12 前上傳）

- **GitHub 儲存庫**（開源、含英文 README）
- **可運行 Demo**：
  - Windows `.exe`
  - macOS `.app`
  - Linux `.flatpak`（或 Flathub 測試安裝指令）
- **Demo 影片（≤3 分鐘）**：展示三平台截圖提問流程
- **文字說明（PDF/Markdown）**：簡述架構、創新點、如何重現

---

> 💡 此設計平衡了 **技術可行性** 與 **競賽創新性**，特別在 Linux 上尊重 Wayland 安全模型，同時確保 Windows/macOS 提供無縫體驗，符合「真實問題解決」精神。

如需我協助產生：
- Wails 專案結構
- 各平台截圖程式碼範例
- Flatpak manifest 檔案
- OpenAI Vision API 呼叫格式（Go）
README.md
6 KB
