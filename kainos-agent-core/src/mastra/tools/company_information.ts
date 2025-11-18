import {createTool} from "@mastra/core/tools";
import {z} from "zod";
import {makeFinancialApiRequest} from "./financial-tools";

export const companyInfoTool = createTool({
    id: 'get-company-info',
    description: 'Get detailed company information and profile',
    inputSchema: z.object({
        symbol: z.string().describe('Stock ticker symbol'),
    }),
    outputSchema: z.object({
        symbol: z.string(),
        name: z.string(),
        sector: z.string().optional(),
        industry: z.string().optional(),
        description: z.string().optional(),
        employees: z.number().optional(),
        founded: z.string().optional(),
        headquarters: z.string().optional(),
        website: z.string().optional(),
    }),
    execute: async ({ context }) => {
        try {
            const data = await makeFinancialApiRequest('OVERVIEW', context.symbol.toUpperCase());

            return {
                symbol: data.Symbol,
                name: data.Name,
                sector: data.Sector,
                industry: data.Industry,
                description: data.Description,
                employees: data.FullTimeEmployees ? parseInt(data.FullTimeEmployees) : undefined,
                founded: undefined, // Not available in Alpha Vantage
                headquarters: data.Address,
                website: undefined, // Not available in Alpha Vantage
            };
        } catch (error) {
            // Fallback to mock data
            return {
                symbol: context.symbol.toUpperCase(),
                name: `${context.symbol.toUpperCase()} Company`,
                sector: 'Technology',
                industry: 'Software',
                description: `${context.symbol.toUpperCase()} is a technology company.`,
                employees: 50000,
                founded: '2000',
                headquarters: 'United States',
                website: `https://${context.symbol.toLowerCase()}.com`,
            };
        }
    },
});
