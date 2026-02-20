const project = "my-prototype";

const agent = {
  name: "coder",
  model: "qwen-3.5-397b",
  instructions: `You are a helpful coding assistant for project: ${project}`,
};

const task = {
  prompt: "Say hello in 3 languages",
  agent: "coder",
};

console.log("OpenCode Minimal Template Ready");
console.log("Edit src/index.ts to customize");
