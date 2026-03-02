# Human Desktop App — Build Prompt for Claude Code

> Copy this entire file into Claude Code as the starting prompt. It contains everything needed to build the desktop application from scratch.

---

## Project Overview

Build **Human Studio** — a cross-platform desktop application for the Human programming language. Human compiles structured English into production-ready full-stack applications. This desktop app is the primary IDE for non-technical users to write .human files, prompt with LLMs, build, and deploy — all without touching a terminal.

**Target users:** Non-coders, founders, product managers, designers — people who can describe what they want but can't write React/Node/SQL.

**Ship targets:** Windows (.exe / .msi), macOS (.dmg), Ubuntu (.deb / .AppImage)

**Think:** Cursor meets Lovable meets VS Code — but for the Human language.

---

## Technology Stack

Use **Electron** with the following stack:

| Layer | Technology | Why |
|-------|-----------|-----|
| Shell | Electron 33+ | Cross-platform desktop, native menus, file system access |
| Renderer | React 19 + TypeScript | Component model, ecosystem, team familiarity |
| Styling | Tailwind CSS 4 + CSS variables | Design token system, utility-first, dark/light mode |
| State | Zustand | Lightweight, no boilerplate, good for multi-window |
| Editor | CodeMirror 6 | Extensible, custom language modes, lightweight |
| Terminal | xterm.js | Embedded terminal for build output |
| File tree | Custom (virtual scroll) | Performance with large projects |
| IPC | Electron IPC + contextBridge | Secure main↔renderer communication |
| Build | electron-builder | Windows/Mac/Linux packaging |
| Auto-update | electron-updater | Squirrel (Win), Sparkle (Mac), AppImage (Linux) |

---

## Brand Identity (MANDATORY — Do Not Deviate)

### The Brand

Human's brand is **monochrome with a single accent color**. The design is predominantly black/white/gray with `#E85D3A` (warm coral) as the sole pop of color. Color is used sparingly and intentionally.

**Brand personality:** The smartest, kindest person in the room. Listens carefully, never condescends, makes complex things feel simple. Not flashy. Quietly, confidently excellent.

### Logo

The logo is the word **human** in lowercase Nunito Bold, followed by a blinking underscore in accent color.

```
human_
```

| Variant | Usage |
|---------|-------|
| `human_` | Primary wordmark — sidebar header, splash screen |
| `h_` | Compact — app icon, favicon, taskbar, dock |
| Underscore blinks in digital contexts (1.2s step-end infinite) |

**Favicon / App Icon SVG** (use this exact SVG for all icon sizes):
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 120">
  <rect width="120" height="120" rx="24" fill="#0D0D0D"/>
  <text x="24" y="84" font-family="Nunito, sans-serif" font-weight="700" font-size="72" letter-spacing="-1">
    <tspan fill="#F5F5F3">h</tspan><tspan fill="#E85D3A">_</tspan>
  </text>
</svg>
```

Generate platform-specific icons from this:
- **Windows:** .ico containing 16x16, 32x32, 48x48, 256x256
- **macOS:** .icns containing 16x16 through 1024x1024 @1x and @2x
- **Linux:** 128x128 and 512x512 PNG

### Color Tokens

```css
:root {
  /* Core brand */
  --accent: #E85D3A;
  --accent-dark: #C44A2D;
  --accent-light: #FFF0EC;
  --accent-dim: rgba(232, 93, 58, 0.08);
  --accent-border: rgba(232, 93, 58, 0.25);

  /* Dark theme (default) */
  --bg: #0D0D0D;
  --bg-raised: #141414;
  --bg-surface: #1A1A1A;
  --bg-hover: #1F1F1F;
  --text: #E8E8E4;
  --text-bright: #F5F5F3;
  --text-muted: #7A7A7A;
  --text-dim: #4A4A4A;
  --border: rgba(255, 255, 255, 0.08);
  --border-hover: rgba(255, 255, 255, 0.15);

  /* Semantic */
  --success: #2D8C5A;
  --error: #C43030;
  --warning: #D4940A;
  --info: #3B82F6;

  /* Code syntax (for .human files) */
  --syn-keyword: #E85D3A;       /* app, data, page, api, authentication, build, when */
  --syn-property: #93C5FD;      /* has, show, accepts, requires, check, create */
  --syn-string: #2D8C5A;        /* "pending", "admin" */
  --syn-type: #5EEAD4;          /* text, email, date, number, boolean */
  --syn-modifier: #A78BFA;      /* unique, required, optional, encrypted */
  --syn-conjunction: #4A4A4A;   /* which is, that, with, using */
  --syn-comment: #4A4A4A;       /* // comments */

  /* Code syntax (for generated TypeScript/JS) */
  --ts-keyword: #C084FC;        /* import, export, const, async, return, if */
  --ts-function: #67E8F9;       /* useState, useEffect, create, findMany */
  --ts-string: #FCD34D;         /* 'string values' */
  --ts-type: #5EEAD4;           /* string, number, Router, PrismaClient */
  --ts-number: #F9A8D4;         /* 42, 3.14 */
  --ts-comment: #4A4A4A;        /* // comments */
}

