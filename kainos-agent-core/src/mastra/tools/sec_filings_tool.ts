import { createTool } from '@mastra/core/tools';
import { z } from 'zod';

export const secFilingsTool = createTool({
  id: 'sec-filings-tool',
  description: 'Get SEC filings for a company including 10-K, 10-Q, 8-K, and other regulatory documents',
  inputSchema: z.object({
    ticker: z.string().optional().describe('The ticker symbol of the company'),
    cik: z.string().optional().describe('The Central Index Key (CIK) of the company'),
    filing_type: z.enum(['10-K', '10-Q', '8-K', '4', '144']).optional().describe('The type of SEC filing'),
  }),
  outputSchema: z.object({
    filings: z.array(z.object({
      cik: z.number(),
      accession_number: z.string(),
      filing_type: z.string(),
      report_date: z.string(),
      ticker: z.string(),
      url: z.string(),
    })),
  }),
  execute: async ({ context }) => {
    const { ticker, cik, filing_type } = context;
    
    const params = new URLSearchParams();
    if (ticker) params.append('ticker', ticker);
    if (cik) params.append('cik', cik);
    if (filing_type) params.append('filing_type', filing_type);

    const response = await fetch(
      `https://api.financialdatasets.ai/filings?${params.toString()}`,
      {
        headers: {
          'X-API-KEY': process.env.FINANCIAL_DATASETS_API_KEY || '',
        },
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch SEC filings: ${response.statusText}`);
    }

    return await response.json();
  },
});
