# 🧠 GCI — Go & Next.js Project

A small project built with **Go (Golang)** and **Next.js (TypeScript)** for practicing concurrency, debugging, and building clean RESTful APIs.

---

## 🏗️ Overview

This repository integrates both **backend (Go)** and **frontend (Next.js)** under a single monorepo structure.

| Layer | Tech Stack | Description |
|-------|-------------|--------------|
| **Backend** | Go 1.23+, `net/http`, `sync`, `context` | Demonstrates concurrency via worker pool, debugging with WaitGroup, and a REST API for managing classes and tasks. |
| **Frontend** | Next.js 15, TypeScript, TailwindCSS, React Query | Provides a user interface for creating, editing, and deleting classes and tasks through API integration. |

---

## 📂 Project Structure

```
gci/
├── cmd/
│   ├── workerpool/   # Go worker pool concurrency example
│   ├── debugfix/     # Debugging & sync.WaitGroup example
│   └── server/       # RESTful API for Classes & Tasks
└── gci-web/          # Next.js frontend (App Router + React Query)
```

---

## ⚙️ Backend (Go)

### Run Examples
```bash
# Worker Pool example
go run ./cmd/workerpool

# Debugging & WaitGroup example
go run ./cmd/debugfix

# RESTful API Server
go run ./cmd/server
```

### API Base URL
```
http://localhost:8080
```

### API Endpoints

#### 📘 Classes
| Method | Endpoint | Description |
|--------|-----------|-------------|
| `GET` | `/classes` | Get all classes |
| `POST` | `/classes` | Create new class |
| `GET` | `/classes/{id}` | Get class by ID |
| `PUT` | `/classes/{id}` | Update class |
| `DELETE` | `/classes/{id}` | Delete class |
| `GET` | `/classes/{id}/tasks` | Get all tasks by class |

#### 🧩 Tasks
| Method | Endpoint | Description |
|--------|-----------|-------------|
| `GET` | `/tasks` | Get all tasks |
| `POST` | `/tasks` | Create new task |
| `GET` | `/tasks/{id}` | Get task by ID |
| `PUT` | `/tasks/{id}` | Update task |
| `DELETE` | `/tasks/{id}` | Delete task |

### Key Features
✅ Thread-safe using `sync.RWMutex`  
✅ Graceful shutdown with `context.Context`  
✅ Built-in CORS middleware  
✅ JSON responses and error handling  
✅ In-memory data store (no external DB dependency)  

---

## 💻 Frontend (Next.js)

### Tech Stack
- **Next.js (App Router)**
- **TypeScript**
- **TailwindCSS**
- **TanStack React Query**

### Setup & Run
```bash
cd gci-web
npm install
npm run dev
```

Frontend runs at:
```
http://localhost:3000
```

### Environment Variables
Create `.env.local` inside `gci-web`:
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### Pages
| Path | Description |
|------|--------------|
| `/` | Home page |
| `/classes` | Manage classes (CRUD) |
| `/tasks` | Manage tasks (CRUD) |

### Features
✅ Real-time data sync with React Query  
✅ Simple UI 
✅ Form validation & inline feedback  
✅ Modular components for reusability  
✅ Fully connected to Go REST API  

---

## 🚀 Run Both Services

### 1️⃣ Start Backend
```bash
go run ./cmd/server
```

### 2️⃣ Start Frontend
```bash
cd gci-web
npm run dev
```

Then open the app:
👉 [http://localhost:3000](http://localhost:3000)

---

## 📦 Future Improvements

- [ ] Integrate MySQL or PostgreSQL for persistence  
- [ ] Add JWT-based authentication  
- [ ] Containerize with Docker  
- [ ] Add automated unit & API tests  
- [ ] CI/CD deployment with GitHub Actions or Vercel  

---

## 📜 License
MIT © 2025

---

## 🧩 Commit Convention
All commits follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` → new features  
- `fix:` → bug fixes  
- `docs:` → documentation changes  
- `refactor:` → code improvements  
- `chore:` → maintenance updates  

Example:
```bash
git commit -m "feat(api): implement CRUD for classes and tasks"
```