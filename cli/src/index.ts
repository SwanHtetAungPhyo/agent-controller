#!/usr/bin/env node

import { Command } from 'commander';
import { config } from 'dotenv';
import chalk from 'chalk';
import ora from 'ora';
import { MastraClient } from '@mastra/client-js';
import { InteractiveCLI } from './interactive-enhanced.js';
import { readFileSync, existsSync } from 'fs';
import { homedir } from 'os';
import { join } from 'path';

config();

// Load config file if exists
function loadConfigFile(): { apiUrl?: string } {
  const configPath = join(homedir(), '.financial-agent-config.json');
  if (existsSync(configPath)) {
    try {
      const configData = readFileSync(configPath, 'utf-8');
      return JSON.parse(configData);
    } catch (error) {
      // Ignore config file errors
    }
  }
  return {};
}

const configFile = loadConfigFile();

const program = new Command();

// Helper to get API URL from options (priority: CLI flag > env var > config file > default)
function getApiUrl(options?: { url?: string; port?: string }): string {
  let url = options?.url || process.env.MASTRA_API_URL || configFile.apiUrl || 'http://localhost:4111';
  
  if (options?.port || process.env.MASTRA_PORT) {
    const port = options?.port || process.env.MASTRA_PORT;
    const urlObj = new URL(url);
    urlObj.port = port!;
    url = urlObj.toString();
  }
  
  return url;
}

const mastraClient = new MastraClient({
  baseUrl: getApiUrl(),
});

program
  .name('financial-agent')
  .description('CLI tool for Financial Agent - Get stock data, SEC filings, and insider trades')
  .version('1.0.2')
  .option('-u, --url <url>', 'Mastra API URL', process.env.MASTRA_API_URL || 'http://localhost:4111')
  .option('-p, --port <port>', 'Mastra API port (overrides URL port)', process.env.MASTRA_PORT)
  .allowUnknownOption(false);

// Interactive mode command
program
  .command('chat')
  .description('Start interactive chat mode (like Gemini)')
  .action(async () => {
    const opts = program.opts();
    const apiUrl = getApiUrl(opts);
    console.log(chalk.dim(`Connecting to: ${apiUrl}\n`));
    const interactive = new InteractiveCLI(apiUrl);
    await interactive.start();
  });

// Stock price command
program
  .command('price <ticker>')
  .description('Get current stock price for a ticker')
  .action(async (ticker: string) => {
    const spinner = ora(`Fetching price for ${ticker.toUpperCase()}...`).start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: `Get the current stock price for ${ticker}`,
          },
        ],
      });

      spinner.succeed('Price fetched successfully');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to fetch price');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// Financial metrics command
program
  .command('metrics <ticker>')
  .description('Get financial metrics for a ticker')
  .action(async (ticker: string) => {
    const spinner = ora(`Fetching metrics for ${ticker.toUpperCase()}...`).start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: `Show me the financial metrics for ${ticker}`,
          },
        ],
      });

      spinner.succeed('Metrics fetched successfully');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to fetch metrics');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// SEC filings command
program
  .command('filings <ticker>')
  .description('Get SEC filings for a ticker')
  .option('-t, --type <type>', 'Filing type (10-K, 10-Q, 8-K)', '10-K')
  .action(async (ticker: string, options: { type: string }) => {
    const spinner = ora(`Fetching ${options.type} filings for ${ticker.toUpperCase()}...`).start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: `Get the latest ${options.type} filings for ${ticker}`,
          },
        ],
      });

      spinner.succeed('Filings fetched successfully');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to fetch filings');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// Insider trades command
program
  .command('insider <ticker>')
  .description('Get insider trading activity for a ticker')
  .action(async (ticker: string) => {
    const spinner = ora(`Fetching insider trades for ${ticker.toUpperCase()}...`).start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: `Show me recent insider trading activity for ${ticker}`,
          },
        ],
      });

      spinner.succeed('Insider trades fetched successfully');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to fetch insider trades');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// News command
