package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/the-witch-king/conjugreater/vocabulary"
	"github.com/the-witch-king/conjugreater/wanikani"
)

//go:embed web/build/*
var buildFS embed.FS

var subjectTypes = []string{"vocabulary", "kana_vocabulary"}

var (
	fetchMu    sync.Mutex
	fetchRunning bool
)

func main() {
	static, err := fs.Sub(buildFS, "web/build")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8666")
	if err != nil {
		// Port in use, pick a random available one
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/config", handleGetConfig)
	mux.HandleFunc("POST /api/config", handlePostConfig)
	mux.HandleFunc("POST /api/fetch", handleFetch)
	mux.HandleFunc("GET /vocabulary.json", handleVocabulary)
	mux.Handle("/", http.FileServer(http.FS(static)))

	addr := listener.Addr().String()
	url := "http://" + addr
	fmt.Printf("Serving at %s\n", url)

	openBrowser(url)

	log.Fatal(http.Serve(listener, mux))
}

func dataDir() string {
	return filepath.Join(".", "data")
}

func envPath() string {
	return filepath.Join(".", ".env")
}

func vocabPath() string {
	return filepath.Join(dataDir(), "vocabulary.json")
}

func readAPIToken() string {
	data, err := os.ReadFile(envPath())
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "WANIKANI_API_TOKEN=") {
			return strings.TrimPrefix(line, "WANIKANI_API_TOKEN=")
		}
	}
	return ""
}

func writeAPIToken(token string) error {
	content := "WANIKANI_API_TOKEN=" + token + "\n"
	return os.WriteFile(envPath(), []byte(content), 0o600)
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	token := readAPIToken()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"has_token": token != "",
		"token":     token,
	})
}

func handlePostConfig(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	body.Token = strings.TrimSpace(body.Token)
	if body.Token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	if err := writeAPIToken(body.Token); err != nil {
		http.Error(w, "failed to save token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	fetchMu.Lock()
	if fetchRunning {
		fetchMu.Unlock()
		http.Error(w, "fetch already in progress", http.StatusConflict)
		return
	}
	fetchRunning = true
	fetchMu.Unlock()

	defer func() {
		fetchMu.Lock()
		fetchRunning = false
		fetchMu.Unlock()
	}()

	token := readAPIToken()
	if token == "" {
		http.Error(w, "API token not configured", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client := wanikani.NewClient(token)
	outputPath := vocabPath()

	// Try to load existing output for incremental updates
	var existing *vocabulary.Output
	if data, err := os.ReadFile(outputPath); err == nil {
		var out vocabulary.Output
		if json.Unmarshal(data, &out) == nil {
			existing = &out
		}
	}

	var subjectsAfter, assignmentsAfter string
	if existing != nil {
		subjectsAfter = existing.SubjectsUpdatedAt
		assignmentsAfter = existing.AssignmentsUpdatedAt
		log.Printf("Incremental update (subjects after %s, assignments after %s)", subjectsAfter, assignmentsAfter)
	} else {
		log.Println("Full fetch")
	}

	subjects, subjectsUpdatedAt, err := client.FetchSubjects(ctx, subjectTypes, subjectsAfter)
	if err != nil {
		http.Error(w, "failed to fetch subjects: "+err.Error(), http.StatusBadGateway)
		return
	}

	assignments, assignmentsUpdatedAt, err := client.FetchAssignments(ctx, subjectTypes, assignmentsAfter)
	if err != nil {
		http.Error(w, "failed to fetch assignments: "+err.Error(), http.StatusBadGateway)
		return
	}

	assignmentMap := make(map[int]wanikani.AssignmentData, len(assignments))
	for _, a := range assignments {
		assignmentMap[a.Data.SubjectID] = a.Data
	}

	newWords := vocabulary.TransformSubjects(subjects, assignmentMap)

	var finalWords []vocabulary.Word
	if existing != nil && len(existing.Words) > 0 {
		finalWords = vocabulary.MergeWords(existing.Words, newWords)
	} else {
		finalWords = newWords
	}

	if subjectsUpdatedAt == "" && existing != nil {
		subjectsUpdatedAt = existing.SubjectsUpdatedAt
	}
	if assignmentsUpdatedAt == "" && existing != nil {
		assignmentsUpdatedAt = existing.AssignmentsUpdatedAt
	}

	out := vocabulary.Output{
		GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
		SubjectsUpdatedAt:    subjectsUpdatedAt,
		AssignmentsUpdatedAt: assignmentsUpdatedAt,
		Words:                finalWords,
	}

	if err := os.MkdirAll(dataDir(), 0o755); err != nil {
		http.Error(w, "failed to create data dir: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		http.Error(w, "failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		http.Error(w, "failed to write output: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"ok":          true,
		"total_words": len(finalWords),
		"new_words":   len(newWords),
	})
}

// handleVocabulary serves vocabulary.json from disk, or an empty word list if not yet fetched.
func handleVocabulary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if data, err := os.ReadFile(vocabPath()); err == nil {
		w.Write(data)
		return
	}
	w.Write([]byte(`{"words":[]}`))
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
