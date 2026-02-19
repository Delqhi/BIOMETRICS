#!/bin/bash
set -euo pipefail

# ============================================================================
# BIOMETRICS BLUE-GREEN DEPLOYMENT SCRIPT
# ============================================================================
# Version: 1.0.0
# Description: Zero-downtime blue-green deployment for BIOMETRICS
# Usage: ./blue-green-deploy.sh [deploy|switch|rollback|status]
# ============================================================================

# Configuration
NAMESPACE="${NAMESPACE:-biometrics}"
RELEASE_NAME="${RELEASE_NAME:-biometrics}"
CHART_PATH="${CHART_PATH:-./helm/biometrics}"
TIMEOUT="${TIMEOUT:-300s}"
HELM="${HELM:-helm}"
KUBECTL="${KUBECTL:-kubectl}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Helper functions
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v $HELM &> /dev/null; then
        log_error "Helm is not installed. Please install Helm first."
        exit 1
    fi
    
    if ! command -v $KUBECTL &> /dev/null; then
        log_error "kubectl is not installed. Please install kubectl first."
        exit 1
    fi
    
    if ! $KUBECTL cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

get_current_deployment() {
    local current=$($KUBECTL get deployments -n $NAMESPACE -l app=biometrics-opencode-server -o jsonpath='{.items[0].metadata.labels.deployment-type}' 2>/dev/null || echo "")
    echo "$current"
}

get_active_color() {
    local current=$(get_current_deployment)
    if [ "$current" == "blue" ]; then
        echo "blue"
    elif [ "$current" == "green" ]; then
        echo "green"
    else
        echo "none"
    fi
}

get_inactive_color() {
    local active=$(get_active_color)
    if [ "$active" == "blue" ]; then
        echo "green"
    elif [ "$active" == "green" ]; then
        echo "blue"
    else
        echo "blue"
    fi
}

# Deployment functions
deploy() {
    local color=$(get_inactive_color)
    local target_release="${RELEASE_NAME}-${color}"
    
    log_info "Starting blue-green deployment..."
    log_info "Current active: $(get_active_color)"
    log_info "Deploying to: $color"
    
    # Create namespace if not exists
    $KUBECTL create namespace $NAMESPACE --dry-run=client -o yaml | $KUBECTL apply -f -
    
    # Deploy to inactive color
    log_info "Installing chart to ${target_release}..."
    $HELM upgrade --install ${target_release} ${CHART_PATH} \
        --namespace ${NAMESPACE} \
        --set deploymentType=${color} \
        --set replicaCount=1 \
        --timeout ${TIMEOUT} \
        --wait \
        --wait-for-jobs
    
    log_success "Deployment to ${color} completed"
    
    # Run health checks
    log_info "Running health checks..."
    if ! run_health_checks ${color}; then
        log_error "Health checks failed for ${color} deployment"
        log_info "Cleaning up failed deployment..."
        $HELM uninstall ${target_release} -n ${NAMESPACE}
        exit 1
    fi
    
    log_success "Health checks passed for ${color}"
    
    # Scale up new deployment
    log_info "Scaling up ${color} deployment..."
    $HELM upgrade ${target_release} ${CHART_PATH} \
        --namespace ${NAMESPACE} \
        --set deploymentType=${color} \
        --reuse-values \
        --wait
    
    log_success "Blue-green deployment completed successfully!"
    log_info "Active deployment: ${color}"
    log_info "To switch traffic, run: $0 switch"
}

switch() {
    local current=$(get_active_color)
    local target=$(get_inactive_color)
    local target_release="${RELEASE_NAME}-${target}"
    
    if [ "$current" == "none" ]; then
        log_error "No active deployment found. Please deploy first."
        exit 1
    fi
    
    log_info "Switching traffic from ${current} to ${target}..."
    
    # Verify target deployment exists
    if ! $HELM status ${target_release} -n ${NAMESPACE} &> /dev/null; then
        log_error "Target deployment ${target} does not exist. Please deploy first."
        exit 1
    fi
    
    # Update ingress to point to new deployment
    log_info "Updating ingress..."
    $KUBECTL patch ingress ${RELEASE_NAME}-ingress -n ${NAMESPACE} \
        --type='json' \
        -p="[{\"op\": \"replace\", \"path\": \"/spec/rules/0/http/paths/0/backend/service/name\", \"value\": \"${RELEASE_NAME}-opencode-server-${target}\"}]"
    
    # Wait for ingress to update
    sleep 5
    
    # Verify switch
    log_info "Verifying traffic switch..."
    if verify_switch ${target}; then
        log_success "Traffic successfully switched to ${target}"
        
        # Scale down old deployment
        log_info "Scaling down old deployment (${current})..."
        $HELM upgrade ${RELEASE_NAME}-${current} ${CHART_PATH} \
            --namespace ${NAMESPACE} \
            --set replicaCount=0 \
            --reuse-values
        
        log_success "Blue-green switch completed!"
        log_info "Active deployment: ${target}"
    else
        log_error "Traffic switch verification failed. Rolling back..."
        rollback
        exit 1
    fi
}

