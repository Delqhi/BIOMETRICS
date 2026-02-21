#!/bin/bash
# BIOMETRICS ORCHESTRATOR - LIVE DASHBOARD
# Shows real-time metrics from Go orchestrator

clear
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "   BIOMETRICS ENTERPRISE ORCHESTRATOR - LIVE STATUS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if orchestrator is running
PID=$(pgrep -f "./biometrics" | head -1)
if [ -z "$PID" ]; then
    echo "❌ ERROR: Go orchestrator NOT running!"
    echo ""
    echo "Start it with:"
    echo "  cd /Users/jeremy/dev/BIOMETRICS/biometrics-cli"
    echo "  ./biometrics"
    exit 1
fi

echo "✅ STATUS: RUNNING (PID: $PID)"
echo ""

# Get metrics
CYCLES=$(curl -s http://localhost:59002/metrics | grep "^biometrics_orchestrator_cycles_total " | awk '{print $2}')
MODEL_ACQ=$(curl -s http://localhost:59002/metrics | grep "^biometrics_orchestrator_model_acquisitions_total{model=\"qwen3.5\"}" | awk '{print $2}')
GOROUTINES=$(curl -s http://localhost:59002/metrics | grep "^go_goroutines " | awk '{print $2}')
CHAOS_EVENTS=$(curl -s http://localhost:59002/metrics | grep "^biometrics_orchestrator_chaos_events_total" | awk -F'} ' '{print $2}' | awk '{sum+=$1} END {print sum}')
UPTIME=$(ps -o etime= -p $PID)

echo "📊 METRICS:"
echo "  ┌─────────────────────────────────────┐"
echo "  │ Cycles Completed:      $(printf "%6s" "$CYCLES")          │"
echo "  │ Model Acquisitions:    $(printf "%6s" "$MODEL_ACQ")          │"
echo "  │ Active Goroutines:     $(printf "%6s" "$GOROUTINES")          │"
echo "  │ Chaos Events:          $(printf "%6s" "$CHAOS_EVENTS")          │"
echo "  │ Uptime:                $UPTIME │"
echo "  └─────────────────────────────────────┘"
echo ""

echo "🔗 ENDPOINTS:"
echo "  • Metrics: http://localhost:59002/metrics"
echo "  • Process: PID $PID"
echo ""

echo "📈 RECENT LOGS:"
tail -5 /tmp/biometrics-orchestrator.log 2>/dev/null | while read line; do
    echo "  $line"
done
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Press Ctrl+C to exit"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
