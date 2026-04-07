package wanikani

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchSubjects_SinglePage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer test-token")
		}
		if got := r.Header.Get("Wanikani-Revision"); got != Revision {
			t.Errorf("Wanikani-Revision header = %q, want %q", got, Revision)
		}

		resp := Collection[Resource[SubjectData]]{
			Object:        "collection",
			TotalCount:    1,
			DataUpdatedAt: "2026-04-01T00:00:00Z",
			Pages:         Pages{PerPage: 500},
			Data: []Resource[SubjectData]{
				{
					ID:     100,
					Object: "vocabulary",
					Data: SubjectData{
						Characters:   "大きい",
						PartOfSpeech: []string{"い adjective"},
						Level:        3,
						Readings:     []Reading{{Reading: "おおきい", Primary: true}},
						Meanings:     []Meaning{{Meaning: "big", Primary: true, AcceptedAnswer: true}},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.BaseURL = server.URL

	subjects, updatedAt, err := client.FetchSubjects(context.Background(), []string{"vocabulary"}, "")
	if err != nil {
		t.Fatalf("FetchSubjects error: %v", err)
	}

	if updatedAt != "2026-04-01T00:00:00Z" {
		t.Errorf("updatedAt = %q, want %q", updatedAt, "2026-04-01T00:00:00Z")
	}

	if len(subjects) != 1 {
		t.Fatalf("got %d subjects, want 1", len(subjects))
	}
	if subjects[0].Data.Characters != "大きい" {
		t.Errorf("characters = %q, want %q", subjects[0].Data.Characters, "大きい")
	}
}

func TestFetchSubjects_Pagination(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		switch callCount {
		case 1:
			// Page 1: return one item with next_url
			nextURL := r.Host // Will be replaced below
			_ = nextURL
			resp := Collection[Resource[SubjectData]]{
				Object:        "collection",
				TotalCount:    2,
				DataUpdatedAt: "2026-04-01T00:00:00Z",
				Pages:         Pages{PerPage: 1},
				Data: []Resource[SubjectData]{
					{ID: 1, Data: SubjectData{Characters: "大きい"}},
				},
			}
			// Set next_url to the same server
			next := "http://" + r.Host + "/subjects?page_after_id=1"
			resp.Pages.NextURL = &next
			json.NewEncoder(w).Encode(resp)

		case 2:
			// Page 2: return one item with no next_url
			resp := Collection[Resource[SubjectData]]{
				Object:        "collection",
				TotalCount:    2,
				DataUpdatedAt: "2026-04-01T00:00:00Z",
				Pages:         Pages{PerPage: 1},
				Data: []Resource[SubjectData]{
					{ID: 2, Data: SubjectData{Characters: "小さい"}},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.BaseURL = server.URL

	subjects, _, err := client.FetchSubjects(context.Background(), []string{"vocabulary"}, "")
	if err != nil {
		t.Fatalf("FetchSubjects error: %v", err)
	}

	if len(subjects) != 2 {
		t.Fatalf("got %d subjects, want 2", len(subjects))
	}
	if subjects[0].Data.Characters != "大きい" || subjects[1].Data.Characters != "小さい" {
		t.Errorf("unexpected subjects: %v, %v", subjects[0].Data.Characters, subjects[1].Data.Characters)
	}
}

func TestFetchSubjects_UpdatedAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "updated_after=2026-04-01T00%3A00%3A00Z") &&
			!strings.Contains(r.URL.RawQuery, "updated_after=2026-04-01T00:00:00Z") {
			t.Errorf("expected updated_after in query, got: %s", r.URL.RawQuery)
		}

		resp := Collection[Resource[SubjectData]]{
			Object:        "collection",
			TotalCount:    0,
			DataUpdatedAt: "2026-04-02T00:00:00Z",
			Pages:         Pages{PerPage: 500},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.BaseURL = server.URL

	_, _, err := client.FetchSubjects(context.Background(), []string{"vocabulary"}, "2026-04-01T00:00:00Z")
	if err != nil {
		t.Fatalf("FetchSubjects error: %v", err)
	}
}

func TestFetchSubjects_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient("bad-token")
	client.BaseURL = server.URL

	_, _, err := client.FetchSubjects(context.Background(), []string{"vocabulary"}, "")
	if err == nil {
		t.Fatal("expected error for 401, got nil")
	}
	if !strings.Contains(err.Error(), "invalid API token") {
		t.Errorf("error = %q, want it to contain 'invalid API token'", err.Error())
	}
}

func TestFetchSubjects_RateLimitRetry(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		resp := Collection[Resource[SubjectData]]{
			Object:        "collection",
			TotalCount:    0,
			DataUpdatedAt: "2026-04-01T00:00:00Z",
			Pages:         Pages{PerPage: 500},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.BaseURL = server.URL

	_, _, err := client.FetchSubjects(context.Background(), []string{"vocabulary"}, "")
	if err != nil {
		t.Fatalf("FetchSubjects error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (1 retry), got %d", callCount)
	}
}