/* Light theme */
[data-theme="light"] {
  --bg: #FAFAF8;
  --bg-raised: #FFFFFF;
  --bg-surface: #F0F0EC;
  --bg-hover: #E8E8E4;
  --text: #2D2D2D;
  --text-bright: #1A1A1A;
  --text-muted: #6B6B6B;
  --text-dim: #9A9A9A;
  --border: rgba(0, 0, 0, 0.08);
  --border-hover: rgba(0, 0, 0, 0.15);
  --accent: #D04E2D;
  --accent-dark: #B8432A;
  --accent-light: #FFF5F2;
}
```

### Typography

| Role | Typeface | Weight | Fallback |
|------|----------|--------|----------|
| Logo | Nunito | 700 | Rounded sans-serif |
| Headings | Nunito Sans | 700 | system-ui |
| Body / UI | Nunito Sans | 400, 600 | system-ui |
| Code / Editor | JetBrains Mono | 400 | monospace |
| Terminal | JetBrains Mono | 400 | monospace |

Bundle these fonts with the app (don't rely on Google Fonts CDN).

### Spacing System

Base unit: 4px. Use tokens: `xs(4)`, `sm(8)`, `md(16)`, `lg(24)`, `xl(32)`, `2xl(48)`, `3xl(64)`.

### Component Design Rules

- **Border radius:** 8px for buttons/inputs, 12px for cards/panels, 24px for modals
- **Borders:** 1px `var(--border)`, hover `var(--border-hover)`
- **Shadows:** Minimal. Only modals/dropdowns: `0 8px 32px rgba(0,0,0,0.3)`
- **Accent usage:** ONLY for primary CTAs, active states, the underscore, links, and important highlights. Never as large area backgrounds.
- **Icons:** Line style, 1.5px stroke, rounded caps. Use Lucide icons.
- **Animations:** Subtle. 150ms transitions. No bouncing/elastic. The build process can be more animated.
- **Focus rings:** 2px accent outline, 2px offset

---

## Application Architecture

### Window Structure

Single window, four-column resizable layout with top bar and bottom panel:

```
┌─────────────────────────────────────────────────────────┐
│  Logo  │ Project ▾ │ ···spacer··· │ Git │ LLM │ Actions │ Avatar │  ← Top Bar (48px)
├────────┼───────────┼──────────────┼─────────────────────┤
│        │           │              │                     │
│ Project│  Prompt   │   Human      │  Generated          │
│ Tree   │  Chat     │   Editor     │  Output             │
│        │           │              │                     │
│ (200px)│  (320px)  │   (flex)     │  (320px)            │
│        │           │              │                     │
├────────┴───────────┴──────────────┴─────────────────────┤
│  Build Output / Terminal                        [toggle]│  ← Bottom Panel (collapsible)
└─────────────────────────────────────────────────────────┘
```

**Every column border is a drag-to-resize handle.** Columns have minimum widths (120px) and maximum widths (600px). The editor (col 3) is flex and absorbs remaining space.

**Each column can be popped out into a separate window** via a button in its header. When popped out, it becomes a floating window that stays synced with the main app via IPC.

### Top Bar (48px fixed)

Left section:
- `h_` logo icon (28x28, from logo-compact.svg) + "Human" text (Nunito Bold 15px) + version badge
- Divider
- Project selector dropdown (current project name, switch/open recent)

Center: spacer

Right section:
- **Git** button: shows branch name, dropdown with push/pull/create branch/disconnect
- **LLM** selector: dropdown with Anthropic Claude, OpenAI GPT-4, Google Gemini, Ollama (local), Groq, OpenRouter. "Configure API keys..." at bottom
- Divider
- **Check** button (cyan accent) — validates .human syntax, no file generation
- **Build** button (accent/primary) — full code generation pipeline
- **Run** button (green accent) — build + docker-compose up
- Divider
- **Avatar** (28x28 circle, user initials, accent background) — opens profile panel

**Check vs Build vs Run behavior:**

| Action | What it does | Duration | Build log shows |
|--------|-------------|----------|-----------------|
| Check | Lexer → Parser → IR → Analyzer | ~0.1s | Token count, block count, validation results |
| Build | Check + all code generators + quality engine | ~3-5s | Each generator step with timing, file counts |
| Run | Build + docker-compose up | ~8-15s | Build steps, then container startup (DB → Backend → Frontend), ends with "App running at localhost:5333" link |

### Column 1: Project Tree (200px default)

Header: "PROJECT" label + link folder button + new project button + pop-out button

**Link folder:** Opens native OS folder picker dialog. Selected folder becomes the project root. Shows all files in the tree.

**New project:** Opens modal with:
- Project name input
- Directory input (defaults to ~/Documents/human-projects/)
- Shows computed path: "A new folder will be created at ~/Documents/human-projects/{name}/"
- Cancel / Create buttons
- On create: runs `human init {name}` in the selected directory, adds to project tree

Tree features:
- Expand/collapse folders
- File type icons with color coding (human=accent, tsx=blue, sql=green, json=yellow, css=purple, yml=cyan)
- Click file to open in editor (col 3) or preview (col 4 for generated files)
- Right-click context menu: rename, delete, new file, new folder

Footer: build status badge (Idle / Checking... / Building... / Running / Build OK / Error)

### Column 2: Prompt Playground (320px default)

Header: "PROMPT" label with robot/chat icon + pop-out button

Chat interface:
- Messages with role labels ("You" in blue, "Human AI" in accent)
- User bubbles: dark surface background, standard border
- AI bubbles: faint accent-tinted background, accent-tinted border
- Auto-scroll to bottom on new messages
- Messages support markdown rendering (code blocks, lists, bold)

**File attachments:**
- Paperclip button next to input opens file picker
- Accepts: images (png, jpg, gif, webp), videos (mp4, webm — **max 50MB**), code files (.human, .ts, .tsx, .js, .json, .yml, .md, .txt), PDFs, CSVs
- If a video exceeds 50MB, show error toast: "Video files must be under 50MB. '{filename}' is {size}MB"
- Attached files appear as removable tags below the input area: `filename.ext (1.2MB) ×`
- Files are sent as context to the LLM along with the message

Quick-action chips above input: "Generate app", "Add feature", "Fix error", "Explain code"

Input area:
- Multi-line textarea (auto-grows, max 120px height)
- Shift+Enter for newline, Enter to send
- Placeholder: "Describe what you want to build..."
- Send button (accent background, arrow icon)

LLM integration:
- Uses the provider selected in the top bar
- System prompt includes the Human language spec, current .human file content, and IR output
- AI responses suggest .human code changes, which user can accept (auto-applies to editor)

### Column 3: Human Code Editor (flex, min 280px)

Tab bar:
- **app.human** tab (with h_ file icon) — the source editor
- **IR Preview** tab — shows compiled Intent IR (YAML format, read-only)
- **Changes** tab — git-style diff view of recent modifications
- Pop-out button
- Cursor position indicator: "Ln {n}, Col {n}"

Editor (CodeMirror 6):
- Custom language mode for .human syntax with highlighting rules:
  - **Block keywords** (accent, bold): `app`, `data`, `page`, `api`, `authentication`, `theme`, `build`, `policy`, `when`, `component`, `integrate`, `workflow`, `error`, `deploy`, `architecture`, `monitor`
  - **Property keywords** (blue): `has`, `belongs`, `requires`, `accepts`, `check`, `create`, `fetch`, `update`, `delete`, `respond`, `set`, `if`, `on`, `enable`, `method`, `expiration`, `show`, `each`, `clicking`, `there`, `send`, `navigate`
  - **Build/config keywords** (blue): `frontend`, `backend`, `database`, `deploy`, `design`, `primary`, `dark`, `border`, `spacing`
  - **Strings** (green): anything in double quotes
  - **Types** (cyan): `text`, `email`, `date`, `number`, `boolean`, `encrypted`
  - **Modifiers** (purple): `unique`, `required`, `optional`
  - **Conjunctions** (dim): `which is`, `that`, `with`, `to`, `using`, `from`, `and`, `or`, `for the`, `is a`, `is either`
- Line numbers
- Active line highlight
- Bracket matching
- Auto-indent (2 spaces)
- Search/replace (Cmd/Ctrl+F)
- Minimap (optional, toggleable)

IR Preview tab:
- Shows the compiled Intent IR in YAML format
- Syntax highlighted (cyan text on dark background)
- Read-only
- Updates live as user edits .human file (debounced 500ms)

Changes tab:
- Shows diff between current file and last saved/built version
- Green for additions, red for deletions
- Line numbers from both versions

### Column 4: Generated Output (320px default)

Header: "GENERATED OUTPUT" label + file count ("69 files") + pop-out button

Tree view of generated files:
- Same tree component as project tree
- Organized by stack: react/, node/, database/, docker/, quality/
- Click any file to open a code preview below the tree

Code preview:
- Back button to return to tree
- File name with type icon
- **Full syntax highlighting for TypeScript/JavaScript:**
  - Keywords (purple): `import`, `export`, `const`, `let`, `var`, `function`, `return`, `async`, `await`, `try`, `catch`, `if`, `else`, `new`, `default`
  - Functions (cyan): `useState`, `useEffect`, `create`, `findMany`, `authenticate`, etc.
  - Types (teal): `string`, `number`, `boolean`, `Router`, `PrismaClient`, etc.
  - Strings (yellow): single and double quoted strings
  - Numbers (pink): numeric literals
  - Comments (dim gray): `//` line comments and `/* */` block comments
