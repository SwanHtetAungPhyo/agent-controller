import { config } from 'dotenv';
import { secFilingsTool } from './src/mastra/tools/sec_filings_tool.js';
import { insiderTradesTool } from './src/mastra/tools/insider_trades_tool.js';

config();

async function testNewTools() {
  console.log('🧪 Testing New Financial Tools...\n');

  try {
    // Test SEC Filings Tool
    console.log('📄 Test 1: SEC Filings for Apple (AAPL)...');
    const secResult = await secFilingsTool.execute({
      context: {
        ticker: 'AAPL',
        filing_type: '10-K',
      },
    });
    console.log('SEC Filings Result:', JSON.stringify(secResult, null, 2));
    console.log('\n---\n');

    // Test Insider Trades Tool
    console.log('💼 Test 2: Insider Trades for Tesla (TSLA)...');
    const insiderResult = await insiderTradesTool.execute({
      context: {
        ticker: 'TSLA',
        limit: 5,
      },
    });
    console.log('Insider Trades Result:', JSON.stringify(insiderResult, null, 2));
    console.log('\n---\n');

    console.log('✅ All new tools tested successfully!');
  } catch (error) {
    console.error('❌ Error during testing:', error);
    throw error;
  }
}

testNewTools().catch(console.error);
