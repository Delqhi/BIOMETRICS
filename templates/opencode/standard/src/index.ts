import { config } from 'dotenv';

config();

const PORT = parseInt(process.env.PORT || '50001', 10);
const NODE_ENV = process.env.NODE_ENV || 'development';

async function main() {
  console.log(`🚀 Starting {{PROJECT_NAME}} on port ${PORT}`);
  console.log(`📦 Environment: ${NODE_ENV}`);

  // TODO: Add your application logic here
  console.log('✅ Application started successfully');
}

main().catch((error) => {
  console.error('❌ Failed to start application:', error);
  process.exit(1);
});
