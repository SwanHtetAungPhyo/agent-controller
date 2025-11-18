import { config } from 'dotenv';
import { mastra } from './src/mastra/index.js';

// Load environment variables
config();

async function testAllTools() {
  console.log('🧪 Testing All Financial Tools...\n');
  console.log('='.repeat(60));

  try {
    const agent = mastra.getAgent('financialAgent');
    
    if (!agent) {
      console.error('❌ Financial agent not found!');
      return;
    }

    console.log('✅ Financial agent loaded successfully\n');

    // Test 1: Stock Price Tool
    console.log('\n📊 Test 1: Stock Price Tool (AAPL)');
    console.log('-'.repeat(60));
    const test1 = await agent.generate([
      { role: 'user', content: 'Get the current stock price for AAPL' }
    ]);
    console.log('✅ Result:', test1.text.substring(0, 200) + '...\n');

    // Test 2: Financial Metrics Tool
    console.log('\n📈 Test 2: Financial Metrics Tool (TSLA)');
    console.log('-'.repeat(60));
    const test2 = await agent.generate([
      { role: 'user', content: 'Show me the financial metrics for Tesla (TSLA)' }
    ]);
    console.log('✅ Result:', test2.text.substring(0, 200) + '...\n');

    // Test 3: Company Information Tool
    console.log('\n🏢 Test 3: Company Information Tool (GOOGL)');
    console.log('-'.repeat(60));
    const test3 = await agent.generate([
      { role: 'user', content: 'Tell me about Google (GOOGL) company information' }
    ]);
    console.log('✅ Result:', test3.text.substring(0, 200) + '...\n');

    // Test 4: Market News Tool
    console.log('\n📰 Test 4: Market News Tool (NVDA)');
    console.log('-'.repeat(60));
    const test4 = await agent.generate([
      { role: 'user', content: 'Get the latest market news for NVIDIA (NVDA)' }
    ]);
    console.log('✅ Result:', test4.text.substring(0, 200) + '...\n');

    // Test 5: Market Indices Tool
    console.log('\n📉 Test 5: Market Indices Tool');
    console.log('-'.repeat(60));
    const test5 = await agent.generate([
      { role: 'user', content: 'Show me the current market indices performance' }
    ]);
    console.log('✅ Result:', test5.text.substring(0, 200) + '...\n');

    // Test 6: Portfolio Analysis Tool
    console.log('\n💼 Test 6: Portfolio Analysis Tool');
    console.log('-'.repeat(60));
    const test6 = await agent.generate([
      { role: 'user', content: 'Analyze a portfolio with AAPL (30%), GOOGL (25%), MSFT (25%), and TSLA (20%)' }
    ]);
    console.log('✅ Result:', test6.text.substring(0, 200) + '...\n');

    // Test 7: Economic Indicators Tool
    console.log('\n🌍 Test 7: Economic Indicators Tool');
    console.log('-'.repeat(60));
    const test7 = await agent.generate([
      { role: 'user', content: 'Show me the current economic indicators' }
    ]);
    console.log('✅ Result:', test7.text.substring(0, 200) + '...\n');

    console.log('\n' + '='.repeat(60));
    console.log('✅ All 7 tools tested successfully!');
    console.log('='.repeat(60));

  } catch (error) {
    console.error('❌ Error during testing:', error);
    throw error;
  }
}

// Run the test
testAllTools().catch(console.error);
