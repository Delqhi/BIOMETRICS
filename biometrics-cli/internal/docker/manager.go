package docker

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

type Manager struct {
	mu           sync.RWMutex
	containers   map[string]*Container
	networks     map[string]*Network
	images       map[string]*Image
	pollInterval time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

type Container struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Ports   []PortMapping
	Created time.Time
}

type PortMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

type Network struct {
	Name    string
	Driver  string
	Subnet  string
	Gateway string
}

type Image struct {
	Name    string
	Tag     string
	Size    int64
	Created time.Time
}

var defaultManager = &Manager{
	containers:   make(map[string]*Container),
	networks:     make(map[string]*Network),
	images:       make(map[string]*Image),
	pollInterval: 30 * time.Second,
	stopChan:     make(chan struct{}),
}

var ManagerInstance = defaultManager

func (m *Manager) StartPolling() {
	m.wg.Add(1)
	go m.pollContainers()
	state.GlobalState.Log("INFO", "Started container polling")
}

func (m *Manager) StopPolling() {
	close(m.stopChan)
	m.wg.Wait()
	state.GlobalState.Log("INFO", "Stopped container polling")
}

func (m *Manager) pollContainers() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			if err := m.RefreshContainers(); err != nil {
				state.GlobalState.Log("ERROR", fmt.Sprintf("Failed to refresh containers: %v", err))
			}
		}
	}
}

func (m *Manager) RefreshContainers() error {
	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}|{{.Ports}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.containers = make(map[string]*Container)

	lines := bytes.Split(out.Bytes(), []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := bytes.Split(line, []byte("|"))
		if len(parts) < 4 {
			continue
		}

		container := &Container{
			ID:     string(parts[0]),
			Name:   string(parts[1]),
			Image:  string(parts[2]),
			Status: string(parts[3]),
		}

		m.containers[container.ID] = container
		metrics.DockerContainersRunning.Set(float64(len(m.containers)))
	}

	return nil
}

func (m *Manager) GetContainer(nameOrID string) (*Container, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if c, ok := m.containers[nameOrID]; ok {
		return c, nil
	}

	for _, c := range m.containers {
		if c.Name == nameOrID {
			return c, nil
		}
	}

	return nil, fmt.Errorf("container %s not found", nameOrID)
}

func (m *Manager) ListContainers() []*Container {
	m.mu.RLock()
	defer m.mu.RUnlock()

	containers := make([]*Container, 0, len(m.containers))
	for _, c := range m.containers {
		containers = append(containers, c)
	}
	return containers
}

func (m *Manager) StartContainer(nameOrID string) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Starting container: %s", nameOrID))
	cmd := exec.Command("docker", "start", nameOrID)
	if err := cmd.Run(); err != nil {
		metrics.DockerContainerStartsFailedTotal.Inc()
		return err
	}
	metrics.DockerContainerStartsTotal.Inc()
	return m.RefreshContainers()
}

func (m *Manager) StopContainer(nameOrID string) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Stopping container: %s", nameOrID))
	cmd := exec.Command("docker", "stop", nameOrID)
	if err := cmd.Run(); err != nil {
		metrics.DockerContainerStopsFailedTotal.Inc()
		return err
	}
	metrics.DockerContainerStopsTotal.Inc()
	return m.RefreshContainers()
}

func (m *Manager) RestartContainer(nameOrID string) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Restarting container: %s", nameOrID))
	cmd := exec.Command("docker", "restart", nameOrID)
	if err := cmd.Run(); err != nil {
		return err
	}
	metrics.DockerContainerRestartsTotal.Inc()
	return m.RefreshContainers()
}

func (m *Manager) RemoveContainer(nameOrID string, force bool) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Removing container: %s", nameOrID))
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, nameOrID)
	cmd := exec.Command("docker", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	metrics.DockerContainersRemovedTotal.Inc()
	return m.RefreshContainers()
}

func (m *Manager) GetContainerLogs(nameOrID string, tail int) (string, error) {
	args := []string{"logs", "--tail", fmt.Sprintf("%d", tail), nameOrID}
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (m *Manager) InspectContainer(nameOrID string) (map[string]interface{}, error) {
	cmd := exec.Command("docker", "inspect", nameOrID)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("container not found")
	}

	return result[0], nil
}

