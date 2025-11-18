# Financial Agent CLI

Interactive command-line interface for the Financial Agent. Get real-time stock data, SEC filings, insider trades, and comprehensive market analysis powered by AI.

## Features

- 🎨 **Beautiful Interactive Mode** - Gemini-like conversational interface
- 📊 **Real-time Stock Data** - Current prices, volume, 52-week ranges
- 💰 **Financial Metrics** - P/E ratios, ROE, margins, and more
- 📄 **SEC Filings** - Access 10-K, 10-Q, 8-K, Form 4, Form 144
- 💼 **Insider Trades** - Track insider trading activity
- 📰 **Market News** - Latest news with sentiment analysis
- 📈 **Portfolio Analysis** - Multi-stock portfolio insights
- 🌍 **Economic Indicators** - GDP, inflation, unemployment data
- ✨ **Markdown Rendering** - Formatted, colorful output
- ⚡ **Smooth Animations** - Typing effects and spinners

## Installation

### Global Installation

```bash
npm install -g kainos-financial-agent
```

### Local Installation

```bash
npm install kainos-financial-agent
```

### Upgrade to Latest Version

```bash
# Update to latest version
npm update -g kainos-financial-agent

# Or force reinstall
npm install -g kainos-financial-agent@latest

# Check your current version
financial-agent --version
```

## Prerequisites

You need a running Mastra server with the Financial Agent deployed.

1. Set up the `kainos-agent-core` project
2. Start the Mastra server:
   ```bash
   cd kainos-agent-core
   npm run dev
   ```
3. The server runs on `http://localhost:4111` by default

## Configuration

Configure the Mastra API URL using any of these methods:

### Method 1: Environment Variable (Recommended)

```bash
export MASTRA_API_URL=http://localhost:4111
financial-agent
```

Or create a `.env` file in your project:
```env
MASTRA_API_URL=http://localhost:4111
```

### Method 2: Command Line Flags

```bash
# Specify full URL
financial-agent --url http://localhost:4111

# Or just the port
financial-agent --port 4111

# Works with any command
financial-agent price AAPL --url http://localhost:4111
```

### Method 3: Config File

Save your configuration permanently:

```bash
# Set the API URL in config file
financial-agent config --set-url http://localhost:4111

# View current configuration
financial-agent config
```

This creates `~/.financial-agent-config.json` with your settings.

### Method 4: Per-Command

```bash
MASTRA_API_URL=http://your-server:8080 financial-agent chat
```

## Quick Start

### Interactive Mode (Recommended)

Simply run the command to start an interactive chat session:

```bash
financial-agent
```

Or explicitly:

```bash
financial-agent chat
```

### Example Conversation

```
You: What is the current price of Apple stock?

Agent: The current price of Apple Inc. (AAPL) stock is $267.46.

Here are some additional details:
  • Change: -$4.95
  • Change Percentage: -1.82%
  • Volume: 42,973,636 shares traded
  • 52-Week High: $270.49
  • 52-Week Low: $265.73

You: Now show me Tesla's financial metrics

Agent: Here are the key financial metrics for Tesla, Inc. (TSLA):
  • P/E Ratio: 276.95
  • P/B Ratio: 16.82
  • Return on Equity (ROE): 6.79%
  ...
```

### Command Mode

For quick one-off queries:

```bash
# Stock price
financial-agent price AAPL

# Financial metrics
financial-agent metrics TSLA

# SEC filings
financial-agent filings MSFT --type 10-K

# Insider trades
financial-agent insider NVDA

# Latest news
financial-agent news AMZN

# Comprehensive analysis
financial-agent analyze GOOGL

# Custom query
financial-agent query "Compare Apple and Microsoft"
```

## Commands

### Interactive Mode Commands

| Command | Description |
|---------|-------------|
| `/help` | Show help message |
| `/clear` | Clear conversation history |
| `/history` | Show conversation history |
| `/exit` | Exit interactive mode |

### CLI Commands

| Command | Description | Example |
|---------|-------------|---------|
| `chat` | Start interactive mode | `financial-agent chat` |
| `price <ticker>` | Get stock price | `financial-agent price AAPL` |
| `metrics <ticker>` | Get financial metrics | `financial-agent metrics TSLA` |
| `filings <ticker>` | Get SEC filings | `financial-agent filings MSFT` |
| `insider <ticker>` | Get insider trades | `financial-agent insider NVDA` |
| `news <ticker>` | Get latest news | `financial-agent news AMZN` |
| `analyze <ticker>` | Comprehensive analysis | `financial-agent analyze GOOGL` |
| `query <question>` | Custom query | `financial-agent query "..."` |

### Filing Types

Use `--type` or `-t` flag with the `filings` command:

