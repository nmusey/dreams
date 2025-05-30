import { Dream } from '@/lib/types/dream';
import { Api } from './api';

export class DreamService extends Api {
  private static instance: DreamService;

  static getInstance(): DreamService {
    if (!DreamService.instance) {
      DreamService.instance = new DreamService();
    }
    return DreamService.instance;
  }

  async create(content: string): Promise<void> {
    await this.post<void>('/api/dreams', { dream: content });
  }

  async getAll(): Promise<Dream[]> {
    return this.get<Dream[]>('/api/dreams');
  }
} 