import { createStep, createWorkflow } from '@mastra/core/workflows';
import { z } from 'zod';

const stockAnalysisSchema = z.object({
  symbol: z.string(),
  price: z.number(),
  change: z.number(),
  changePercent: z.number(),
  volume: z.number(),
  marketCap: z.number().optional(),
  peRatio: z.number().optional(),
  pbRatio: z.number().optional(),
  roe: z.number().optional(),
  companyName: z.string(),
  sector: z.string().optional(),
  industry: z.string().optional(),
  newsCount: z.number(),
  sentiment: z.string().optional(),
});

const fetchStockData = createStep({
  id: 'fetch-stock-data',
  description: 'Fetches comprehensive stock data including price, metrics, and company info',
  inputSchema: z.object({
    symbol: z.string().describe('Stock ticker symbol'),
  }),
  outputSchema: stockAnalysisSchema,
  execute: async ({ inputData, mastra }) => {
    if (!inputData) {
      throw new Error('Input data not found');
    }

    const agent = mastra?.getAgent('financialAgent');
    if (!agent) {
      throw new Error('Financial agent not found');
    }

    // Fetch stock price
    const priceResponse = await agent.stream([
      {
        role: 'user',
        content: `Get the current stock price for ${inputData.symbol}`,
      },
    ]);

    // Fetch financial metrics
    const metricsResponse = await agent.stream([
      {
        role: 'user',
        content: `Get financial metrics for ${inputData.symbol}`,
      },
    ]);

    // Fetch company info
    const companyResponse = await agent.stream([
      {
        role: 'user',
        content: `Get company information for ${inputData.symbol}`,
      },
    ]);

    // Fetch recent news
    const newsResponse = await agent.stream([
      {
        role: 'user',
        content: `Get recent news for ${inputData.symbol}`,
      },
    ]);

    // For this example, we'll return mock data structure
    // In a real implementation, you'd parse the agent responses
    return {
      symbol: inputData.symbol.toUpperCase(),
      price: 150.00, // This would come from actual API response
      change: 2.50,
      changePercent: 1.69,
      volume: 1000000,
      marketCap: 2500000000,
      peRatio: 25.5,
      pbRatio: 3.2,
      roe: 15.8,
      companyName: "Sample Company",
      sector: "Technology",
      industry: "Software",
      newsCount: 5,
      sentiment: "positive",
    };
  },
});

const generateAnalysis = createStep({
  id: 'generate-analysis',
  description: 'Generates comprehensive stock analysis report',
  inputSchema: stockAnalysisSchema,
  outputSchema: z.object({
    analysis: z.string(),
  }),
  execute: async ({ inputData, mastra }) => {
    const stockData = inputData;

    if (!stockData) {
      throw new Error('Stock data not found');
    }

    const agent = mastra?.getAgent('financialAgent');
    if (!agent) {
      throw new Error('Financial agent not found');
    }

    const prompt = `Based on the following stock data for ${stockData.symbol}, provide a comprehensive analysis:

    📊 STOCK DATA SUMMARY
    ${JSON.stringify(stockData, null, 2)}

    Please structure your analysis as follows:

    🏢 COMPANY OVERVIEW
    • Company: [Name]
    • Sector: [Sector] | Industry: [Industry]
    • Current Price: $[Price] ([Change]% change)

    📈 PRICE PERFORMANCE
    • Current Price: $[Price]
    • Daily Change: $[Change] ([ChangePercent]%)
    • Volume: [Volume] shares
    • Market Cap: $[MarketCap]

    🔍 FINANCIAL METRICS ANALYSIS
    • P/E Ratio: [PE] - [Interpretation: overvalued/undervalued/fair]
    • P/B Ratio: [PB] - [What this indicates about book value]
    • ROE: [ROE]% - [Performance assessment]
    • [Additional context about what these metrics mean]

    📰 MARKET SENTIMENT
    • Recent News Articles: [Count]
    • Overall Sentiment: [Sentiment]
    • [Brief interpretation of news impact]

    ⚖️ INVESTMENT CONSIDERATIONS
    Strengths:
    • [List 2-3 positive factors based on the data]

    Risks:
    • [List 2-3 potential concerns or risks]

    📋 KEY TAKEAWAYS
    • [3-4 bullet points summarizing the most important insights]

    ⚠️ DISCLAIMER
    This analysis is based on current market data and should not be considered as investment advice. Always conduct your own research and consult with financial professionals before making investment decisions.

    Keep the analysis factual, balanced, and educational. Focus on interpreting the data rather than making buy/sell recommendations.`;

    const response = await agent.stream([
      {
        role: 'user',
        content: prompt,
      },
    ]);

    let analysisText = '';

    for await (const chunk of response.textStream) {
      process.stdout.write(chunk);
      analysisText += chunk;
    }

    return {
      analysis: analysisText,
    };
  },
});

const financialWorkflow = createWorkflow({
  id: 'financial-workflow',
  inputSchema: z.object({
    symbol: z.string().describe('Stock ticker symbol to analyze'),
  }),
  outputSchema: z.object({
    analysis: z.string(),
  }),
})
  .then(fetchStockData)
  .then(generateAnalysis);

financialWorkflow.commit();

export { financialWorkflow };