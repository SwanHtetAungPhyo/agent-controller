import {createTool} from "@mastra/core/tools";
import {z} from "zod";

export const economicIndicatorsTool = createTool({
    id: 'get-economic-indicators',
    description: 'Get key economic indicators and market data',
    inputSchema: z.object({
        indicators: z.array(z.string()).optional().default(['GDP', 'INFLATION', 'UNEMPLOYMENT', 'INTEREST_RATES']).describe('Economic indicators to fetch'),
    }),
    outputSchema: z.object({
        indicators: z.array(z.object({
            name: z.string(),
            value: z.number(),
            unit: z.string(),
            change: z.number().optional(),
            lastUpdated: z.string(),
        })),
    }),
    execute: async ({ context }) => {
        // Mock economic indicators data (Alpha Vantage free tier doesn't include this)
        const mockIndicators = [
            {
                name: 'GDP Growth Rate',
                value: 2.1,
                unit: '%',
                change: 0.1,
                lastUpdated: new Date().toISOString(),
            },
            {
                name: 'Inflation Rate',
                value: 3.2,
                unit: '%',
                change: -0.2,
                lastUpdated: new Date().toISOString(),
            },
            {
                name: 'Unemployment Rate',
                value: 3.8,
                unit: '%',
                change: -0.1,
                lastUpdated: new Date().toISOString(),
            },
            {
                name: 'Federal Funds Rate',
                value: 5.25,
                unit: '%',
                change: 0.0,
                lastUpdated: new Date().toISOString(),
            },
        ];

        const requestedIndicators = context.indicators || ['GDP', 'INFLATION', 'UNEMPLOYMENT', 'INTEREST_RATES'];

        return {
            indicators: mockIndicators.filter(indicator =>
                requestedIndicators.some(req =>
                    indicator.name.toLowerCase().includes(req.toLowerCase())
                )
            ),
        };
    },
});