- Read-only
- Line numbers

Footer: stack badges showing counts per technology — React {n}, Node {n}, SQL {n}, Docker {n}, Tests {n}

### Bottom Panel: Build Output (collapsible)

Toggle header: chevron + "Build Output" + status badge (Passed / Running... / Failed)

Content: xterm.js terminal emulator showing:
- Build step output with checkmarks and timing
- Docker container logs when running
- Clickable URLs (localhost:5333 opens in default browser)
- Color-coded: green for success, red for errors, yellow for warnings, gray for info

### Profile Panel (slide-in from right, 360px)

Triggered by clicking the avatar in the top bar. Slides in with overlay.

Sections:

**User Profile**
- Name field
- Email field
- Save changes button

**Password**
- Current password field
- New password field
- Reset password button

**MCP Connections**
- List of integrations: Figma, GitHub, Slack, Vercel, AWS, etc.
- Each shows: name, connected/disconnected status dot, connect/disconnect button
- Connected services show green dot
- Disconnected show gray dot with "Connect" button

**Subscription**
- Current plan name + "Active" badge
- Price and next billing date
- Payment Method subsection: card type + last 4 digits + "Change" button
- Billing History subsection: table with date, amount, status for recent invoices

**Logout** — link/button, red-ish styling

**Delete Account** — link/button, danger styling with confirmation modal ("This will permanently delete your account and all projects. Type your email to confirm.")

