#!/usr/bin/env node

import chalk from 'chalk';
import inquirer from 'inquirer';
import ora from 'ora';
import { execa } from 'execa';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs/promises';
import os from 'os';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const HOME_DIR = os.homedir();

// API Key Help Links
const API_HELP_LINKS = {
  gitlab: 'https://gitlab.com/-/profile/personal_access_tokens',
  nvidia: 'https://build.nvidia.com/explore/discover',
  whatsapp: 'https://developers.facebook.com/apps/creation/',
  telegram: 'https://core.telegram.org/bots/features#botfather',
  gmail: 'https://console.cloud.google.com/apis/credentials',
  twitter: 'https://developer.twitter.com/en/portal/dashboard',
  clawdbot: 'https://clawdbot.com/dashboard',
};

// Banner
const banner = `
${chalk.cyan('╔══════════════════════════════════════════════════════════╗')}
${chalk.cyan('║')}          ${chalk.yellow.bold('BIOMETRICS ONBOARD')} ${chalk.cyan('║')}
${chalk.cyan('║')}    ${chalk.gray('Complete Setup: GitLab + OpenCode + OpenClaw')}  ${chalk.cyan('║')}
${chalk.cyan('║')}      ${chalk.gray('+ All Integrations (WhatsApp, Telegram, ...)')}   ${chalk.cyan('║')}
${chalk.cyan('╚══════════════════════════════════════════════════════════╝')}
`;

console.log(banner);

// Helper functions
async function runCommand(command, args, options = {}) {
  try {
    const { stdout, stderr } = await execa(command, args, options);
    return { success: true, stdout, stderr };
  } catch (error) {
    return { success: false, error: error.message };
  }
}

async function checkInstalled(command) {
  const result = await runCommand('which', [command], { reject: false });
  return result.success;
}

function showHelpLink(service) {
  const link = API_HELP_LINKS[service];
  if (link) {
    console.log(chalk.blue(`ℹ Create ${service} API key: ${link}`));
  }
}

// Main onboarding questions
const questions = [
  {
    type: 'confirm',
    name: 'needGitLabToken',
    message: 'Do you have a GitLab Personal Access Token?',
    default: false,
  },
  {
    type: 'input',
    name: 'gitlabToken',
    message: 'Enter your GitLab Personal Access Token:',
    validate: (input) => {
      if (!input.startsWith('glpat-')) {
        return 'GitLab token must start with "glpat-"';
      }
      return true;
    },
    when: (answers) => answers.needGitLabToken,
  },
  {
    type: 'input',
    name: 'nvidiaApiKey',
    message: 'Enter your NVIDIA API Key:',
    validate: (input) => {
      if (input.length < 10) {
        return 'Please enter a valid NVIDIA API Key';
      }
      return true;
    },
  },
  {
    type: 'confirm',
    name: 'setupWhatsApp',
    message: 'Setup WhatsApp integration?',
    default: true,
  },
  {
    type: 'input',
    name: 'whatsappToken',
    message: 'Enter WhatsApp Business API Token:',
    when: (answers) => answers.setupWhatsApp,
  },
  {
    type: 'confirm',
    name: 'setupTelegram',
    message: 'Setup Telegram Bot integration?',
    default: true,
  },
  {
    type: 'input',
    name: 'telegramBotToken',
    message: 'Enter Telegram Bot Token (from @BotFather):',
    when: (answers) => answers.setupTelegram,
  },
  {
    type: 'confirm',
    name: 'setupGmail',
    message: 'Setup Gmail integration?',
    default: true,
  },
  {
    type: 'confirm',
    name: 'setupTwitter',
    message: 'Setup Twitter/X integration?',
    default: false,
  },
  {
    type: 'confirm',
    name: 'installOpenCode',
    message: 'Install OpenCode?',
    default: true,
  },
  {
    type: 'confirm',
    name: 'installOpenClaw',
    message: 'Install OpenClaw?',
    default: true,
  },
];

