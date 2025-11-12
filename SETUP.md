# SnapAsk Setup Guide

Complete setup instructions for building and running SnapAsk on all supported platforms.

## 📋 Prerequisites

### All Platforms
- **Go 1.22+**: [Download from golang.org](https://golang.org/dl/)
- **Node.js 18+**: [Download from nodejs.org](https://nodejs.org/)
- **Wails CLI**: Install with `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Platform-Specific Requirements

#### Windows
- **Go 1.22+**
- **Node.js 18+**
- **WebView2 Runtime**: Usually pre-installed on Windows 10/11
- **NSIS** (optional, for installer): [Download from nsis.sourceforge.io](https://nsis.sourceforge.io/)

#### macOS
- **Go 1.22+**
- **Node.js 18+**
- **Xcode Command Line Tools**: `xcode-select --install`

---

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/SnapAsk.git
cd SnapAsk
```

### 2. Install Frontend Dependencies
```bash
cd frontend
npm install
cd ..
```

### 3. Initialize Go Modules
```bash
go mod download
```

### 4. Run in Development Mode
```bash
wails dev
```

The application will launch with hot-reload enabled for frontend changes.

---

## 🔨 Building for Production

### Windows (.exe)
```bash
wails build -platform windows/amd64
```

Output: `build/bin/Korner.exe`

**Create Installer:**
```bash
wails build -platform windows/amd64 -nsis
```

### macOS (.app)
```bash
# Intel
wails build -platform darwin/amd64

# Apple Silicon
wails build -platform darwin/arm64

# Universal Binary (both architectures)
wails build -platform darwin/universal
```

Output: `build/bin/Korner.app`

**Sign and Notarize (for distribution):**
```bash
# Sign
codesign --deep --force --verify --verbose --sign "Developer ID Application: Your Name" build/bin/Korner.app

# Notarize (requires Apple Developer account)
xcrun notarytool submit build/bin/Korner.app.zip --keychain-profile "AC_PASSWORD" --wait
```

---

## ⚙️ Configuration

### AMD LLM Endpoint Configuration

Create a `.env` file in the project root:

```env
AMD_LLM_ENDPOINT=https://your-amd-endpoint.com/v1/chat/completions
AMD_API_KEY=your-api-key-here
MODEL_NAME=gpt-oss-120b
```

**Load in Go backend (`app.go`):**
```go
import (
    "os"
    "github.com/joho/godotenv"
)

func init() {
    godotenv.Load()
}
```

### Global Hotkey Setup

#### Windows/macOS
Hotkeys are automatically registered:
- **Windows**: `Ctrl+Alt+Q`
- **macOS**: `Cmd+Option+Q`

## 🧪 Development Workflow

### Frontend Development
```bash
cd frontend
npm run dev
```

Visit `http://localhost:34115` to see live changes.

### Hot Reload with Wails
```bash
wails dev
```

Changes to Vue files reload automatically. Go changes require restart.

### Build Frontend Only
```bash
cd frontend
npm run build
```

### Testing

#### Run Frontend Tests
```bash
cd frontend
npm test  # (tests to be added)
```

#### Run Go Tests
```bash
go test ./...
```

---

## 📁 Project Structure

```
SnapAsk/
├── main.go                 # Wails entry point
├── app.go                  # App struct and backend methods
├── go.mod                  # Go dependencies
├── wails.json             # Wails configuration
├── .env                   # Environment variables (not in git)
├── frontend/              # Vue 3 application
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── components/
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
└── build/                 # Build assets and outputs
    ├── bin/              # Compiled binaries
    ├── appicon.png       # Linux icon
    ├── appicon.ico       # Windows icon
    └── appicon.icns      # macOS icon
```

---

## 🐛 Troubleshooting

### Windows

**Issue**: WebView2 not found
```
Solution: Install WebView2 Runtime from Microsoft
https://developer.microsoft.com/microsoft-edge/webview2/
```

**Issue**: Build fails with "gcc not found"
```
Solution: Install TDM-GCC or MinGW-w64
```

### macOS

**Issue**: "xcrun: error: invalid active developer path"
```bash
xcode-select --install
```

**Issue**: App won't open (Gatekeeper)
```bash
xattr -cr build/bin/SnapAsk.app
```
---

## 🚢 Distribution

### Windows
- Distribute `.exe` directly or create NSIS installer
- Consider code signing for Windows SmartScreen

### macOS
- Distribute `.app` in `.dmg` disk image
- **Must** be signed and notarized for macOS 10.15+
---

## 📝 Next Steps

1. ✅ Setup development environment
2. ✅ Run `wails dev` to test
3. 🔄 Implement screenshot capture backend
4. 🔄 Connect to AMD LLM API
5. 🔄 Test on all platforms
6. 📦 Build production binaries
7. 🚀 Create demo video
8. 📤 Submit to AMD Hackathon

---

## 🆘 Getting Help

- **Wails Documentation**: https://wails.io/docs/introduction
- **Vue 3 Documentation**: https://vuejs.org/guide/introduction.html
- **Tailwind CSS**: https://tailwindcss.com/docs
- **Project Issues**: https://github.com/yourusername/SnapAsk/issues

---

**Good luck with the AMD Hackathon! 🚀**