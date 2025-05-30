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

  async getById(id: string): Promise<Dream> {
    console.log('Fetching dream by id:', id);
    const dream = await this.get<Dream>(`/api/dreams/${id}`);
    console.log('Received dream:', dream);
    return dream;
  }

  async generateImage(id: string): Promise<Dream> {
    console.log('Generating image for dream:', id);
    const dream = await this.post<Dream>(`/api/dreams/${id}/generate-image`, {});
    console.log('Image generation response:', dream);
    return dream;
  }

  async update(id: number, dream: string): Promise<Dream> {
    console.log('Sending update request:', { id, dream });
    const response = await this.put<Dream>(`/api/dreams/${id}`, { dream });
    console.log('Update response:', response);
    return response;
  }

  async removeDelete(id: number): Promise<void> {
    console.log('Deleting dream:', id);
    try {
      const result = await this.delete<void>(`/api/dreams/${id}`);
      console.log('Dream deleted successfully');
      return result;
    } catch (error) {
      console.error('Error deleting dream:', error);
      throw error;
    }
  }
} 