func (m *Manager) GetContainerStats(nameOrID string) (map[string]interface{}, error) {
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.CPUPerc}}|{{.MemUsage}}|{{.NetIO}}|{{.BlockIO}}", nameOrID)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	parts := bytes.Split(bytes.TrimSpace(out.Bytes()), []byte("|"))
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid stats output")
	}

	return map[string]interface{}{
		"cpu_percent": string(parts[0]),
		"memory":      string(parts[1]),
		"net_io":      string(parts[2]),
		"block_io":    string(parts[3]),
	}, nil
}

func (m *Manager) PullImage(imageName string) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Pulling image: %s", imageName))
	cmd := exec.Command("docker", "pull", imageName)
	if err := cmd.Run(); err != nil {
		metrics.DockerImagePullsFailedTotal.Inc()
		return err
	}
	metrics.DockerImagePullsTotal.Inc()
	return nil
}

func (m *Manager) ListImages() ([]*Image, error) {
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}|{{.Tag}}|{{.Size}}|{{.CreatedAt}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var images []*Image
	lines := bytes.Split(out.Bytes(), []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := bytes.Split(line, []byte("|"))
		if len(parts) < 4 {
			continue
		}

		img := &Image{
			Name: string(parts[0]),
			Tag:  string(parts[1]),
		}
		images = append(images, img)
	}

	return images, nil
}

func (m *Manager) CreateNetwork(name, driver, subnet, gateway string) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Creating network: %s", name))
	args := []string{"network", "create", "--driver", driver}
	if subnet != "" {
		args = append(args, "--subnet", subnet)
	}
	if gateway != "" {
		args = append(args, "--gateway", gateway)
	}
	args = append(args, name)

	cmd := exec.Command("docker", args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	m.mu.Lock()
	m.networks[name] = &Network{
		Name:    name,
		Driver:  driver,
		Subnet:  subnet,
		Gateway: gateway,
	}
	m.mu.Unlock()

	metrics.DockerNetworksCreatedTotal.Inc()
	return nil
}

func (m *Manager) RemoveNetwork(name string) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Removing network: %s", name))
	cmd := exec.Command("docker", "network", "rm", name)
	if err := cmd.Run(); err != nil {
		return err
	}

	m.mu.Lock()
	delete(m.networks, name)
	m.mu.Unlock()

	metrics.DockerNetworksRemovedTotal.Inc()
	return nil
}

func (m *Manager) ConnectContainerToNetwork(containerName, networkName string) error {
	cmd := exec.Command("docker", "network", "connect", networkName, containerName)
	if err := cmd.Run(); err != nil {
		return err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Connected %s to %s", containerName, networkName))
	return nil
}

func (m *Manager) DisconnectContainerFromNetwork(containerName, networkName string) error {
	cmd := exec.Command("docker", "network", "disconnect", networkName, containerName)
	if err := cmd.Run(); err != nil {
		return err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Disconnected %s from %s", containerName, networkName))
	return nil
}

func (m *Manager) HealthCheck() error {
	if err := m.RefreshContainers(); err != nil {
		return fmt.Errorf("docker daemon not accessible: %w", err)
	}

	healthy := 0
	unhealthy := 0

	for _, c := range m.containers {
		if c.Status == "running" {
			healthy++
		} else {
			unhealthy++
		}
	}

	state.GlobalState.Log("INFO", fmt.Sprintf("Docker health: %d healthy, %d unhealthy", healthy, unhealthy))

	if healthy == 0 && len(m.containers) > 0 {
		return fmt.Errorf("no healthy containers found")
	}

	return nil
}

func (m *Manager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"containers": len(m.containers),
		"networks":   len(m.networks),
		"images":     len(m.images),
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := ManagerInstance.HealthCheck(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func SetupHealthCheck() {
	http.HandleFunc("/health/docker", HealthCheckHandler)
}

func Start() {
	ManagerInstance.RefreshContainers()
	ManagerInstance.StartPolling()
}

func Stop() {
	ManagerInstance.StopPolling()
}
