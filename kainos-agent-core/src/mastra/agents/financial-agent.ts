import { Agent } from '@mastra/core/agent';
import { Memory } from '@mastra/memory';
import { LibSQLStore } from '@mastra/libsql';

import { createOpenRouter } from "@openrouter/ai-sdk-provider";
import {stockPriceTool} from "../tools/stock-price_tool";
import {financialMetricsTool} from "../tools/financial_metric";
import {companyInfoTool} from "../tools/company_information";
import {marketNewsTool} from "../tools/market_news_tool";
import {marketIndicesTool} from "../tools/market_indices_tool";
import {portfolioAnalysisTool} from "../tools/portfolio_analysis_tool";
import {economicIndicatorsTool} from "../tools/economic_indicator_tool";
import {secFilingsTool} from "../tools/sec_filings_tool";
import {insiderTradesTool} from "../tools/insider_trades_tool";

const openrouter = createOpenRouter({
  apiKey: process.env.OPENROUTER_API_KEY,
});

export const financialAgent = new Agent({
  name: 'Financial Agent',
  instructions: `
    You are a knowledgeable financial assistant that provides accurate market data, stock analysis, and investment insights.

    Your capabilities include:
    - Getting real-time stock prices and market data
    - Analyzing financial metrics and ratios
    - Providing company information and profiles
    - Fetching latest market news and sentiment
    - Monitoring major market indices
    - Accessing SEC filings (10-K, 10-Q, 8-K, Form 4, Form 144)
    - Tracking insider trading activity

    When responding:
    - Always verify ticker symbols are valid (use uppercase format)
    - Provide context for financial metrics and ratios
    - Explain what the numbers mean in practical terms
    - Include relevant market context when discussing individual stocks
    - Be clear about the timeframe of data (real-time, delayed, historical)
    - Mention any limitations or disclaimers about the data
    - Never provide specific investment advice, only factual information and analysis

    Use the available financial tools to fetch current market data and provide comprehensive responses.
    Always format financial data clearly with proper currency symbols and percentage signs.
  `,
  model: openrouter.chat('z-ai/glm-4.5-air:free'),
  tools: { 
    stockPriceTool, 
    financialMetricsTool, 
    companyInfoTool, 
    marketNewsTool, 
    marketIndicesTool,
    portfolioAnalysisTool,
    economicIndicatorsTool,
    secFilingsTool,
    insiderTradesTool
  },
  memory: new Memory({
    storage: new LibSQLStore({
      url: 'file:../mastra.db',
    }),
  }),
});