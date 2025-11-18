import { createTool } from '@mastra/core/tools';
import { z } from 'zod';

export function getIndexName(symbol: string): string {
    const indexNames: Record<string, string> = {
        'SPY': 'S&P 500 ETF',
        'QQQ': 'NASDAQ-100 ETF',
        'DIA': 'Dow Jones Industrial Average ETF',
        'VTI': 'Total Stock Market ETF',
        'IWM': 'Russell 2000 ETF',
    };
    return indexNames[symbol] || symbol;
}

const FINANCIAL_API_BASE = 'https://www.alphavantage.co/query';

export async function makeFinancialApiRequest(functionType: string, symbol?: string, params: Record<string, string> = {}) {
    const url = new URL(FINANCIAL_API_BASE);
    url.searchParams.append('function', functionType);
    url.searchParams.append('apikey', process.env.FINANCIAL_DATASETS_API_KEY || 'demo');

    if (symbol) {
        url.searchParams.append('symbol', symbol);
    }

    Object.entries(params).forEach(([key, value]) => {
        url.searchParams.append(key, value);
    });

    const response = await fetch(url.toString());

    if (!response.ok) {
        throw new Error(`Financial API error: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();

    if (data['Error Message'] || data['Note']) {
        throw new Error(data['Error Message'] || data['Note'] || 'API limit reached');
    }

    return data;
}

