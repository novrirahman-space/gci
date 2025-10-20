package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ===== Models =====
type Class struct {
	ID string `json:"id"`
	ClassName string `json:"class_name"`
	Teacher string `json:"teacher"`
}

type Task struct {
	ID string `json:"id"`
	ClassID string `json:"class_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	DueAt *time.Time `json:"due_at,omitempty"`
	IsClosed bool `json:"is_closed"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
}

// ===== Store (thread-safe) =====
type Store struct {
	muClasses sync.RWMutex
	classes map[string]Class

	muTasks sync.RWMutex
	tasks map[string]Task
}

func NewStore() *Store {
	return &Store{
		classes: make(map[string]Class),
		tasks: make(map[string]Task),
	}
}

// ===== Utilities =====
func newID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{
		"error": http.StatusText(status),
		"message": msg,
	})
}

func trim(s string) string { return strings.TrimSpace(s) }

// pathID mengekstrak ID dari path
func pathID(path, base string) (string, bool) {
	if !strings.HasPrefix(path, base) {
		return "", false
	}
	id := strings.TrimPrefix(path, base)
	if id == "" || strings.Contains(id, "/") {
		return "", false
	}
	return id, true
}

func parseRFC339Ptr(s string) (*time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ===== Handlers: Classes =====
func (s *Store) handleClasses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var p struct {
			ClassName string `json:"class_name"`
			Teacher string `json:"teacher"`
		}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if trim(p.ClassName) == "" || trim(p.Teacher) == "" {
			writeError(w, http.StatusBadRequest, "class_name and teacher are required")
			return
		}

		id, err := newID()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to generate id")
			return
		}

		c := Class{ID: id, ClassName: p.ClassName, Teacher: p.Teacher}

		s.muClasses.Lock()
		s.classes[c.ID] = c
		s.muClasses.Unlock()

		writeJSON(w, http.StatusCreated, c)

	case http.MethodGet:
		s.muClasses.RLock()
		out := make([]Class, 0, len(s.classes))
		for _, c := range s.classes {
			out = append(out, c)
		}
		s.muClasses.RUnlock()
		writeJSON(w, http.StatusOK, out)
	
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Store) handleClassByID(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/classes/")
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.muClasses.RLock()
		c, ok := s.classes[id]
		s.muClasses.RUnlock()
		if !ok {
			writeError(w, http.StatusNotFound, "class not found")
			return
		}
		writeJSON(w, http.StatusOK, c)

	case http.MethodPut:
		var p struct {
			ClassName string `json:"class_name"`
			Teacher string `json:"teacher"`
		}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if trim(p.ClassName) == "" || trim(p.Teacher) == "" {
			writeError(w, http.StatusBadRequest, "class_name and teacher are required")
			return
		}

		s.muClasses.Lock()
		c, ok := s.classes[id]
		if !ok {
			s.muClasses.Unlock()
			writeError(w, http.StatusNotFound, "class not found")
			return
		}
		c.ClassName = p.ClassName
		c.Teacher = p.Teacher
		s.classes[id] = c
		s.muClasses.Unlock()

		writeJSON(w, http.StatusOK, c)

	case http.MethodDelete:
		// hapus class
		s.muClasses.Lock()
		_, ok := s.classes[id]
		if ok {
			delete(s.classes, id)
		}
		s.muClasses.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "class not found")
			return
		}

		// bersihkan tasks milik class tsb
		s.muTasks.Lock()
		for tid, t := range s.tasks {
			if t.ClassID == id {
				delete(s.tasks, tid)
			}
		}
		s.muTasks.Unlock()

		writeJSON(w, http.StatusOK, map[string]string{"message": "Class deleted successfully"})

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// GET /classes/{id}/tasks
func (s *Store) handleClassTasks(w http.ResponseWriter, r * http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	// path contoh: /classes/abc/tasks
	if !strings.HasSuffix(r.URL.Path, "/tasks") {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	base := strings.TrimSuffix(r.URL.Path, "/tasks")
	id, ok := pathID(base, "/classes/")
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	// pastikan class ada
	s.muClasses.RLock()
	_, exists := s.classes[id]
	s.muClasses.RUnlock()
	if !exists {
		writeError(w, http.StatusNotFound, "class not found")
		return
	}

	// kumpulkan task milik class
	s.muTasks.RLock()
	out := make([]Task, 0)
	for _, t := range s.tasks {
		if t.ClassID == id {
			out = append(out, t)
		}
	}
	s.muTasks.RUnlock()

	writeJSON(w, http.StatusOK, out)
}

// ===== Handlers: Tasks =====
func (s *Store) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var p struct {
			ClassID string `json:"class_id"`
			Title string `json:"title"`
			Description string `json:"description"`
			DueAt string `json:"due_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if trim(p.ClassID) == "" || trim(p.Title) == "" {
			writeError(w, http.StatusBadRequest, "class_id and title are required")
			return
		}

		// validasi class_id
		s.muClasses.RLock()
		_, ok := s.classes[p.ClassID]
		s.muClasses.RUnlock()
		if !ok {
			writeError(w, http.StatusBadRequest, "class_id not found")
			return
		}

		dueAt, err := parseRFC339Ptr(p.DueAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid due_at (use RFC3339)")
			return
		}

		id, err := newID()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to generate id")
			return
		}

		t := Task{
			ID: id,
			ClassID: p.ClassID,
			Title: p.Title,
			Description: p.Description,
			DueAt: dueAt,
			IsClosed: false,
			ClosedAt: nil,
		}

		s.muTasks.Lock()
		s.tasks[t.ID] = t
		s.muTasks.Unlock()

		writeJSON(w, http.StatusCreated, t)

	case http.MethodGet:
		s.muTasks.RLock()
		out := make([]Task, 0, len(s.tasks))
		for _, t := range s.tasks {
			out = append(out, t)
		}
		s.muTasks.RUnlock()
		writeJSON(w, http.StatusOK, out)
		
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Store) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/tasks/")
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.muTasks.RLock()
		t, ok := s.tasks[id]
		s.muTasks.RUnlock()
		if !ok {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeJSON(w, http.StatusOK, t)

	case http.MethodPut:
		var p struct {
			ClassID string `json:"class_id"`
			Title string `json:"title"`
			Description string `json:"description"`
			DueAt string `json:"due_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if trim(p.ClassID) == "" || trim(p.Title) == "" {
			writeError(w, http.StatusBadRequest, "class_id and title are required")
			return
		}

		// pastikan class_id valid
		s.muClasses.RLock()
		_, ok := s.classes[p.ClassID]
		s.muClasses.RUnlock()
		if !ok {
			writeError(w, http.StatusBadRequest, "class_id not found")
			return
		}

		dueAt, err := parseRFC339Ptr(p.DueAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid due_at (use RFC3339)")
			return
		}

		s.muTasks.Lock()
		t, ok := s.tasks[id]
		if !ok {
			s.muTasks.Unlock()
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		t.ClassID = p.ClassID
		t.Title = p.Title
		t.Description = p.Description
		t.DueAt = dueAt
		s.tasks[id] = t
		s.muTasks.Unlock()

		writeJSON(w, http.StatusOK, t)

	case http.MethodDelete:
		s.muTasks.Lock()
		_, ok := s.tasks[id]
		if ok {
			delete(s.tasks, id)
		}
		s.muTasks.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Store) handleTaskClose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", "PATCH")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !strings.HasSuffix(r.URL.Path, "/close") {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	base := strings.TrimSuffix(r.URL.Path, "/close")
	id, ok := pathID(base, "/tasks/")
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	s.muTasks.Lock()
	t, exists := s.tasks[id]
	if !exists {
		s.muTasks.Unlock()
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if !t.IsClosed {
		now := time.Now()
		t.IsClosed = true
		t.ClosedAt = &now
		s.tasks[id] = t
	}
	out := t
	s.muTasks.Unlock()
	
	writeJSON(w, http.StatusOK, out)
}

func (s *Store) handleTaskOpen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", "PATCH")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !strings.HasSuffix(r.URL.Path, "/open"){
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	base := strings.TrimSuffix(r.URL.Path, "/open")
	id, ok := pathID(base, "/tasks/")
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	s.muTasks.Lock()
	t, exists := s.tasks[id]
	if !exists {
		s.muTasks.Unlock()
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if t.IsClosed {
		t.IsClosed = false
		t.ClosedAt = nil
		s.tasks[id] = t
	}
	out := t
	s.muTasks.Unlock()

	writeJSON(w, http.StatusOK, out)
}

// ===== Server & Middleware =====
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "Content-Type")
		h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := NewStore()
	mux := http.NewServeMux()

	// classes
	mux.HandleFunc("/classes", store.handleClasses)
	mux.HandleFunc("/classes/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/tasks") {
			store.handleClassTasks(w, r)
			return
		}
		store.handleClassByID(w, r)
	})

	// tasks
	mux.HandleFunc("/tasks", store.handleTasks)
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/close"):
			store.handleTaskClose(w, r)
			return
		case strings.HasSuffix(r.URL.Path, "/open"):
			store.handleTaskOpen(w, r)
			return
		default:
			store.handleTaskByID(w, r)
		}
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: cors(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// run server
	go func ()  {
		log.Printf("[server] listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("[server] shutdown gracefully")
}