---

## File Structure

```
human-studio/
├── package.json
├── electron-builder.yml           # Build config for all platforms
├── tsconfig.json
├── tailwind.config.ts
├── vite.config.ts                 # Vite for renderer process
│
├── resources/                     # Platform assets
│   ├── icon.ico                   # Windows
│   ├── icon.icns                  # macOS
│   ├── icon.png                   # Linux (512x512)
│   ├── icon-16.png
│   ├── icon-32.png
│   ├── icon-256.png
│   └── tray-icon.png              # System tray (if needed)
│
├── src/
│   ├── main/                      # Electron main process
│   │   ├── index.ts               # App entry, window management
│   │   ├── ipc.ts                 # IPC handlers
│   │   ├── menu.ts                # Native menu (File, Edit, View, Build, Help)
│   │   ├── updater.ts             # Auto-update logic
│   │   ├── window.ts              # Window creation, pop-out management
│   │   └── services/
│   │       ├── compiler.ts        # Spawns `human check/build/run` CLI
│   │       ├── git.ts             # Git operations via simple-git
│   │       ├── project.ts         # Project CRUD, file watching
│   │       ├── llm.ts             # LLM API calls (Anthropic, OpenAI, etc.)
│   │       └── docker.ts          # Docker/compose operations
│   │
│   ├── renderer/                  # React app (renderer process)
│   │   ├── index.html
│   │   ├── main.tsx               # React root
│   │   ├── App.tsx                # Layout shell
│   │   │
│   │   ├── components/
│   │   │   ├── TopBar.tsx
│   │   │   ├── ProjectTree.tsx
│   │   │   ├── PromptChat.tsx
│   │   │   ├── HumanEditor.tsx    # CodeMirror wrapper
│   │   │   ├── OutputViewer.tsx
│   │   │   ├── BuildPanel.tsx     # xterm.js wrapper
│   │   │   ├── ProfilePanel.tsx
│   │   │   ├── ResizeHandle.tsx
│   │   │   ├── FileTree.tsx       # Shared tree component
│   │   │   ├── SyntaxHighlighter.tsx
│   │   │   └── ui/               # Primitive components
│   │   │       ├── Button.tsx
│   │   │       ├── Dropdown.tsx
│   │   │       ├── Modal.tsx
│   │   │       ├── Toast.tsx
│   │   │       ├── Badge.tsx
│   │   │       ├── Input.tsx
│   │   │       └── Avatar.tsx
│   │   │
│   │   ├── stores/               # Zustand stores
│   │   │   ├── project.ts
│   │   │   ├── editor.ts
│   │   │   ├── chat.ts
│   │   │   ├── build.ts
│   │   │   └── settings.ts
│   │   │
│   │   ├── hooks/
│   │   │   ├── useResize.ts      # Column resize logic
│   │   │   ├── useKeyboard.ts    # Global keyboard shortcuts
│   │   │   └── useTheme.ts
│   │   │
│   │   ├── lib/
│   │   │   ├── human-lang.ts     # CodeMirror language definition
│   │   │   ├── ts-highlight.ts   # TypeScript syntax highlighting
│   │   │   └── ipc.ts            # Typed IPC client
│   │   │
│   │   └── styles/
│   │       ├── globals.css        # CSS variables, base styles
│   │       ├── editor.css         # CodeMirror theme overrides
│   │       └── terminal.css       # xterm.js theme
│   │
│   └── preload/
│       └── index.ts               # contextBridge API exposure
│
├── scripts/
│   ├── generate-icons.ts          # SVG → ico/icns/png
│   └── notarize.ts                # macOS notarization
│
└── .github/
    └── workflows/
        └── release.yml            # Build + release for all platforms
```

