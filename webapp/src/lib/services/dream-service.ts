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
    console.log('Fetching all dreams...');
    const dreams = await this.get<Dream[]>('/api/dreams');
    console.log('Received dreams:', dreams);
    return dreams;
  }

  async update(id: number, dream: string): Promise<Dream> {
    console.log('Sending update request:', { id, dream });
    const response = await this.put<Dream>(`/api/dreams/${id}`, { dream });
    console.log('Update response:', response);
    return response;
  }

  async removeDelete(id: number): Promise<void> {
    await this.delete<void>(`/api/dreams/${id}`);
  }
} 