/**
 * Task definitions for OpenCode project
 * 
 * Define reusable tasks that agents can execute.
 */

export interface Task {
  id: string;
  name: string;
  description: string;
  agent: string;
  inputSchema: Record<string, unknown>;
  outputSchema: Record<string, unknown>;
}

export const tasks: Record<string, Task> = {
  generate_code: {
    id: 'generate_code',
    name: 'Generate Code',
    description: 'Generate code based on specification',
    agent: 'sisyphus',
    inputSchema: {
      specification: 'string',
      language: 'string',
    },
    outputSchema: {
      code: 'string',
      files: 'array',
    },
  },
  create_plan: {
    id: 'create_plan',
    name: 'Create Plan',
    description: 'Create a detailed project plan',
    agent: 'prometheus',
    inputSchema: {
      goal: 'string',
      constraints: 'array',
    },
    outputSchema: {
      plan: 'object',
      milestones: 'array',
    },
  },
  research: {
    id: 'research',
    name: 'Research',
    description: 'Research a topic thoroughly',
    agent: 'librarian',
    inputSchema: {
      query: 'string',
      depth: 'string',
    },
    outputSchema: {
      findings: 'array',
      sources: 'array',
    },
  },
};

export function getTask(id: string): Task | undefined {
  return tasks[id];
}

export function listTasks(): Task[] {
  return Object.values(tasks);
}