---

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Cmd/Ctrl + B | Build |
| Cmd/Ctrl + Shift + B | Run (Build + Start) |
| Cmd/Ctrl + Shift + C | Check |
| Cmd/Ctrl + K | Focus prompt input |
| Cmd/Ctrl + 1/2/3/4 | Focus column 1/2/3/4 |
| Cmd/Ctrl + S | Save current file |
| Cmd/Ctrl + P | Quick open file |
| Cmd/Ctrl + Shift + P | Command palette |
| Cmd/Ctrl + F | Search in editor |
| Cmd/Ctrl + , | Open settings/profile |
| Cmd/Ctrl + ` | Toggle build panel |

---

## Native Menu Bar

```
File
  New Project...          Cmd+N
  Open Project...         Cmd+O
  Open Recent            →
  Link Folder...
  ─────────────
  Save                    Cmd+S
  Save All                Cmd+Shift+S
  ─────────────
  Settings                Cmd+,
  ─────────────
  Exit                    Cmd+Q

Edit
  Undo                    Cmd+Z
  Redo                    Cmd+Shift+Z
  ─────────────
  Cut                     Cmd+X
  Copy                    Cmd+C
  Paste                   Cmd+V
  Select All              Cmd+A
  ─────────────
  Find                    Cmd+F
  Replace                 Cmd+H

