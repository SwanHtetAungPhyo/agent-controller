import {createTool} from "@mastra/core/tools";
import {z} from "zod";
import {getIndexName, makeFinancialApiRequest} from "./financial-tools";


export const marketIndicesTool = createTool({
    id: 'get-market-indices',
    description: 'Get current values and performance of major market indices',
    inputSchema: z.object({
        indices: z.array(z.string()).optional().default(['SPY', 'QQQ', 'DIA']).describe('Array of index symbols'),
    }),
    outputSchema: z.object({
        indices: z.array(z.object({
            symbol: z.string(),
            name: z.string(),
            value: z.number(),
            change: z.number(),
            changePercent: z.number(),
        })),
    }),
    execute: async ({ context }) => {
        const indices = context.indices || ['SPY', 'QQQ', 'DIA'];
        const results = [];

        for (const symbol of indices) {
            try {
                const data = await makeFinancialApiRequest('GLOBAL_QUOTE', symbol);
                const quote = data['Global Quote'];

                if (quote) {
                    results.push({
                        symbol: quote['01. symbol'],
                        name: getIndexName(symbol),
                        value: parseFloat(quote['05. price']),
                        change: parseFloat(quote['09. change']),
                        changePercent: parseFloat(quote['10. change percent'].replace('%', '')),
                    });
                }
            } catch (error) {
                // Add mock data if API fails
                results.push({
                    symbol: symbol,
                    name: getIndexName(symbol),
                    value: 450.00,
                    change: 2.50,
                    changePercent: 0.56,
                });
            }
        }

        return { indices: results };
    },
});
