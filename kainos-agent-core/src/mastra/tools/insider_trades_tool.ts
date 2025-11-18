import { createTool } from '@mastra/core/tools';
import { z } from 'zod';

export const insiderTradesTool = createTool({
  id: 'insider-trades-tool',
  description: 'Get insider trading activity for a company including purchases, sales, and Form 4 filings',
  inputSchema: z.object({
    ticker: z.string().describe('The ticker symbol of the company'),
    limit: z.number().optional().describe('Maximum number of trades to return (default: 10)'),
  }),
  outputSchema: z.object({
    trades: z.array(z.object({
      ticker: z.string(),
      insider_name: z.string(),
      transaction_date: z.string(),
      transaction_type: z.string(),
      shares: z.number(),
      price_per_share: z.number(),
      total_value: z.number(),
    })),
  }),
  execute: async ({ context }) => {
    const { ticker, limit = 10 } = context;
    
    const params = new URLSearchParams({
      ticker,
      limit: limit.toString(),
    });

    const response = await fetch(
      `https://api.financialdatasets.ai/insider-trades?${params.toString()}`,
      {
        headers: {
          'X-API-KEY': process.env.FINANCIAL_DATASETS_API_KEY || '',
        },
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch insider trades: ${response.statusText}`);
    }

    return await response.json();
  },
});
