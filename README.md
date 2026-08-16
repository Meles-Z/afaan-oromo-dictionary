# 📖 Kitaaba Jechootaa — Afaan Oromo Dictionary

A fast, fully offline English ↔ Afaan Oromo desktop dictionary, built as a native cross-platform app.

No internet connection. No ads. No accounts. Just a dictionary that works.

---

## Why I Built This

I searched for a proper Afaan Oromo desktop dictionary — something I could give my sister so she could learn and look up words without needing an internet connection every time. I couldn't find one that actually worked well, worked offline, or felt like a real application instead of a bloated web page wrapped in Electron.

I'm a backend and systems engineer by trade, not a desktop or frontend developer. But once I decided this needed to exist, I built the first working version in a matter of hours — a real, native, searchable, bidirectional dictionary running fully on-device.

The dictionary data itself came from a PDF/Word document I had access to, which I parsed and imported into a proper local database. It's a solid start, but it's incomplete. **If you have access to a more complete Afaan Oromo word list — PDF, Word document, spreadsheet, or anything else — I would genuinely appreciate you sharing it.** The more complete the dictionary, the more people it can help.

---

## Features

- 🔄 **Bidirectional search** — English → Afaan Oromo and Afaan Oromo → English, with a single tap to switch direction
- ⚡ **Full-text search (SQLite FTS5)** — fast, ranked, relevance-based results, not a slow linear scan
- 📴 **100% offline** — all data lives in a local SQLite database bundled with the app; no network calls, ever
- 🎯 **Smart matching** — exact match → prefix match → contains match, with automatic fallback so typos and partial words still return results
- 📚 **Detailed entries** — part of speech, pronunciation, and example sentences in both languages where available
- 🖥️ **Native desktop app** — not a browser tab, not Electron bloat; a real lightweight native binary for Windows, macOS, and Linux

---

## Tech Stack

| Layer | Technology |
|---|---|
| Desktop framework | [Wails v2](https://wails.io) — Go backend + native OS webview, no Electron/Chromium bundling |
| Backend | Go 1.25, layered architecture (repository → service → app bindings) |
| Database | SQLite ([`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite), pure-Go, no CGO) with [FTS5](https://sqlite.org/fts5.html) full-text search |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate), embedded via `go:embed` — schema ships inside the binary, no external files |
| Frontend | Svelte 5 + TypeScript + Vite |
| Data import | Python-based parsing pipeline to extract and clean dictionary entries from source documents |

### Architecture

┌─────────────────────┐
│ Svelte 5 UI │ ← search, direction toggle, word detail view
└──────────┬───────────┘
│ Wails-generated JS bindings (typed, auto-generated)
┌──────────▼───────────┐
│ app.go │ ← thin binding layer exposed to frontend
└──────────┬───────────┘
┌──────────▼───────────┐
│ Service layer │ ← validation, business rules
└──────────┬───────────┘
┌──────────▼───────────┐
│ Repository layer │ ← raw SQL, FTS5 queries
└──────────┬───────────┘
┌──────────▼───────────┐
│ SQLite (embedded) │ ← versioned migrations, local file, offline
└───────────────────────┘

---

## Development

**Live development** — hot reload for frontend changes, auto-rebuild for Go changes:

```bash
wails dev
```

This runs a Vite dev server with fast hot reload. A separate dev server also runs at `http://localhost:34115` — open it in a browser to call your Go methods directly from devtools.

**Production build:**

```bash
wails build
```

Produces a redistributable native binary in `build/bin/`.

---

## Data

Dictionary entries were extracted and parsed from an existing Afaan Oromo–English word document, cleaned, deduplicated, and imported into the local SQLite database via a Python parsing script. The dataset currently contains ~3,100 word pairs, covering common vocabulary across everyday topics.

**Have a more complete word list?** Please reach out — see below.

---

## Contact

Developed by **Meles Zewude**

📧 meles.zewde.tech@gmail.com
📱 +251-92-022-7833

If you have Afaan Oromo dictionary data you're willing to share, or want to contribute, I'd love to hear from you.

---

## License

MIT — see [LICENSE](./LICENSE) for details.