View
  Toggle Sidebar          Cmd+B (with no editor focus)
  Toggle Build Panel      Cmd+`
  Toggle Theme            (no shortcut)
  ─────────────
  Zoom In                 Cmd+=
  Zoom Out                Cmd+-
  Reset Zoom              Cmd+0
  ─────────────
  Focus Project           Cmd+1
  Focus Prompt            Cmd+2
  Focus Editor            Cmd+3
  Focus Output            Cmd+4

Build
  Check                   Cmd+Shift+C
  Build                   Cmd+Shift+B (note: different from View toggle)
  Run                     Cmd+Shift+R
  Stop                    Cmd+Shift+.
  ─────────────
  Clean Output

Help
  Documentation           (opens docs site)
  Language Spec            (opens spec page)
  Keyboard Shortcuts
  ─────────────
  Check for Updates
  About Human Studio
```

---

## Cross-Platform Packaging

### electron-builder.yml

```yaml
appId: com.humanlang.studio
productName: Human Studio
copyright: Copyright © 2026 Human Language Project

directories:
  output: dist
  buildResources: resources

mac:
  category: public.app-category.developer-tools
  icon: resources/icon.icns
  darkModeSupport: true
  hardenedRuntime: true
  gatekeeperAssess: false
  entitlements: build/entitlements.mac.plist
  entitlementsInherit: build/entitlements.mac.plist
  target:
    - target: dmg
      arch: [x64, arm64]
    - target: zip
      arch: [x64, arm64]

win:
  icon: resources/icon.ico
  target:
    - target: nsis
      arch: [x64]
    - target: portable
      arch: [x64]

linux:
  icon: resources/icon.png
  category: Development
  target:
    - target: AppImage
      arch: [x64]
    - target: deb
      arch: [x64]

nsis:
  oneClick: false
  allowToChangeInstallationDirectory: true
  installerIcon: resources/icon.ico
  uninstallerIcon: resources/icon.ico

dmg:
  contents:
    - x: 130
      y: 220
    - x: 410
      y: 220
      type: link
      path: /Applications
```

### Platform-Specific Considerations

**macOS:**
- Notarize with Apple Developer ID
- Support both Intel (x64) and Apple Silicon (arm64)
- Respect system dark mode preference
- Native window controls (traffic lights)
- Touch Bar support (Check / Build / Run buttons) if applicable

**Windows:**
- NSIS installer with custom install directory option
- Portable .exe option for no-install usage
- Add to PATH option during install (for `human` CLI)
- Windows 10+ required

**Linux (Ubuntu):**
- .AppImage for universal compatibility
- .deb for apt-based systems
- Desktop entry file with proper icon and categories
- Respect system GTK/Qt theme where possible

---

## Implementation Order

Build in this order. Each phase should be fully working before moving to the next.

### Phase 1: Shell (Day 1-2)
1. Scaffold Electron + Vite + React + TypeScript + Tailwind
2. Generate app icons from the SVG logo (all platform formats)
3. Implement the four-column layout with resize handles
4. CSS variables for full dark/light theme system
5. Top bar with logo, version badge, placeholder buttons
6. Native menu bar for all three platforms
7. Basic window management (size, position persistence)

### Phase 2: Editor Core (Day 3-4)
8. CodeMirror 6 integration with custom .human language mode
9. Full syntax highlighting for .human files (all keyword categories)
10. Line numbers, active line, cursor position
11. Tab bar (Editor / IR Preview / Changes)
12. File watching — reload when .human file changes on disk
13. Search/replace

### Phase 3: Project & Files (Day 5-6)
14. Project tree component with virtual scrolling
15. File type icons with color coding
16. "New project" modal (calls `human init`)
17. "Link folder" with native OS dialog
18. Recent projects list (persisted)
19. File CRUD (create, rename, delete) via context menu