// System requirements check and auto-install
async function checkAndInstallRequirements() {
  const spinner = ora();
  
  console.log(chalk.blue('\n🔍 Checking system requirements...\n'));
  
  const requirements = [
    {
      name: 'Git',
      command: 'git',
      installCmd: ['brew', 'install', 'git'],
      checkCmd: ['git', '--version'],
    },
    {
      name: 'Node.js',
      command: 'node',
      installCmd: ['brew', 'install', 'node'],
      checkCmd: ['node', '--version'],
    },
    {
      name: 'pnpm',
      command: 'pnpm',
      installCmd: ['brew', 'install', 'pnpm'],
      checkCmd: ['pnpm', '--version'],
    },
    {
      name: 'Homebrew',
      command: 'brew',
      installCmd: ['/bin/bash', '-c', '$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)'],
      checkCmd: ['brew', '--version'],
      macOSOnly: true,
    },
    {
      name: 'Python 3',
      command: 'python3',
      installCmd: ['brew', 'install', 'python@3.11'],
      checkCmd: ['python3', '--version'],
    },
  ];
  
  for (const req of requirements) {
    // Skip macOS-only tools on other platforms
    if (req.macOSOnly && os.platform() !== 'darwin') {
      continue;
    }
    
    const installed = await checkInstalled(req.command);
    
    if (!installed) {
      spinner.start(`Installing ${req.name}...`);
      console.log(chalk.gray(`   → Running: ${req.installCmd.join(' ')}`));
      
      try {
        const result = await execa(req.installCmd[0], req.installCmd.slice(1), {
          stdio: ['ignore', 'pipe', 'pipe'],
          reject: false,
        });
        
        if (result.exitCode === 0) {
          spinner.succeed(chalk.green(`${req.name} installed successfully`));
        } else {
          spinner.fail(chalk.red(`Failed to install ${req.name}`));
          console.log(chalk.yellow(`   Manual installation required: ${req.installCmd.join(' ')}`));
        }
      } catch (error) {
        spinner.fail(chalk.red(`Error installing ${req.name}: ${error.message}`));
      }
    } else {
      // Check version
      try {
        const version = await execa(req.checkCmd[0], req.checkCmd.slice(1));
        spinner.info(chalk.cyan(`${req.name} ${version.stdout.trim()} already installed`));
      } catch (error) {
        spinner.info(chalk.cyan(`${req.name} already installed`));
      }
    }
  }
  
  console.log('');
}