program
  .command('news <ticker>')
  .description('Get latest news for a ticker')
  .action(async (ticker: string) => {
    const spinner = ora(`Fetching news for ${ticker.toUpperCase()}...`).start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: `What is the latest news about ${ticker}?`,
          },
        ],
      });

      spinner.succeed('News fetched successfully');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to fetch news');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// Analyze command (comprehensive)
program
  .command('analyze <ticker>')
  .description('Get comprehensive analysis for a ticker')
  .action(async (ticker: string) => {
    const spinner = ora(`Analyzing ${ticker.toUpperCase()}...`).start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: `Give me a comprehensive analysis of ${ticker} including price, metrics, recent news, and insider activity`,
          },
        ],
      });

      spinner.succeed('Analysis complete');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to analyze');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// Custom query command
program
  .command('query <question>')
  .description('Ask a custom question to the financial agent')
  .action(async (question: string) => {
    const spinner = ora('Processing query...').start();
    
    try {
      const agent = mastraClient.getAgent('financialAgent');
      const response = await agent.generate({
        messages: [
          {
            role: 'user',
            content: question,
          },
        ],
      });

      spinner.succeed('Query processed');
      console.log('\n' + chalk.cyan(response.text) + '\n');
    } catch (error) {
      spinner.fail('Failed to process query');
      console.error(chalk.red('Error:'), error);
      process.exit(1);
    }
  });

// Config command
program
  .command('config')
  .description('Show current configuration')
  .option('--set-url <url>', 'Set the Mastra API URL in config file')
  .action(async (options: { setUrl?: string }) => {
    const configPath = join(homedir(), '.financial-agent-config.json');
    
    if (options.setUrl) {
      // Save config
      const { writeFileSync } = await import('fs');
      const config = { apiUrl: options.setUrl };
      writeFileSync(configPath, JSON.stringify(config, null, 2));
      console.log(chalk.green('✓ Configuration saved to:'), chalk.dim(configPath));
      console.log(chalk.cyan('API URL:'), options.setUrl);
    } else {
      // Show current config
      const opts = program.opts();
      const currentUrl = getApiUrl(opts);
      
      console.log(chalk.bold('\n📋 Current Configuration\n'));
      console.log(chalk.cyan('API URL:'), currentUrl);
      console.log(chalk.dim('\nConfiguration priority:'));
      console.log(chalk.dim('  1. --url flag'));
      console.log(chalk.dim('  2. MASTRA_API_URL environment variable'));
      console.log(chalk.dim('  3. ~/.financial-agent-config.json'));
      console.log(chalk.dim('  4. Default (http://localhost:4111)'));
      console.log(chalk.dim('\nConfig file:'), configPath);
      console.log(chalk.dim('Exists:'), existsSync(configPath) ? chalk.green('Yes') : chalk.yellow('No'));
      
      if (existsSync(configPath)) {
        console.log(chalk.dim('Content:'), JSON.stringify(configFile, null, 2));
      }
      console.log('');
    }
  });

// Check if a command was provided
const args = process.argv.slice(2);
const commands = ['chat', 'price', 'metrics', 'filings', 'insider', 'news', 'analyze', 'query', 'config', 'help'];
const hasCommand = args.some(arg => commands.includes(arg));

// If no command provided (only flags or nothing), start interactive mode
if (!hasCommand && args.length === 0) {
  // No arguments at all - start interactive
  const apiUrl = getApiUrl();
  console.log(chalk.dim(`Connecting to: ${apiUrl}\n`));
  const interactive = new InteractiveCLI(apiUrl);
  interactive.start();
} else if (!hasCommand && args.every(arg => arg.startsWith('-'))) {
  // Only flags provided - parse them and start interactive
  program.parse();
  const opts = program.opts();
  const apiUrl = getApiUrl(opts);
  console.log(chalk.dim(`Connecting to: ${apiUrl}\n`));
  const interactive = new InteractiveCLI(apiUrl);
  interactive.start();
} else {
  // Command provided - parse and execute
  program.parse();
}
