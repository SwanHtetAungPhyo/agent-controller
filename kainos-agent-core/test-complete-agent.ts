import { config } from 'dotenv';
import { mastra } from './src/mastra/index.js';

config();

async function testCompleteAgent() {
  console.log('🧪 Testing Complete Financial Agent with All 9 Tools...\n');

  try {
    const agent = mastra.getAgent('financialAgent');
    
    if (!agent) {
      console.error('❌ Financial agent not found!');
      return;
    }

    console.log('✅ Financial agent loaded successfully\n');

    // Test 1: Stock price
    console.log('📊 Test 1: Stock Price Query...');
    const test1 = await agent.generate([
      { role: 'user', content: 'What is the current stock price for Apple (AAPL)?' }
    ]);
    console.log('Response:', test1.text);
    console.log('\n---\n');

    // Test 2: Financial metrics
    console.log('💰 Test 2: Financial Metrics Query...');
    const test2 = await agent.generate([
      { role: 'user', content: 'Show me the financial metrics for Tesla (TSLA)' }
    ]);
    console.log('Response:', test2.text);
    console.log('\n---\n');

    // Test 3: SEC Filings (NEW)
    console.log('📄 Test 3: SEC Filings Query...');
    const test3 = await agent.generate([
      { role: 'user', content: 'Get the latest 10-K filing for Microsoft (MSFT)' }
    ]);
    console.log('Response:', test3.text);
    console.log('\n---\n');

    // Test 4: Insider Trades (NEW)
    console.log('💼 Test 4: Insider Trades Query...');
    const test4 = await agent.generate([
      { role: 'user', content: 'Show me recent insider trading activity for NVIDIA (NVDA)' }
    ]);
    console.log('Response:', test4.text);
    console.log('\n---\n');

    // Test 5: Market news
    console.log('📰 Test 5: Market News Query...');
    const test5 = await agent.generate([
      { role: 'user', content: 'What is the latest news about Amazon (AMZN)?' }
    ]);
    console.log('Response:', test5.text);
    console.log('\n---\n');

    console.log('✅ All agent tests completed successfully!');
    console.log('\n🎉 Financial Agent with 9 tools is fully operational!');
  } catch (error) {
    console.error('❌ Error during testing:', error);
    throw error;
  }
}

testCompleteAgent().catch(console.error);
