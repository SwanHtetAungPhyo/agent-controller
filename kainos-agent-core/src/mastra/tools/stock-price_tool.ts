import {createTool} from "@mastra/core/tools";
import {z} from "zod";
import {makeFinancialApiRequest} from "./financial-tools";


export const stockPriceTool = createTool({
    id: 'get-stock-price',
    description: 'Get current stock price and basic information for a given ticker symbol',
    inputSchema: z.object({
        symbol: z.string().describe('Stock ticker symbol (e.g., AAPL, TSLA, MSFT)'),
    }),
    outputSchema: z.object({
        symbol: z.string(),
        price: z.number(),
        change: z.number(),
        changePercent: z.number(),
        volume: z.number(),
        marketCap: z.number().optional(),
        high52Week: z.number().optional(),
        low52Week: z.number().optional(),
    }),
    execute: async ({ context }) => {
        const data = await makeFinancialApiRequest('GLOBAL_QUOTE', context.symbol.toUpperCase());

        const quote = data['Global Quote'];
        if (!quote) {
            throw new Error(`No data found for symbol ${context.symbol}`);
        }

        return {
            symbol: quote['01. symbol'],
            price: parseFloat(quote['05. price']),
            change: parseFloat(quote['09. change']),
            changePercent: parseFloat(quote['10. change percent'].replace('%', '')),
            volume: parseInt(quote['06. volume']),
            marketCap: undefined, // Not available in this endpoint
            high52Week: parseFloat(quote['03. high']),
            low52Week: parseFloat(quote['04. low']),
        };
    },
});
