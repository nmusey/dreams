import AppDataSource from '../db-config';

export async function initDb() {
  try {
    await AppDataSource.initialize();
    console.log('Connection to the database has been established.');
  } catch (err) {
    console.error('Error connecting to database: ', err);
  }
}