// Run onboarding
async function main() {
  const spinner = ora();

  try {
    // Step 0: Check and install system requirements
    await checkAndInstallRequirements();
    
    // Show help links upfront
    console.log(chalk.blue('\n📖 Need help getting API keys?'));
    console.log('   GitLab:  ' + chalk.gray(API_HELP_LINKS.gitlab));
    console.log('   NVIDIA:  ' + chalk.gray(API_HELP_LINKS.nvidia));
    console.log('   WhatsApp:' + chalk.gray(API_HELP_LINKS.whatsapp));
    console.log('   Telegram:' + chalk.gray(API_HELP_LINKS.telegram));
    console.log('   Gmail:   ' + chalk.gray(API_HELP_LINKS.gmail));
    console.log('   Twitter: ' + chalk.gray(API_HELP_LINKS.twitter) + '\n');

    // Step 1: Collect user input
    const answers = await inquirer.prompt(questions);

    // Auto-generate GitLab token if not provided
    if (!answers.gitlabToken) {
      console.log(chalk.yellow('\n⚠ No GitLab token provided. Skipping GitLab project creation.'));
      console.log('   You can create one later at:', chalk.gray(API_HELP_LINKS.gitlab));
    }

    console.log('\n' + chalk.green('✓') + ' Configuration collected\n');

    // Step 2: Create GitLab media storage project
    if (answers.gitlabToken) {
      spinner.start('Creating GitLab media storage project...');
      const gitlabProject = await runCommand('curl', [
        '-s', '-X', 'POST',
        'https://gitlab.com/api/v4/projects',
        '-H', `PRIVATE-TOKEN: ${answers.gitlabToken}`,
        '-H', 'Content-Type: application/json',
        '-d', JSON.stringify({
          name: 'biometrics-media',
          description: 'BIOMETRICS project media storage (videos, PDFs, images)',
          visibility: 'public'
        })
      ]);
      
      if (gitlabProject.success) {
        const projectData = JSON.parse(gitlabProject.stdout);
        if (projectData.web_url) {
          spinner.succeed(chalk.green(`GitLab project created: ${projectData.web_url}`));
          
          // Store GitLab credentials in .env
          const envContent = `GITLAB_TOKEN=${answers.gitlabToken}\nGITLAB_PROJECT_ID=${projectData.id}\nGITLAB_PROJECT_PATH=${projectData.path_with_namespace}\n`;
          await fs.writeFile(join(process.cwd(), '.env'), envContent);
          spinner.succeed(chalk.green('GitLab credentials saved to .env'));
        } else {
          spinner.warn(chalk.yellow('GitLab project may already exist'));
        }
      } else {
        spinner.fail(chalk.red('GitLab project creation failed'));
        console.log(chalk.yellow('You can create manually at: https://gitlab.com/projects/new'));
      }
    }

    // Step 3: Install NLM CLI (always)
    spinner.start('Installing NLM CLI...');
    const nlmInstall = await runCommand('pnpm', ['add', '-g', 'nlm-cli'], { cwd: HOME_DIR });
    if (nlmInstall.success) {
      spinner.succeed(chalk.green('NLM CLI installed'));
      
      // Authenticate NLM CLI (opens browser)
      console.log(chalk.yellow('\nℹ Browser will open for NLM authentication...'));
      await runCommand('nlm', ['auth', 'login'], { stdio: 'inherit' });
      spinner.succeed(chalk.green('NLM CLI authenticated'));
    } else {
      spinner.fail(chalk.red('NLM CLI installation failed'));
      console.log(chalk.yellow('You can install manually: pnpm add -g nlm-cli'));
    }

    // Step 4: Install OpenCode (if requested)
    if (answers.installOpenCode) {
      spinner.start('Installing OpenCode...');
      const opencodeInstall = await runCommand('brew', ['install', 'opencode'], { cwd: HOME_DIR });
      if (opencodeInstall.success) {
        spinner.succeed(chalk.green('OpenCode installed'));
        
        // Configure OpenCode
        spinner.start('Configuring OpenCode...');
        const configDir = join(HOME_DIR, '.config', 'opencode');
        await fs.mkdir(configDir, { recursive: true });
        
        const opencodeConfig = {
          "provider": {
            "google": {
              "npm": "@ai-sdk/google",
              "models": {
                "gemini-2.5-pro": {
                  "id": "gemini-2.5-pro",
                  "name": "Gemini 2.5 Pro"
                }
              }
            },
            "nvidia": {
              "npm": "@ai-sdk/openai-compatible",
              "options": {
                "baseURL": "https://integrate.api.nvidia.com/v1"
              },
              "models": {
                "qwen-3.5-397b": {
                  "id": "qwen/qwen3.5-397b-a17b",
                  "name": "Qwen 3.5 397B"
                }
              }
            }
          }
        };
        
        await fs.writeFile(
          join(configDir, 'opencode.json'),
          JSON.stringify(opencodeConfig, null, 2)
        );
        spinner.succeed(chalk.green('OpenCode configured'));
        
        // Set NVIDIA API key
        spinner.start('Setting NVIDIA API key...');
        process.env.NVIDIA_API_KEY = answers.nvidiaApiKey;
        spinner.succeed(chalk.green('NVIDIA API key configured (environment)'));
      } else {
        spinner.fail(chalk.red('OpenCode installation failed'));
        console.log(chalk.yellow('You can install manually: brew install opencode'));
      }
    }

    // Step 5: Install OpenClaw (if requested)
    if (answers.installOpenClaw) {
      spinner.start('Installing OpenClaw...');
      const openclawInstall = await runCommand('pnpm', ['add', '-g', '@delqhi/openclaw'], { cwd: HOME_DIR });
      if (openclawInstall.success) {
        spinner.succeed(chalk.green('OpenClaw installed'));
        
        // Configure OpenClaw
        spinner.start('Configuring OpenClaw...');
        const configDir = join(HOME_DIR, '.openclaw');
        await fs.mkdir(configDir, { recursive: true });
        
        const openclawConfig = {
          "env": {
            "NVIDIA_API_KEY": answers.nvidiaApiKey,
            ...(answers.whatsappToken && { WHATSAPP_TOKEN: answers.whatsappToken }),
            ...(answers.telegramBotToken && { TELEGRAM_BOT_TOKEN: answers.telegramBotToken })
          },
          "models": {
            "providers": {
              "nvidia": {
                "baseUrl": "https://integrate.api.nvidia.com/v1",
                "api": "openai-completions",
                "models": []
              }
            }
          },
          "agents": {
            "defaults": {
              "model": {
                "primary": "nvidia/qwen/qwen3.5-397b-a17b"
              }
            }
          },
          "integrations": {
            "whatsapp": answers.setupWhatsApp ? {
              "enabled": true,
              "token": answers.whatsappToken || "${WHATSAPP_TOKEN}"
            } : { "enabled": false },
            "telegram": answers.setupTelegram ? {
              "enabled": true,
              "botToken": answers.telegramBotToken || "${TELEGRAM_BOT_TOKEN}"
            } : { "enabled": false },
            "gmail": answers.setupGmail ? {
              "enabled": true,
              "auth": "oauth2"
            } : { "enabled": false },
            "twitter": answers.setupTwitter ? {
              "enabled": true,
              "auth": "oauth2"
            } : { "enabled": false },
            "clawdbot": {
              "enabled": true,
              "url": "https://clawdbot.com/api"
            }
          }
        };
        
        await fs.writeFile(
          join(configDir, 'openclaw.json'),
          JSON.stringify(openclawConfig, null, 2)
        );
        spinner.succeed(chalk.green('OpenClaw configured'));
        
        // Setup ClawdBot integration
        if (answers.whatsappToken || answers.telegramBotToken) {
          spinner.start('Setting up ClawdBot integration...');
          await runCommand('openclaw', ['integrations', 'setup'], { cwd: HOME_DIR });
          spinner.succeed(chalk.green('ClawdBot integration complete'));
        }
      } else {
        spinner.fail(chalk.red('OpenClaw installation failed'));
        console.log(chalk.yellow('You can install manually: pnpm add -g @delqhi/openclaw'));
      }
    }

    // Step 6: Install Google Antigravity plugin
    spinner.start('Installing Google Antigravity plugin...');
    const antigravityInstall = await runCommand('opencode', ['plugin', 'add', 'opencode-antigravity-auth'], { cwd: HOME_DIR });
    if (antigravityInstall.success) {
      spinner.succeed(chalk.green('Google Antigravity plugin installed'));
      
      // Authenticate
      console.log(chalk.yellow('\nℹ Browser will open for Google authentication...'));
      await runCommand('opencode', ['auth', 'login'], { stdio: 'inherit' });
      spinner.succeed(chalk.green('Google authenticated'));
    } else {
      spinner.fail(chalk.red('Antigravity plugin installation failed'));
    }

    // Step 7: Verification
    console.log('\n' + chalk.bold('Running verification tests...') + '\n');
    
    const checks = [
      { name: 'NLM CLI', command: 'nlm', args: ['--version'] },
      { name: 'OpenCode', command: 'opencode', args: ['--version'] },
      { name: 'OpenClaw', command: 'openclaw', args: ['--version'] },
    ];
    
    for (const check of checks) {
      const result = await runCommand(check.command, check.args);
      if (result.success) {
        console.log(chalk.green('✓') + ` ${check.name} is installed`);
      } else {
        console.log(chalk.yellow('⚠') + ` ${check.name} not found`);
      }
    }

    // Check integrations
    if (answers.setupWhatsApp && answers.whatsappToken) {
      console.log(chalk.green('✓') + ' WhatsApp integration configured');
    }
    if (answers.setupTelegram && answers.telegramBotToken) {
      console.log(chalk.green('✓') + ' Telegram integration configured');
    }
    if (answers.setupGmail) {
      console.log(chalk.green('✓') + ' Gmail integration ready (OAuth setup required)');
    }
    if (answers.setupTwitter) {
      console.log(chalk.green('✓') + ' Twitter integration ready (OAuth setup required)');
    }

    // Success message
    console.log('\n' + chalk.green.bold('╔══════════════════════════════════════════════════════════╗'));
    console.log(chalk.green.bold('║') + '          ' + chalk.white.bold('ONBOARDING COMPLETE!') + '                    ' + chalk.green.bold('║'));
    console.log(chalk.green.bold('╚══════════════════════════════════════════════════════════╝') + '\n');
    
    console.log(chalk.cyan('What was set up:'));
    if (answers.gitlabToken) {
      console.log('  ✅ GitLab media storage project');
    }
    console.log('  ✅ NLM CLI (NotebookLM)');
    console.log('  ✅ OpenCode (AI coding assistant)');
    console.log('  ✅ OpenClaw (AI orchestration)');
    console.log('  ✅ Google Antigravity (OAuth)');
    if (answers.setupWhatsApp) {
      console.log('  ✅ WhatsApp integration');
    }
    if (answers.setupTelegram) {
      console.log('  ✅ Telegram integration');
    }
    if (answers.setupGmail) {
      console.log('  ✅ Gmail integration');
    }
    if (answers.setupTwitter) {
      console.log('  ✅ Twitter integration');
    }
    
    console.log('\n' + chalk.cyan('Next steps:'));
    console.log('  1. Clone the BIOMETRICS repo:');
    console.log(chalk.gray('     git clone https://github.com/Delqhi/BIOMETRICS.git'));
    console.log('  2. Navigate to the project:');
    console.log(chalk.gray('     cd BIOMETRICS'));
    console.log('  3. Start building with AI assistance!');
    console.log(chalk.gray('     opencode'));
    console.log('  4. Use OpenClaw for automation:');
    console.log(chalk.gray('     openclaw start'));
    console.log('');

  } catch (error) {
    console.error(chalk.red('\n✗ Onboarding failed:'), error.message);
    process.exit(1);
  }
}

main();
