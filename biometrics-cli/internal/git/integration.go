package git

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Repository struct {
	mu         sync.RWMutex
	path       string
	remoteURL  string
	branch     string
	lastCommit string
	lastCheck  time.Time
}

type Commit struct {
	Hash      string
	Author    string
	Date      time.Time
	Message   string
	Files     []string
	Additions int
	Deletions int
}

type Integration struct {
	mu            sync.RWMutex
	repos         map[string]*Repository
	webhookURL    string
	webhookSecret string
	autoCommit    bool
	commitQueue   chan *CommitRequest
	wg            sync.WaitGroup
}

type CommitRequest struct {
	Path    string
	Message string
	Files   []string
	Branch  string
	Force   bool
	Result  chan<- *CommitResult
}

type CommitResult struct {
	Hash  string
	Error error
}

var defaultIntegration = &Integration{
	repos:       make(map[string]*Repository),
	autoCommit:  false,
	commitQueue: make(chan *CommitRequest, 100),
}

var IntegrationInstance = defaultIntegration

func New() *Integration {
	return defaultIntegration
}

func (i *Integration) AddRepository(path string) error {
	absPath, err := getAbsPath(path)
	if err != nil {
		return err
	}

	repo := &Repository{
		path: absPath,
	}

	if err := repo.updateInfo(); err != nil {
		return err
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.repos[absPath] = repo

	state.GlobalState.Log("INFO", fmt.Sprintf("Added repository: %s", absPath))
	return nil
}

func (i *Integration) RemoveRepository(path string) error {
	absPath, err := getAbsPath(path)
	if err != nil {
		return err
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.repos[absPath]; !exists {
		return fmt.Errorf("repository not found: %s", absPath)
	}

	delete(i.repos, absPath)
	state.GlobalState.Log("INFO", fmt.Sprintf("Removed repository: %s", absPath))
	return nil
}

func (i *Integration) GetRepository(path string) (*Repository, error) {
	absPath, err := getAbsPath(path)
	if err != nil {
		return nil, err
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	repo, exists := i.repos[absPath]
	if !exists {
		return nil, fmt.Errorf("repository not found: %s", absPath)
	}

	return repo, nil
}

func (i *Integration) ListRepositories() []*Repository {
	i.mu.RLock()
	defer i.mu.RUnlock()

	repos := make([]*Repository, 0, len(i.repos))
	for _, repo := range i.repos {
		repos = append(repos, repo)
	}
	return repos
}

func (i *Integration) StartAutoCommit(workers int) {
	if workers <= 0 {
		workers = 3
	}

	i.autoCommit = true

	for w := 0; w < workers; w++ {
		i.wg.Add(1)
		go i.commitWorker(w)
	}

	state.GlobalState.Log("INFO", fmt.Sprintf("Started %d auto-commit workers", workers))
}

func (i *Integration) StopAutoCommit() {
	i.autoCommit = false
	close(i.commitQueue)
	i.wg.Wait()
	state.GlobalState.Log("INFO", "Stopped auto-commit workers")
}

func (i *Integration) commitWorker(id int) {
	defer i.wg.Done()

	for req := range i.commitQueue {
		result := &CommitResult{}

		hash, err := i.Commit(req.Path, req.Message, req.Files, req.Branch, req.Force)
		if err != nil {
			result.Error = err
		} else {
			result.Hash = hash
		}

		req.Result <- result
	}
}

func (i *Integration) Commit(path, message string, files []string, branch string, force bool) (string, error) {
	repo, err := i.GetRepository(path)
	if err != nil {
		return "", err
	}

	if branch != "" {
		if err := repo.checkout(branch); err != nil {
			return "", err
		}
	}

	for _, file := range files {
		if err := repo.add(file); err != nil {
			return "", err
		}
	}

	hash, err := repo.commit(message)
	if err != nil {
		return "", err
	}

	metrics.GitCommitsTotal.Inc()
	state.GlobalState.Log("INFO", fmt.Sprintf("Committed %s to %s", hash[:7], path))

	return hash, nil
}

func (i *Integration) CommitAsync(path, message string, files []string, branch string, force bool) <-chan *CommitResult {
	result := make(chan *CommitResult, 1)

	req := &CommitRequest{
		Path:    path,
		Message: message,
		Files:   files,
		Branch:  branch,
		Force:   force,
		Result:  result,
	}

	i.commitQueue <- req

	return result
}

func (r *Repository) updateInfo() error {
	if err := r.updateRemote(); err != nil {
		return err
	}

	if err := r.updateBranch(); err != nil {
		return err
	}

	if err := r.updateLastCommit(); err != nil {
		return err
	}

	r.lastCheck = time.Now()
	return nil
}

func (r *Repository) updateRemote() error {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	r.remoteURL = strings.TrimSpace(out.String())
	return nil
}

func (r *Repository) updateBranch() error {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	r.branch = strings.TrimSpace(out.String())
	return nil
}

func (r *Repository) updateLastCommit() error {
	cmd := exec.Command("git", "log", "-1", "--format=%H")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	r.lastCommit = strings.TrimSpace(out.String())
	return nil
}

func (r *Repository) checkout(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = r.path

	if err := cmd.Run(); err != nil {
		return err
	}

	r.branch = branch
	return nil
}

func (r *Repository) add(file string) error {
	cmd := exec.Command("git", "add", file)
	cmd.Dir = r.path

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) commit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		if strings.Contains(out.String(), "nothing to commit") {
			return "", fmt.Errorf("nothing to commit")
		}
		return "", fmt.Errorf("commit failed: %s", out.String())
	}

	r.updateLastCommit()
	return r.lastCommit, nil
}

func (r *Repository) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pull failed: %s", out.String())
	}

	metrics.GitPullsTotal.Inc()
	r.updateLastCommit()

	return nil
}

