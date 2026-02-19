class DashboardApp {
    constructor() {
        this.ws = null;
        this.reconnectInterval = 5000;
        this.charts = {};
        this.startTime = Date.now();
        this.init();
    }
    init() {
        this.setupWebSocket();
        this.createCharts();
        this.setupEventListeners();
        this.startUptimeTimer();
        this.fetchInitialData();
    }
    setupWebSocket() {
        const wsUrl = `ws://${window.location.host}/ws/dashboard`;
        this.ws = new WebSocket(wsUrl);
        this.ws.onopen = () => { console.log('Dashboard connected'); this.showAlert('Connected to BIOMETRICS', 'success'); };
        this.ws.onmessage = (event) => { const data = JSON.parse(event.data); this.handleUpdate(data); };
        this.ws.onclose = () => { console.log('Disconnected, reconnecting...'); setTimeout(() => this.setupWebSocket(), this.reconnectInterval); };
        this.ws.onerror = (error) => { console.error('WebSocket error:', error); this.showAlert('Connection error - retrying...', 'error'); };
    }
    createCharts() {
        const chartConfig = (type, data, options = {}) => ({ type, data, options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#a0a0b0' } }, y: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#a0a0b0' } } }, ...options } });
        this.charts.request = new Chart(document.getElementById('request-chart'), chartConfig('line', { labels: Array(10).fill(''), datasets: [{ data: Array(10).fill(0), borderColor: '#00f5ff', tension: 0.4 }] }));
        this.charts.response = new Chart(document.getElementById('response-chart'), chartConfig('line', { labels: Array(10).fill(''), datasets: [{ data: Array(10).fill(0), borderColor: '#ff00ff', tension: 0.4 }] }));
        this.charts.duration = new Chart(document.getElementById('duration-histogram'), chartConfig('bar', { labels: ['P50', 'P95', 'P99'], datasets: [{ data: [0, 0, 0], backgroundColor: ['#00f5ff', '#ff00ff', '#ff3366'] }] }));
    }
    handleUpdate(data) {
        if (data.type === 'metrics') { this.updateMetrics(data.payload); }
        else if (data.type === 'agents') { this.updateAgents(data.payload); }
        else if (data.type === 'alert') { this.showAlert(data.payload.message, data.payload.severity); }
    }
    updateMetrics(metrics) {
        document.getElementById('request-rate').textContent = metrics.requestRate || 0;
        document.getElementById('avg-response').textContent = Math.round(metrics.avgResponse || 0);
        document.getElementById('error-rate').textContent = (metrics.errorRate || 0).toFixed(2) + '%';
        document.getElementById('queue-size').textContent = metrics.queueSize || 0;
        this.charts.request.data.datasets[0].data.push(metrics.requestRate || 0);
        this.charts.request.data.datasets[0].data.shift();
        this.charts.request.update('none');
        this.charts.response.data.datasets[0].data.push(metrics.avgResponse || 0);
        this.charts.response.data.datasets[0].data.shift();
        this.charts.response.update('none');
    }
    updateAgents(agents) {
        const grid = document.getElementById('agents-grid');
        grid.innerHTML = '';
        document.getElementById('total-agents').textContent = agents.filter(a => a.status === 'active').length;
        agents.forEach(agent => { const card = this.createAgentCard(agent); grid.appendChild(card); });
    }
    createAgentCard(agent) {
        const template = document.getElementById('agent-card-template');
        const clone = template.content.cloneNode(true);
        const card = clone.querySelector('.agent-card');
        card.dataset.agentId = agent.id;
        card.querySelector('.agent-name').textContent = agent.name;
        card.querySelector('.agent-role').textContent = agent.role;
        card.querySelector('.agent-status-badge').className = `agent-status-badge ${agent.status}`;
        card.querySelector('.agent-status-badge').textContent = agent.status;
        const metrics = card.querySelectorAll('.agent-metric .metric-value');
        metrics[0].textContent = agent.tasksCompleted || 0;
        metrics[1].textContent = (agent.avgTime || 0) + 'ms';
        metrics[2].textContent = agent.errors || 0;
        const progressFill = card.querySelector('.progress-fill');
        progressFill.style.width = (agent.progress || 0) + '%';
        card.querySelector('.progress-text').textContent = (agent.progress || 0).toFixed(0) + '%';
        card.querySelector('.task-name').textContent = agent.currentTask || 'Idle';
        return card;
    }
    showAlert(message, severity = 'info') {
        const banner = document.getElementById('alert-banner');
        banner.textContent = message;
        banner.className = `alert-banner ${severity}`;
        setTimeout(() => { banner.classList.add('hidden'); }, 5000);
    }
    startUptimeTimer() {
        setInterval(() => {
            const elapsed = Date.now() - this.startTime;
            const hours = Math.floor(elapsed / 3600000);
            const minutes = Math.floor((elapsed % 3600000) / 60000);
            const seconds = Math.floor((elapsed % 60000) / 1000);
            document.getElementById('uptime').textContent = `${hours.toString().padStart(2,'0')}:${minutes.toString().padStart(2,'0')}:${seconds.toString().padStart(2,'0')}`;
            document.getElementById('last-update').textContent = new Date().toLocaleTimeString();
        }, 1000);
    }
    async fetchInitialData() {
        try {
            const response = await fetch('/api/dashboard/data');
            if (response.ok) { const data = await response.json(); this.handleUpdate({ type: 'metrics', payload: data.metrics }); this.handleUpdate({ type: 'agents', payload: data.agents }); }
        } catch (error) { console.log('Using mock data for demo'); this.loadMockData(); }
    }
    loadMockData() {
        const mockAgents = [
            { id: 'sisyphus', name: 'Sisyphus', role: 'Main Coder', status: 'active', tasksCompleted: 127, avgTime: 2340, errors: 2, progress: 75, currentTask: 'Refactoring API endpoints' },
            { id: 'prometheus', name: 'Prometheus', role: 'Planner', status: 'active', tasksCompleted: 89, avgTime: 1890, errors: 0, progress: 45, currentTask: 'Creating sprint plan' },
            { id: 'oracle', name: 'Oracle', role: 'Architect', status: 'idle', tasksCompleted: 56, avgTime: 3200, errors: 1, progress: 0, currentTask: 'Idle' },
            { id: 'librarian', name: 'Librarian', role: 'Documentation', status: 'active', tasksCompleted: 203, avgTime: 1560, errors: 0, progress: 90, currentTask: 'Writing API docs' },
        ];
        const mockMetrics = { requestRate: 42, avgResponse: 187, errorRate: 0.12, queueSize: 15 };
        this.updateAgents(mockAgents);
        this.updateMetrics(mockMetrics);
        setInterval(() => { mockMetrics.requestRate = Math.floor(Math.random() * 50) + 20; mockMetrics.avgResponse = Math.floor(Math.random() * 100) + 150; this.updateMetrics(mockMetrics); }, 3000);
    }
    setupEventListeners() {
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
            });
        });
    }
}
document.addEventListener('DOMContentLoaded', () => { window.dashboard = new DashboardApp(); });
