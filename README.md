# 🐷 AI Gamification Elearning-Server — Backend

> Open source backend for the **AI Gamification Elearning Server** e-learning platform: a gamified, AI-powered environment to learn web and backend development.

---

## 📖 About the project

**AI Gamification Elearning Server** is an open source, gamified e-learning platform designed to teach software development in a practical and engaging way. The platform's initial focus is **web and backend development**, with a roadmap to expand into mobile development and artificial intelligence.

The platform is built around two core ideas:

- **Learning by doing** — students write real code, solve real challenges, and progress through levels as they build skills.
- **AI-powered feedback** — an integrated AI code reviewer analyses the student's submissions and provides contextual, constructive feedback in real time, acting as a always-available mentor.

This repository is the **Go backend** of the AI Gamification Elearning Server platform. It exposes the REST API consumed by the frontend and handles all business logic, persistence, and integration with external services — including the AI code review engine.

This project is **free to use and open to contributions**. See the [license](#license) section for details.

---

## 🗺️ Roadmap highlights

| Phase | Scope |
|---|---|
| v1 | Web & backend development tracks |
| v2 | AI code reviewer integration |
| v3 | Mobile development track |
| v4 | AI / ML development track |

---

## ⚙️ Tech stack

| Concern | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Cache | Redis |
| Messaging | Kafka |
| AI integration | Anthropic API (Claude) |
| Logger | [zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |

---

## 🚀 Getting started

```bash
# Clone the repository
git clone https://github.com/your-org/ai-gamification-elearning-server.git
cd ai-gamification-elearning-server

# Copy and configure environment
cp .env.sample .env.local
# Edit .env with your values (optional: create .env.local for local overrides)

# Run the server
make run

# Run with live-reload
make run/watch
```

---

## 🤝 Contributing

AI Gamification Elearning Server is an open source project and contributions are welcome. Please open an issue before submitting a pull request so we can discuss the proposed change first.

When contributing, make sure your code respects the layer boundaries described in this document. A PR that imports GORM into the domain layer or places business logic inside a controller will not be merged.

---

## 📄 License

This project is released under the [MIT License](LICENSE). You are free to use, modify, and distribute it.
