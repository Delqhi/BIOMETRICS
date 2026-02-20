/**
 * Agent definitions for OpenCode project
 * 
 * Define your AI agents here with their roles, models, and capabilities.
 */

export interface Agent {
  name: string;
  description: string;
  model: string;
  capabilities: string[];
  instructions?: string;
}

export const agents: Record<string, Agent> = {
  sisyphus: {
    name: 'Sisyphus',
    description: 'Main coder for implementation tasks',
    model: 'nvidia-nim/qwen-3.5-397b',
    capabilities: ['code', 'implementation', 'refactor'],
    instructions: 'Focus on clean, production-ready code with tests.',
  },
  prometheus: {
    name: 'Prometheus',
    description: 'Strategic planner for project tasks',
    model: 'nvidia-nim/qwen-3.5-397b',
    capabilities: ['planning', 'architecture', 'analysis'],
    instructions: 'Create detailed plans with clear milestones.',
  },
  librarian: {
    name: 'Librarian',
    description: 'Documentation and research agent',
    model: 'opencode-zen/zen-big-pickle',
    capabilities: ['research', 'documentation', 'search'],
    instructions: 'Thorough research with proper citations.',
  },
};

export function getAgent(name: string): Agent | undefined {
  return agents[name];
}

export function listAgents(): Agent[] {
  return Object.values(agents);
}
