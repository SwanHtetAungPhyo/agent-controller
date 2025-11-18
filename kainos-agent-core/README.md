# stock-agent-core

Financial analysis agent powered by Mastra with comprehensive stock market tools and workflows.

## Installation

```bash
npm install stock-agent-core --legacy-peer-deps
```

## Usage

```typescript
import { mastra, financialAgent } from 'stock-agent-core';

// Use the financial agent
const result = await financialAgent.generate('What is the current price of AAPL?');
console.log(result.text);

// Access individual tools
import { 
  getStockPrice, 
  getCompanyInfo,
  getMarketNews 
} from 'stock-agent-core';

const price = await getStockPrice({ symbol: 'AAPL' });
const info = await getCompanyInfo({ symbol: 'TSLA' });
const news = await getMarketNews({ query: 'tech stocks' });
```

## Features

- **Financial Agent**: AI-powered agent for stock market analysis
- **Stock Tools**: Price tracking, company info, market indices
- **Market Data**: News, SEC filings, insider trades
- **Analysis Tools**: Portfolio analysis, financial metrics, economic indicators
- **Workflows**: Pre-built financial analysis workflows

## Available Tools

- `getStockPrice` - Real-time stock prices
- `getCompanyInfo` - Company information and fundamentals
- `getMarketNews` - Latest market news
- `getSecFilings` - SEC filings and reports
- `getInsiderTrades` - Insider trading activity
- `getMarketIndices` - Market index data
- `getEconomicIndicators` - Economic data
- `getFinancialMetrics` - Financial ratios and metrics
- `analyzePortfolio` - Portfolio analysis

## License

ISC
