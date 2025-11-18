import readline from 'readline';
import chalk from 'chalk';
import { MastraClient } from '@mastra/client-js';
import ora, { Ora } from 'ora';
import { animateMarkdown } from './markdown-renderer.js';

interface Message {
  role: 'user' | 'assistant';
  content: string;
}

export class InteractiveCLI {
  private client: MastraClient;
  private conversationHistory: Message[] = [];
  private rl: readline.Interface;
  private spinner: Ora | null = null;

  constructor(apiUrl: string) {
    this.client = new MastraClient({ baseUrl: apiUrl });
    this.rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
      prompt: chalk.cyan.bold('You: '),
    });
  }

  async start() {
    console.clear();
    this.printWelcome();
    
    this.rl.prompt();

    this.rl.on('line', async (input: string) => {
      const trimmedInput = input.trim();

      if (!trimmedInput) {
        this.rl.prompt();
        return;
      }

      // Handle special commands
      if (trimmedInput === '/exit' || trimmedInput === '/quit') {
        this.exit();
        return;
      }

      if (trimmedInput === '/clear') {
        console.clear();
        this.conversationHistory = [];
        console.log(chalk.green('✓ Conversation cleared\n'));
        this.rl.prompt();
        return;
      }

      if (trimmedInput === '/history') {
        this.showHistory();
        this.rl.prompt();
        return;
      }

      if (trimmedInput === '/help') {
        this.showHelp();
        this.rl.prompt();
        return;
      }

      // Pause readline during processing
      this.rl.pause();

      // Process user query
      await this.processQuery(trimmedInput);
      
      // Resume and show prompt
      this.rl.resume();
      this.rl.prompt();
    });

    this.rl.on('close', () => {
      this.exit();
    });

    // Keep process alive
    this.rl.on('SIGINT', () => {
      console.log('\n');
      console.log(chalk.yellow('Press Ctrl+C again or type /exit to quit'));
      this.rl.prompt();
    });
  }

  private async processQuery(query: string) {
    // Add user message to history
    this.conversationHistory.push({
      role: 'user',
      content: query,
    });

    // Start spinner
    this.spinner = ora({
      text: chalk.gray('Thinking...'),
      color: 'cyan',
      spinner: 'dots',
    }).start();

    try {
      const agent = this.client.getAgent('financialAgent');
      
      // Use generate
      const response = await agent.generate({
        messages: this.conversationHistory,
      });

      // Stop spinner
      this.spinner.stop();

      // Check if we got a response
      if (!response || !response.text) {
        throw new Error('No response received from agent');
      }

      const text = response.text;

      // Display response with markdown rendering and animation
      console.log('');
      console.log(chalk.green.bold('🤖 Agent:'));
      console.log('');
      
      await this.animateMarkdownResponse(text);

      console.log('');

      // Add assistant response to history
      this.conversationHistory.push({
        role: 'assistant',
        content: text,
      });
    } catch (error: any) {
      if (this.spinner) {
        this.spinner.stop();
      }
      console.log('');
      
      // Check for connection errors
      if (error?.code === 'ECONNREFUSED' || error?.cause?.code === 'ECONNREFUSED') {
        console.error(chalk.red.bold('✗ Connection Error:'), chalk.red('Cannot connect to Mastra server'));
        console.log('');
        console.log(chalk.yellow('Make sure the Mastra server is running:'));
        console.log(chalk.dim('  1. Navigate to your kainos-agent-core directory'));
        console.log(chalk.dim('  2. Run: npm run dev'));
        console.log('');
        console.log(chalk.yellow('Or specify a different URL:'));
        console.log(chalk.dim('  financial-agent --url http://your-server:port'));
        console.log(chalk.dim('  MASTRA_API_URL=http://your-server:port financial-agent'));
      } else {
        console.error(chalk.red.bold('✗ Error:'), chalk.red(error?.message || 'Unknown error'));
        if (error?.cause) {
          console.error(chalk.gray('Cause:'), error.cause);
        }
      }
      
      console.log('');
      // Remove the failed user message from history
      this.conversationHistory.pop();
    }
  }

  private async animateMarkdownResponse(text: string) {
    await animateMarkdown(text, 15);
  }

  private printWelcome() {
    const banner = `
╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║         ${chalk.bold.cyan('Financial Agent')} - ${chalk.bold.white('Interactive Mode')}                 ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝
`;
    console.log(banner);
    console.log(chalk.white('Ask me anything about stocks, markets, and financial data!'));
    console.log('');
    console.log(chalk.gray('💡 Examples:'));
    console.log(chalk.gray('  • What is the current price of Apple stock?'));
    console.log(chalk.gray('  • Show me Tesla\'s financial metrics'));
    console.log(chalk.gray('  • Get the latest 10-K filing for Microsoft'));
    console.log(chalk.gray('  • What are the recent insider trades for NVIDIA?'));
    console.log(chalk.gray('  • Compare Apple and Microsoft'));
    console.log('');
    console.log(chalk.yellow('⚡ Commands:'));
    console.log(chalk.yellow('  /help    ') + chalk.gray('- Show this help'));
    console.log(chalk.yellow('  /clear   ') + chalk.gray('- Clear conversation history'));
    console.log(chalk.yellow('  /history ') + chalk.gray('- Show conversation history'));
    console.log(chalk.yellow('  /exit    ') + chalk.gray('- Exit interactive mode'));
    console.log('');
  }

  private showHelp() {
    console.log('');
    console.log(chalk.bold.cyan('═══════════════════════════════════════════════════════════'));
    console.log(chalk.bold.cyan('                    Available Commands                      '));
    console.log(chalk.bold.cyan('═══════════════════════════════════════════════════════════'));
    console.log('');
    console.log(chalk.yellow('  /help    ') + chalk.gray('- Show this help message'));
    console.log(chalk.yellow('  /clear   ') + chalk.gray('- Clear conversation history'));
    console.log(chalk.yellow('  /history ') + chalk.gray('- Show conversation history'));
    console.log(chalk.yellow('  /exit    ') + chalk.gray('- Exit interactive mode'));
    console.log('');
    console.log(chalk.bold.cyan('═══════════════════════════════════════════════════════════'));
    console.log(chalk.bold.cyan('                  What I Can Help With                      '));
    console.log(chalk.bold.cyan('═══════════════════════════════════════════════════════════'));
    console.log('');
    console.log(chalk.white('  📊 Stock prices and market data'));
    console.log(chalk.white('  💰 Financial metrics (P/E, ROE, margins, etc.)'));
    console.log(chalk.white('  🏢 Company information'));
    console.log(chalk.white('  📄 SEC filings (10-K, 10-Q, 8-K)'));
    console.log(chalk.white('  💼 Insider trading activity'));
    console.log(chalk.white('  📰 Market news and sentiment'));
    console.log(chalk.white('  📈 Portfolio analysis'));
    console.log(chalk.white('  🌍 Economic indicators'));
    console.log('');
  }

  private showHistory() {
    console.log('');
    console.log(chalk.bold.cyan('═══════════════════════════════════════════════════════════'));
    console.log(chalk.bold.cyan('                  Conversation History                      '));
    console.log(chalk.bold.cyan('═══════════════════════════════════════════════════════════'));
    console.log('');

    if (this.conversationHistory.length === 0) {
      console.log(chalk.gray('  No messages yet'));
      console.log('');
      return;
    }

    this.conversationHistory.forEach((msg, index) => {
      const prefix = msg.role === 'user' 
        ? chalk.cyan.bold('👤 You: ')
        : chalk.green.bold('🤖 Agent: ');
      
      const content = msg.content.length > 100 
        ? msg.content.substring(0, 100) + '...'
        : msg.content;

      console.log(`${chalk.gray(`${index + 1}.`)} ${prefix}`);
      console.log(chalk.white(`   ${content}`));
      console.log('');
    });
  }

  private exit() {
    console.log('');
    console.log(chalk.cyan('═══════════════════════════════════════════════════════════'));
    console.log(chalk.cyan.bold('   Thanks for using Financial Agent! Goodbye! 👋'));
    console.log(chalk.cyan('═══════════════════════════════════════════════════════════'));
    console.log('');
    process.exit(0);
  }
}
