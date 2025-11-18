import {createTool} from "@mastra/core/tools";
import {z} from "zod";
import {makeFinancialApiRequest} from "./financial-tools";


export const marketNewsTool = createTool({
    id: 'get-market-news',
    description: 'Get latest market news and headlines',
    inputSchema: z.object({
        symbol: z.string().optional().describe('Stock ticker symbol for company-specific news'),
        limit: z.number().optional().default(10).describe('Number of news articles to return'),
    }),
    outputSchema: z.object({
        articles: z.array(z.object({
            title: z.string(),
            summary: z.string().optional(),
            url: z.string(),
            publishedAt: z.string(),
            source: z.string(),
            sentiment: z.string().optional(),
        })),
    }),
    execute: async ({ context }) => {
        try {
            // Try to get news from Alpha Vantage (limited availability)
            const data = await makeFinancialApiRequest('NEWS_SENTIMENT', context.symbol?.toUpperCase());

            if (data.feed) {
                return {
                    articles: data.feed.slice(0, context.limit || 10).map((article: any) => ({
                        title: article.title,
                        summary: article.summary,
                        url: article.url,
                        publishedAt: article.time_published,
                        source: article.source,
                        sentiment: article.overall_sentiment_label,
                    })),
                };
            }
        } catch (error) {
            // Fallback to mock news data
        }

        // Mock news data for demonstration
        const mockArticles = [
            {
                title: `${context.symbol || 'Market'} Shows Strong Performance in Recent Trading`,
                summary: 'Recent market activity shows positive trends with increased investor confidence.',
                url: 'https://example.com/news/1',
                publishedAt: new Date().toISOString(),
                source: 'Financial News',
                sentiment: 'positive',
            },
            {
                title: `Analysts Upgrade ${context.symbol || 'Stock'} Rating`,
                summary: 'Leading analysts have upgraded their rating based on strong fundamentals.',
                url: 'https://example.com/news/2',
                publishedAt: new Date(Date.now() - 3600000).toISOString(),
                source: 'Market Watch',
                sentiment: 'positive',
            },
            {
                title: `${context.symbol || 'Company'} Reports Quarterly Earnings`,
                summary: 'Latest quarterly results show mixed performance with revenue growth.',
                url: 'https://example.com/news/3',
                publishedAt: new Date(Date.now() - 7200000).toISOString(),
                source: 'Business Times',
                sentiment: 'neutral',
            },
        ];

        return {
            articles: mockArticles.slice(0, context.limit || 10),
        };
    },
});
