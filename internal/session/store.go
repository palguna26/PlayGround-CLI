package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrInvalidSessionID = errors.New("invalid session ID")
	ErrInvalidRepoPath  = errors.New("invalid repository path")
	ErrEmptyGoal        = errors.New("session goal cannot be empty")
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionLocked    = errors.New("session is locked by another process")
)

const (
	sessionDirName = ".pg"
	sessionSubDir  = "sessions"
	lockSuffix     = ".lock"
)

// Store manages session persistence on the filesystem
type Store struct {
	baseDir string // Repository root directory
	mu      sync.Mutex
}

// NewStore creates a new session store for the given repository
func NewStore(repoDir string) (*Store, error) {
	absPath, err := filepath.Abs(repoDir)
	if err != nil {
		return nil, fmt.Errorf("invalid repo path: %w", err)
	}

	return &Store{
		baseDir: absPath,
	}, nil
}

// getSessionDir returns the path to the .pg/sessions directory
func (s *Store) getSessionDir() string {
	return filepath.Join(s.baseDir, sessionDirName, sessionSubDir)
}

// getSessionPath returns the path to a specific session file
func (s *Store) getSessionPath(sessionID string) string {
	return filepath.Join(s.getSessionDir(), sessionID+".json")
}

// getLockPath returns the path to a session's lock file
func (s *Store) getLockPath(sessionID string) string {
	return filepath.Join(s.getSessionDir(), sessionID+lockSuffix)
}

// Save persists a session to disk using atomic write pattern
func (s *Store) Save(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := session.Validate(); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}

	// Ensure session directory exists
	sessionDir := s.getSessionDir()
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Marshal session to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Atomic write: write to temp file, then rename
	sessionPath := s.getSessionPath(session.ID)
	tempPath := sessionPath + ".tmp"

	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session: %w", err)
	}

	if err := os.Rename(tempPath, sessionPath); err != nil {
		os.Remove(tempPath) // Clean up temp file on failure
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// Load reads a session from disk
func (s *Store) Load(sessionID string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionPath := s.getSessionPath(sessionID)

	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to read session: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	return &session, nil
}

// List returns all session IDs in the store
func (s *Store) List() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionDir := s.getSessionDir()

	// If directory doesn't exist, return empty list
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	var sessionIDs []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Skip lock files and temp files
		if filepath.Ext(name) == ".json" {
			sessionID := name[:len(name)-5] // Remove .json extension
			sessionIDs = append(sessionIDs, sessionID)
		}
	}

	return sessionIDs, nil
}

// Delete removes a session from disk
func (s *Store) Delete(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionPath := s.getSessionPath(sessionID)

	if err := os.Remove(sessionPath); err != nil {
		if os.IsNotExist(err) {
			return ErrSessionNotFound
		}
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Also remove lock file if it exists
	lockPath := s.getLockPath(sessionID)
	os.Remove(lockPath) // Ignore errors - lock may not exist

	return nil
}

// AcquireLock attempts to acquire an exclusive lock on a session
func (s *Store) AcquireLock(sessionID string) (func() error, error) {
	lockPath := s.getLockPath(sessionID)

	// Ensure session directory exists
	sessionDir := s.getSessionDir()
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	// Try to create lock file exclusively
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil, ErrSessionLocked
		}
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	// Write current process ID to lock file
	fmt.Fprintf(lockFile, "%d\n%s", os.Getpid(), time.Now().Format(time.RFC3339))
	lockFile.Close()

	// Return unlock function
	unlock := func() error {
		return os.Remove(lockPath)
	}

	return unlock, nil
}

// GetActiveSessionID returns the currently active session ID, if any
// This is stored in a special .pg/active file
func (s *Store) GetActiveSessionID() (string, error) {
	activePath := filepath.Join(s.baseDir, sessionDirName, "active")

	data, err := os.ReadFile(activePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No active session
		}
		return "", fmt.Errorf("failed to read active session: %w", err)
	}

	return string(data), nil
}

// SetActiveSessionID sets the currently active session
func (s *Store) SetActiveSessionID(sessionID string) error {
	activePath := filepath.Join(s.baseDir, sessionDirName, "active")

	// Ensure .pg directory exists
	pgDir := filepath.Join(s.baseDir, sessionDirName)
	if err := os.MkdirAll(pgDir, 0755); err != nil {
		return fmt.Errorf("failed to create .pg directory: %w", err)
	}

	if err := os.WriteFile(activePath, []byte(sessionID), 0644); err != nil {
		return fmt.Errorf("failed to set active session: %w", err)
	}

	return nil
}

// GenerateSessionID creates a new unique session ID
func (s *Store) GenerateSessionID() (string, error) {
	sessionIDs, err := s.List()
	if err != nil {
		return "", err
	}

	// Simple counter-based ID generation: pg-1, pg-2, etc.
	nextNum := len(sessionIDs) + 1
	for {
		candidate := fmt.Sprintf("pg-%d", nextNum)

		// Check if this ID exists
		exists := false
		for _, id := range sessionIDs {
			if id == candidate {
				exists = true
				break
			}
		}

		if !exists {
			return candidate, nil
		}

		nextNum++
	}
}

// CopyFrom copies a reader's content to a writer (utility function)
func copyFrom(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}
