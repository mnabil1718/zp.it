# ⚡ Zp.it — Fast, Simplified URL Shortener

> *Read: "Zip it!"* — A free, fast, and open-source URL shortener. No logins, no subscriptions, no nonsense.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-1.x-3D72D7?style=flat&logo=htmx&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Cache-FF4438?style=flat&logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

---

## Live Demo
[https://zpit.up.railway.app](https://zpit.up.railway.app)

## Screenshots
> Home
<img width="1848" height="917" alt="Screenshot from 2026-03-07 15-59-54" src="https://github.com/user-attachments/assets/410e4fee-d796-4a2e-a370-ed76113d2290" />



---

## ✨ Features

- **⚡ Fast Lookup & Redirect** — Redis-powered caching ensures near-instant redirects with minimal latency
- **📊 Click Counter** — Track how many times each shortlink has been visited in real time
- **✏️ Custom Aliases** — Choose your own memorable slug instead of a random one (e.g. `zp.it/my-link`)
- **📱 QR Code Generator** — Instantly generate a scannable QR code for any shortened link
- **🌐 No Account Required** — Just paste, shorten, and share. Zero friction.
- **🖥️ Hypermedia-driven UI** — Snappy, SPA-like experience powered by HTMX — no heavy JavaScript framework

---

## 🛠️ Tech Stack

| Layer     | Technology                          |
|-----------|-------------------------------------|
| Backend   | [Go](https://go.dev/)               |
| Frontend  | [HTMX](https://htmx.org/) + [Tailwind CSS](https://tailwindcss.com/) |
| Database  | [SQLite](https://www.sqlite.org/)   |
| Cache     | [Redis](https://redis.io/)          |

> **Why this stack?** Go keeps the binary lean and fast. SQLite means zero-config persistence — just a file. Redis handles hot-path lookups so the database barely breaks a sweat. HTMX delivers interactivity without shipping a JS framework to the browser.

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Redis](https://redis.io/docs/getting-started/)
- Git

### Installation

1. **Clone the repository**
```bash
   git clone https://github.com/yourusername/zp.it.git
   cd zp.it
```

2. **Install dependencies**
```bash
   go mod tidy
```

3. **Configure environment**
```bash
   cp .env.example .env
```
   Then edit `.env` with your settings:
```env
   HOST=http://localhost:8080/
   PORT=8080
   REDIS_URL=redis://localhost:6379
   DB_PATH=./zpit.db
```

4. **Run the app**
```bash
   go run ./cmd/main.go
```

5. **Open your browser**
```
   http://localhost:8080
```

### Docker (optional)
```bash
# Coming soon
docker compose up
```

---

## 📖 Usage

### Shorten a URL
1. Paste any long URL into the input field
2. Optionally enable **custom alias** and enter your preferred slug
3. Optionally check **Generate QR code** to get a scannable image alongside your link
4. Hit **Shorten Link** — done!

### Track Clicks
Navigate to `/counter` or click **"See shortlink click counter"** to look up click stats for any shortlink.

### Custom Alias Rules
- Allowed characters: letters, numbers, hyphens (`-`), underscores (`_`)
- Length: 3–24 characters
- First-come, first-served — aliases are unique

---

## 🗂️ Project Structure
```
zp.it/
├── cmd/
│   └── main.go          # Entry point
├── internal/
│   ├── handler/         # HTTP handlers
│   ├── store/           # SQLite data layer
│   └── cache/           # Redis cache layer
├── ui/                  # Go HTML templates
│   ├── base.html
│   ├── index.html
│   └── ...
├── static/              # Static assets
├── .env.example
├── go.mod
└── README.md
```
---

## 🤝 Contributing

Contributions are what make open source great — all kinds of input are welcome, whether it's a bug report, a feature idea, or a pull request!

1. **Fork** the repository
2. **Create** a feature branch
```bash
   git checkout -b feature/your-feature-name
```
3. **Commit** your changes
```bash
   git commit -m "feat: add your feature"
```
4. **Push** to your branch
```bash
   git push origin feature/your-feature-name
```
5. **Open a Pull Request** — describe what you changed and why

Please make sure your code is formatted with `gofmt` and passes any existing tests before submitting.

### Ideas for Contributions
- [ ] Expiring / time-limited links
- [ ] Link management dashboard
- [ ] REST API with JSON responses
- [ ] Analytics charts on the counter page

---

## 🐛 Reporting Issues

Found a bug? Have a feature request? [Open an issue](https://github.com/yourusername/zp.it/issues) and describe it clearly. Include steps to reproduce for bugs if possible.

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](./LICENSE) file for details.

---

## 🙏 Acknowledgements

- [HTMX](https://htmx.org/) for making server-side rendering fun again
- [go-redis](https://github.com/redis/go-redis) for the Redis client
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite bindings

---
<p align="center">Made with ❤️ and Go · <a href="https://zp.it">zp.it</a></p>
