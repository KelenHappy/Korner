# Korner Frontend - Vue 3 Application

This is the frontend application for SnapAsk, built with Vue 3, Vite, and Tailwind CSS.

## 🏗️ Architecture

### Technology Stack
- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite 5
- **Styling**: Tailwind CSS 3
- **Language**: JavaScript (ES6+)

### Project Structure

```
frontend/
├── src/
│   ├── App.vue                 # Main application component
│   ├── main.js                 # Application entry point
│   ├── style.css               # Global styles (Tailwind imports)
│   └── components/
│       ├── ScreenshotOverlay.vue    # Fullscreen screenshot selection overlay
│       ├── QueryWindow.vue          # Query input window after screenshot
│       ├── ResponseWindow.vue       # Floating AI response window
│       └── FirstRunGuide.vue        # Linux first-run hotkey setup guide
├── public/                     # Static assets
├── index.html                  # HTML entry point
├── vite.config.js             # Vite configuration
├── tailwind.config.js         # Tailwind CSS configuration
├── postcss.config.js          # PostCSS configuration
└── package.json               # Dependencies and scripts
```

## 📦 Components Overview

### App.vue
The main application component that orchestrates all functionality:
- **State Management**: Manages screenshot overlay, conversation history, current query
- **Platform Detection**: Detects OS for appropriate hotkey display
- **First Run**: Shows Linux setup guide on first launch
- **Event Handling**: Coordinates screenshot capture and AI query workflow

**Key Features:**
- Welcome screen with usage examples
- Conversation history display
- Integration with all child components
- Responsive layout with Tailwind CSS

### ScreenshotOverlay.vue
Fullscreen transparent overlay for area selection:
- Click and drag to select screen region
- Shows selection dimensions in real-time
- ESC to cancel
- Minimum 10x10px selection required
- Generates placeholder screenshot (in production, calls Go backend)

### QueryWindow.vue
Query input interface after screenshot capture:
- Displays screenshot preview
- Text area for user question
- Quick prompt buttons for common queries
- Ctrl/Cmd+Enter to submit
- Character counter (1000 char limit)

**Quick Prompts:**
- "Explain this code"
- "What does this mean?"
- "Find the error"
- "Summarize this"
- "Translate this"
- "Improve this code"

### ResponseWindow.vue
Floating window showing AI responses:
- Fixed position (bottom-right corner)
- Loading animation during API call
- Copy, pin, and close actions
- Scrollable content area
- Beautiful gradient styling

### FirstRunGuide.vue
Modal dialog for Linux users on first launch:
- Instructions for manual hotkey setup
- Command to add to system shortcuts
- Stores "seen" state in localStorage
- Only shown once per installation

## 🎨 Styling

### Tailwind CSS Configuration
- Gradient backgrounds (slate-50 to slate-100)
- Primary color: Blue (500-600)
- AI response accent: Purple-Pink gradient
- Responsive breakpoints (mobile-first)
- Custom animations (bounce for loading states)

### Design System
- **Spacing**: Consistent padding/margins (p-4, p-6, space-x-3)
- **Borders**: Rounded corners (rounded-lg, rounded-full)
- **Shadows**: Elevation system (shadow-sm, shadow-lg, shadow-2xl)
- **Typography**: Slate color palette for text hierarchy

## 🔗 Backend Integration

### Wails Bindings
The frontend communicates with the Go backend via Wails:

```javascript
// Example calls (to be implemented):
await window.go.main.App.GetPlatform()
await window.go.main.App.CaptureScreenshot(x, y, width, height)
await window.go.main.App.QueryLLM(query, screenshotBase64)
```

### Current Implementation
- Simulated screenshot capture (gradient placeholder)
- Mock AI responses with 2-second delay
- Platform detection via navigator.userAgent

## 🚀 Development

### Install Dependencies
```bash
cd frontend
npm install
```

### Run Development Server
```bash
npm run dev
```
Server runs on `http://localhost:34115`

### Build for Production
```bash
npm run build
```
Output goes to `frontend/dist/`

## 📝 State Management

### Reactive State (Vue Composition API)
```javascript
showScreenshotOverlay    // Boolean - show/hide screenshot UI
showFirstRunGuide        // Boolean - show Linux setup guide
currentQuery             // Object - pending query with screenshot
conversationHistory      // Array - all Q&A pairs
showResponseWindow       // Boolean - floating response window
latestResponse           // String - most recent AI response
isLoadingResponse        // Boolean - loading state
platform                 // String - detected OS
```

## 🔑 Keyboard Shortcuts

- **Windows/Linux**: `Ctrl+Alt+Q` - Trigger screenshot
- **macOS**: `Cmd+Option+Q` - Trigger screenshot
- **In Query Window**: `Ctrl/Cmd+Enter` - Submit query
- **In Screenshot Overlay**: `ESC` - Cancel selection

## 🎯 User Flow

1. **Launch App** → Welcome screen displayed
2. **Trigger Screenshot** → Click button or press hotkey
3. **Select Area** → Click and drag on overlay
4. **Enter Query** → Type question about screenshot
5. **View Response** → AI response in floating window
6. **History** → All conversations saved in main view

## 🔧 Configuration

### Vite Config
- Port: 34115 (configured to work with Wails)
- Asset handling optimized for production builds
- Vue plugin enabled

### Tailwind Config
- Content paths: `./index.html`, `./src/**/*.{vue,js}`
- Default theme with extensions
- Utilities for common patterns

## 📱 Responsive Design

- Mobile-first approach
- Grid layouts adapt to screen size
- Touch-friendly targets (min 44x44px)
- Overflow scrolling in constrained areas

## 🐛 Error Handling

- Minimum screenshot size validation
- Query text length limits
- Graceful API failure handling
- User-friendly error messages

## 🔮 Future Enhancements

- [ ] Real screenshot capture via backend
- [ ] Stream AI responses (token-by-token)
- [ ] Markdown rendering for code blocks
- [ ] Image zoom/fullscreen viewer
- [ ] Export conversation history
- [ ] Multiple language support (i18n)
- [ ] Dark mode toggle
- [ ] Custom hotkey configuration
- [ ] Screenshot annotation tools

## 📄 License

Part of the SnapAsk project - AMD Hackathon 2024

---

For backend documentation, see the main project README.md
