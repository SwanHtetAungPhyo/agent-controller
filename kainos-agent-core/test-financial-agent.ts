import { config } from 'dotenv';
import { mastra } from './src/mastra/index.js';

// Load environment variables
config();

async function testFinancialAgent() {
  console.log('🧪 Testing Financial Agent...\n');

  try {
    // Get the financial agent
    const agent = mastra.getAgent('financialAgent');
    
    if (!agent) {
      console.error('❌ Financial agent not found!');
      return;
    }

    console.log('✅ Financial agent loaded successfully\n');

    // Test 1: Get stock price
    console.log('📊 Test 1: Getting stock price for AAPL...');
    const priceResponse = await agent.generate([
      {
        role: 'user',
        content: 'Get the current stock price for AAPL'
      }
    ]);

    console.log('Response:', priceResponse.text);
    console.log('\n---\n');

    // Test 2: Get company information
    console.log('🏢 Test 2: Getting company information for Microsoft...');
    const companyResponse = await agent.generate([
      {
        role: 'user',
        content: 'Tell me about Microsoft (MSFT) company information'
      }
    ]);

    console.log('Response:', companyResponse.text);
    console.log('\n---\n');

    console.log('✅ All tests completed successfully!');
  } catch (error) {
    console.error('❌ Error during testing:', error);
    throw error;
  }
}

// Run the test
testFinancialAgent().catch(console.error);