rollback() {
    local current=$(get_active_color)
    local target=$(get_inactive_color)
    
    log_warning "Rolling back from ${current} to ${target}..."
    
    # Switch ingress back
    $KUBECTL patch ingress ${RELEASE_NAME}-ingress -n ${NAMESPACE} \
        --type='json' \
        -p="[{\"op\": \"replace\", \"path\": \"/spec/rules/0/http/paths/0/backend/service/name\", \"value\": \"${RELEASE_NAME}-opencode-server-${target}\"}]"
    
    log_success "Rollback completed"
}

status() {
    local active=$(get_active_color)
    
    echo ""
    echo "=========================================="
    echo "  BIOMETRICS DEPLOYMENT STATUS"
    echo "=========================================="
    echo ""
    echo "Active Deployment: ${active}"
    echo ""
    echo "Helm Releases:"
    $HELM list -n ${NAMESPACE} --all
    echo ""
    echo "Pods:"
    $KUBECTL get pods -n ${NAMESPACE} -l app=biometrics-opencode-server
    echo ""
    echo "Services:"
    $KUBECTL get svc -n ${NAMESPACE} -l app=biometrics-opencode-server
    echo ""
    echo "Ingress:"
    $KUBECTL get ingress -n ${NAMESPACE}
    echo ""
    echo "=========================================="
}

# Health check functions
run_health_checks() {
    local color=$1
    local max_attempts=30
    local attempt=0
    
    log_info "Running health checks for ${color}..."
    
    while [ $attempt -lt $max_attempts ]; do
        attempt=$((attempt + 1))
        log_info "Health check attempt ${attempt}/${max_attempts}..."
        
        # Check if pods are ready
        local ready_pods=$($KUBECTL get pods -n ${NAMESPACE} -l app=biometrics-opencode-server,deployment-type=${color} \
            -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}' | grep -c "True" || echo "0")
        
        if [ "$ready_pods" -gt 0 ]; then
            log_success "Pods are ready"
            
            # Check HTTP health endpoint
            local pod_name=$($KUBECTL get pods -n ${NAMESPACE} -l app=biometrics-opencode-server,deployment-type=${color} \
                -o jsonpath='{.items[0].metadata.name}')
            
            local health_status=$($KUBECTL exec ${pod_name} -n ${NAMESPACE} -- \
                curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/global/health || echo "000")
            
            if [ "$health_status" == "200" ]; then
                log_success "Health endpoint returned 200"
                return 0
            else
                log_warning "Health endpoint returned ${health_status}"
            fi
        fi
        
        sleep 10
    done
    
    log_error "Health checks failed after ${max_attempts} attempts"
    return 1
}

verify_switch() {
    local target=$1
    local max_attempts=10
    local attempt=0
    
    log_info "Verifying traffic switch to ${target}..."
    
    while [ $attempt -lt $max_attempts ]; do
        attempt=$((attempt + 1))
        
        # Get ingress IP
        local ingress_ip=$($KUBECTL get ingress ${RELEASE_NAME}-ingress -n ${NAMESPACE} \
            -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
        
        if [ -n "$ingress_ip" ]; then
            local response=$(curl -s -o /dev/null -w "%{http_code}" http://${ingress_ip}/global/health || echo "000")
            
            if [ "$response" == "200" ]; then
                return 0
            fi
        fi
        
        sleep 5
    done
    
    return 1
}

# Main script
main() {
    local command=${1:-help}
    
    case $command in
        deploy)
            check_prerequisites
            deploy
            ;;
        switch)
            check_prerequisites
            switch
            ;;
        rollback)
            check_prerequisites
            rollback
            ;;
        status)
            status
            ;;
        help|*)
            echo "BIOMETRICS Blue-Green Deployment Script"
            echo ""
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  deploy    Deploy to inactive color (blue/green)"
            echo "  switch    Switch traffic to new deployment"
            echo "  rollback  Rollback to previous deployment"
            echo "  status    Show deployment status"
            echo "  help      Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  NAMESPACE       Kubernetes namespace (default: biometrics)"
            echo "  RELEASE_NAME    Helm release name (default: biometrics)"
            echo "  CHART_PATH    Path to Helm chart (default: ./helm/biometrics)"
            echo "  TIMEOUT       Deployment timeout (default: 300s)"
            ;;
    esac
}

main "$@"