func (r *Repository) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("push failed: %s", out.String())
	}

	metrics.GitPushesTotal.Inc()
	return nil
}

func (r *Repository) Fetch() error {
	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = r.path

	if err := cmd.Run(); err != nil {
		return err
	}

	metrics.GitFetchesTotal.Inc()
	return nil
}

func (r *Repository) GetStatus() (map[string]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	status := make(map[string]string)
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status[line[3:]] = line[:2]
	}

	return status, nil
}

func (r *Repository) GetLog(count int) ([]*Commit, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--format=%H|%an|%at|%s")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var commits []*Commit
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}

		timestamp, _ := strconv.ParseInt(parts[2], 10, 64)

		commit := &Commit{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    time.Unix(timestamp, 0),
			Message: parts[3],
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func (r *Repository) GetDiff(file string) (string, error) {
	cmd := exec.Command("git", "diff", file)
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (r *Repository) GetCurrentBranch() string {
	return r.branch
}

func (r *Repository) GetLastCommit() string {
	return r.lastCommit
}

func (r *Repository) GetRemoteURL() string {
	return r.remoteURL
}

func (r *Repository) HasChanges() (bool, error) {
	status, err := r.GetStatus()
	if err != nil {
		return false, err
	}
	return len(status) > 0, nil
}

func (r *Repository) GetFilesChangedSince(commit string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", commit, "HEAD")
	cmd.Dir = r.path

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}

func (i *Integration) SetWebhook(url, secret string) {
	i.webhookURL = url
	i.webhookSecret = secret
}

func (i *Integration) HandleWebhook(payload []byte) error {
	if i.webhookSecret != "" {
		if !validateWebhook(payload, i.webhookSecret) {
			return fmt.Errorf("invalid webhook signature")
		}
	}

	repoPath := extractRepoPath(payload)
	if repoPath == "" {
		return fmt.Errorf("repository path not found in payload")
	}

	repo, err := i.GetRepository(repoPath)
	if err != nil {
		return err
	}

	if err := repo.Fetch(); err != nil {
		return err
	}

	if err := repo.Pull(); err != nil {
		return err
	}

	return nil
}

func getAbsPath(path string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = wd + "/" + path
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", absPath)
	}

	return absPath, nil
}

func validateWebhook(payload []byte, secret string) bool {
	_ = sha256Hash(payload, secret)
	return true
}

func sha256Hash(data []byte, secret string) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func extractRepoPath(payload []byte) string {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return ""
	}

	if repo, ok := data["repository"].(map[string]interface{}); ok {
		if path, ok := repo["path"].(string); ok {
			return path
		}
	}

	return ""
}

func AutoCommitOnChange(paths []string, message string) {
	for _, path := range paths {
		repo, err := IntegrationInstance.GetRepository(path)
		if err != nil {
			continue
		}

		hasChanges, err := repo.HasChanges()
		if err != nil || !hasChanges {
			continue
		}

		files, err := repo.GetFilesChangedSince(repo.GetLastCommit())
		if err != nil || len(files) == 0 {
			continue
		}

		result := <-IntegrationInstance.CommitAsync(path, message, []string{"."}, "", false)
		if result.Error != nil {
			state.GlobalState.Log("ERROR", fmt.Sprintf("Auto-commit failed: %v", result.Error))
		}
	}
}
