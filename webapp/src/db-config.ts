import path from 'path';
import { DataSource } from 'typeorm';

export const AppDataSource = new DataSource({
  type: 'postgres',
  host: process.env.DB_HOST || 'localhost',
  port: Number(process.env.DB_PORT || 5432),
  username: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
  entities: [path.resolve(__dirname, '../models/*.ts')], 
  synchronize: process.env.NODE_ENV == "development", 
  logging: false,
});

export default AppDataSource;