```bash
financial-agent filings AAPL --type 10-K  # Annual report
financial-agent filings AAPL --type 10-Q  # Quarterly report
financial-agent filings AAPL --type 8-K   # Current report
financial-agent filings AAPL --type 4     # Insider trading
financial-agent filings AAPL --type 144   # Restricted stock
```

## Setup

### 1. Environment Configuration

Create a `.env` file or set environment variable:

```bash
MASTRA_API_URL=http://localhost:4111
```

For production, point to your deployed Mastra server:

```bash
MASTRA_API_URL=https://your-mastra-server.com
```

### 2. Start Mastra Server

The CLI requires a running Mastra server with the Financial Agent. 

If you're developing locally:

```bash
cd your-mastra-project
pnpm dev
```

The server should be accessible at `http://localhost:4111` (or your configured URL).

### 3. Verify Connection

Test the connection:

```bash
curl http://localhost:4111/api/agents
```

You should see the `financialAgent` in the response.

## Usage Examples

### Stock Analysis

```bash
# Quick price check
financial-agent price AAPL

# Deep analysis
financial-agent analyze TSLA

# Compare stocks
financial-agent query "Compare AAPL and MSFT performance"
```

### Regulatory Filings

```bash
# Get annual reports
financial-agent filings MSFT --type 10-K

# Get quarterly reports
financial-agent filings GOOGL --type 10-Q

# Get current reports
financial-agent filings AMZN --type 8-K
```

### Insider Activity

```bash
# Check insider trades
financial-agent insider NVDA

# Ask about specific insider
financial-agent query "What insider trades has Elon Musk made recently?"
```

### Market Research

```bash
# Get latest news
financial-agent news AAPL

# Sector analysis
financial-agent query "What are the top performing tech stocks?"

# Market overview
financial-agent query "Give me a market overview for today"
```

## Interactive Mode Features

### Markdown Rendering

Responses are beautifully formatted with:
- **Bold text** for emphasis
- *Italic text* for notes
- `Code` highlighting
- Colored currency ($100) and percentages (50%)
- Bullet points with •
- Clickable links
- Styled headers

### Animations

- Smooth line-by-line text rendering
- Spinner while processing queries
- Typing effects for dramatic responses

### Context Awareness

The interactive mode maintains conversation history, so you can ask follow-up questions:

```
You: What is Apple's P/E ratio?
Agent: Apple's P/E ratio is 35.2...

You: How does that compare to the industry average?
Agent: Compared to the tech industry average of 28.5...
```

## Configuration

### Custom API URL

Set via environment variable:

```bash
export MASTRA_API_URL=https://your-server.com
financial-agent
```

Or create a `.env` file in your working directory:

```
MASTRA_API_URL=https://your-server.com
```

### Animation Speed

The CLI uses sensible defaults, but you can modify the source code if needed:
- Line delay: 15ms (in `interactive-enhanced.ts`)
- Character delay: 3ms (for typing effect)

## Troubleshooting

### Connection Errors

**Error**: `fetch failed` or `ECONNREFUSED`

**Solution**:
1. Verify Mastra server is running: `curl http://localhost:4111/api/agents`
2. Check `MASTRA_API_URL` environment variable
3. Ensure firewall allows connections

### No Response

**Error**: `No response received from agent`

**Solution**:
1. Check server logs for errors
2. Verify API keys are configured in the Mastra server
3. Test with a simple query first

### Command Not Found

**Error**: `financial-agent: command not found`

**Solution**:
1. Reinstall globally: `npm install -g @kainos/financial-agent-cli`
2. Check npm global bin path: `npm bin -g`
3. Add to PATH if needed

## Development

### Build from Source

```bash
git clone <repository>
cd cli
pnpm install
pnpm build
```

### Run Locally

```bash
pnpm dev
```

### Test Commands

```bash
node dist/index.js price AAPL
node dist/index.js chat
```

## Requirements

- Node.js v20 or higher
- Running Mastra server with Financial Agent
- Internet connection for market data

## Data Sources

The Financial Agent uses:
- Financial Datasets API for market data
- SEC EDGAR for regulatory filings
- Real-time news feeds
- Economic data providers

## License

MIT

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report a bug](https://github.com/yourusername/financial-agent-cli/issues)
- Documentation: [Full docs](https://github.com/yourusername/financial-agent-cli)

## Disclaimer

This tool provides financial data for informational purposes only. It is not financial advice. Always do your own research and consult with financial professionals before making investment decisions.

## Credits

Built with:
- [Mastra](https://mastra.ai) - AI agent framework
- [Commander.js](https://github.com/tj/commander.js) - CLI framework
- [Chalk](https://github.com/chalk/chalk) - Terminal styling
- [Ora](https://github.com/sindresorhus/ora) - Spinners

---

Made with ❤️ by Kainos