### Phase 4: Build Pipeline (Day 7-8)
20. IPC bridge to spawn `human check/build/run` CLI commands
21. Build panel with xterm.js for real terminal output
22. Differentiated Check/Build/Run with proper log output
23. Build status in project tree footer and top bar
24. Output tree rendering after successful build
25. Generated file preview with TypeScript/JS syntax highlighting

### Phase 5: Prompt & AI (Day 9-10)
26. Chat UI with message bubbles, role labels, markdown rendering
27. File attachment system (drag-drop + picker, 50MB video limit)
28. Quick-action chips
29. LLM provider selector + API key configuration modal
30. LLM integration: send .human context + user message, receive suggestions
31. "Accept suggestion" flow: AI response → apply to editor

### Phase 6: Git & Profile (Day 11-12)
32. Git integration via simple-git (branch display, push/pull)
33. Profile slide-in panel (user, password, MCP, subscription, billing)
34. Settings persistence (theme, column widths, recent projects, API keys encrypted)
35. Pop-out windows for each column

### Phase 7: Polish & Ship (Day 13-14)
36. Keyboard shortcuts (all from the table above)
37. Command palette (Cmd+Shift+P)
38. Loading states, error states, empty states for every view
39. Onboarding flow for first launch
40. Auto-updater configuration
41. CI/CD pipeline for Windows/Mac/Linux releases
42. Performance optimization (lazy loading, virtual scrolling, debouncing)

---

## Critical UX Details

### Error Toast System
- Appears centered at bottom, 80px up from edge
- Red background for errors, accent for info, green for success
- Auto-dismisses after 3 seconds
- Multiple toasts stack vertically

### Empty States
- Project tree with no project: "Open a project or create a new one" + buttons
- Prompt with no messages: centered illustration + "Describe what you want to build"
- Output with no build: "Run a build to see generated code"
- Build panel with no output: "Check, Build, or Run your project"

### Loading States
- Check: pulsing cyan dot + "Checking..."
- Build: pulsing accent dot + "Building..." with progress in build log
- Run: pulsing green dot + "Starting..." then "Running" when containers are up
- LLM response: animated dots in chat + "Human AI is thinking..."

### Window Title
Format: `{filename} — {project} — Human Studio`
Example: `app.human — taskflow — Human Studio`

### Persistence (electron-store)
Save and restore:
- Window position and size
- Column widths
- Active tab in editor
- Theme preference
- Recent projects list (max 10)
- Last opened project
- LLM provider selection
- API keys (encrypted with keytar/safeStorage)
- Git credentials

---

## Design Quality Checklist

Before shipping, verify every screen against these criteria:

- [ ] Accent color (#E85D3A) used ONLY for: primary buttons, active states, links, the underscore, important highlights
- [ ] No other colors used as accent — the UI is monochrome + accent
- [ ] Logo renders correctly at all sizes (icon tray, title bar, splash)
- [ ] Both dark and light themes work completely
- [ ] All text is legible (contrast ratio ≥ 4.5:1)
- [ ] All interactive elements have hover, active, and focus states
- [ ] Resize handles work smoothly with no layout jumps
- [ ] All dropdowns close when clicking outside
- [ ] Keyboard navigation works for every interactive element
- [ ] Empty states exist for every view
- [ ] Error states are friendly, specific, and constructive
- [ ] No UI element uses more than 150ms transition
- [ ] Fonts are bundled, not loaded from CDN
- [ ] App icon looks crisp on all platforms and sizes
- [ ] Window remembers position/size across restarts

---

## Reference

- **Brand Guidelines:** `brand/BRAND_GUIDELINES.md` in the human repo
- **Website (for visual reference):** https://barun-bash.github.io/human/index.html
- **Design CSS variables:** `docs/shared.css` in the human repo
- **Logo SVGs:** `brand/logo-compact.svg`, `brand/logo-primary-dark.svg`, `brand/logo-primary-light.svg`
- **Favicon source:** The `<link rel="icon">` in `docs/index.html`
- **CLI commands:** `human check`, `human build`, `human run`, `human init`, `human version`, `human doctor`
- **Language spec:** `docs/language-spec.html` or `LANGUAGE_SPEC.md`

---

*Build this app the way the Human language is built: clear, warm, honest, and with craftsmanship in every detail.*
