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
    return await this.get<Dream[]>('/api/dreams');
  }

  async getById(id: string | number): Promise<Dream> {
    try {
      const response = await fetch(`${this.baseUrl}/api/dreams/${id}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Error fetching dream:', { status: response.status, error: errorText });
        throw new Error(`Failed to fetch dream: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Error in getById:', error);
      throw error;
    }
  }

  async generateImage(id: string | number, onProgress?: (status: { position?: number; message: string }) => void): Promise<{ message: string; position?: number }> {
    try {
      // 1. Start the image generation process
      const url = `${process.env.NEXT_PUBLIC_API_URL || ''}/api/dreams/${id.toString()}/generate-image`;
      
      const startResponse = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      }).catch(error => {
        console.error('Network error during image generation:', error);
        throw new Error(`Network error: ${error.message}`);
      });

      if (!startResponse.ok) {
        const error = await startResponse.json().catch(() => ({}));
        throw new Error(error.message || 'Failed to start image generation');
      }

      const { queuePosition, message } = await startResponse.json();
      
      if (onProgress) {
        if (queuePosition > 0) {
          onProgress({
            position: queuePosition,
            message: `Your dream is in the queue (Position: ${queuePosition}). It will start generating soon...`
          });
        } else {
          onProgress({
            position: 0,
            message: 'Starting image generation. This may take a few minutes...'
          });
        }
      }

      // 2. Poll for completion
      const pollInterval = 10000; // 10 seconds
      const maxAttempts = 300; // 50 minutes max (300 * 10000ms = 50 minutes)
      
      for (let attempt = 0; attempt < maxAttempts; attempt++) {
        await new Promise(resolve => setTimeout(resolve, pollInterval));
        
        try {
          const statusResponse = await fetch(`${process.env.NEXT_PUBLIC_API_URL || ''}/api/dreams/${id}/status`);
          
          if (statusResponse.status === 200) {
            const dream = await statusResponse.json();
            if (dream.image_url) {  // Fixed field name to match backend
              return dream;
            }
          } else if (statusResponse.status === 202) {
            // Still processing
            const status = await statusResponse.json();

            if (onProgress) {
              onProgress({
                position: status.queuePosition || queuePosition,
                message: status.queuePosition > 0 
                  ? `Your dream is in the queue (Position: ${status.queuePosition}). It will start generating soon...`
                  : 'Your dream is being generated. This may take a few minutes...'
              });
            }
            continue;
          } else {
            // Handle other status codes
            const error = await statusResponse.json().catch(() => ({}));
            throw new Error(error.message || 'Error checking image generation status');
          }
        } catch (error) {
          console.error('Error polling image status:', error);
          // Continue polling on network errors
          if (onProgress) {
            onProgress({
              position: queuePosition,
              message: 'Connection issue, retrying...'
            });
          }
        }
      }

      throw new Error('Image generation timed out');
    } catch (error) {
      console.error('Error in generateImage:', error);
      throw error;
    }
  }

  async update(id: string | number, dream: string): Promise<Dream> {
    return await this.put<Dream>(`/api/dreams/${id}`, { dream });
  }

  async removeDelete(id: string | number): Promise<void> {
    try {
      await this.delete(`/api/dreams/${id}`);
    } catch (error) {
      console.error('Error deleting dream:', error);
      throw error;
    }
  }

  async checkImageStatus(id: string): Promise<{ 
    isGenerating: boolean; 
    position?: number; 
    message: string;
    dream?: any;
    error?: string;
    shouldRetry?: boolean;
  }> {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || ''}/api/dreams/${id}/status`);
      
      if (response.status === 200) {
        // Image is ready
        const dream = await response.json();
        return { 
          isGenerating: false, 
          message: 'Image generated successfully',
          dream,
          shouldRetry: false
        };
      } else if (response.status === 202) {
        // Still processing
        const { queuePosition, message } = await response.json();
        return { 
          isGenerating: true, 
          position: queuePosition, 
          message: message || 'Generating image...',
          shouldRetry: true
        };
      } else if (response.status === 204) {
        // No image generation in progress
        return { 
          isGenerating: false, 
          message: 'No image generation in progress',
          shouldRetry: false
        };
      } else {
        // Handle other status codes
        const error = await response.json().catch(() => ({}));
        const errorMessage = error.message || 'Error checking image status';
        return {
          isGenerating: false,
          message: errorMessage,
          error: errorMessage,
          shouldRetry: response.status >= 500 // Only retry on server errors
        };
      }
    } catch (error) {
      console.error('Error checking image status:', error);
      return {
        isGenerating: false,
        message: 'Network error while checking image status',
        error: error instanceof Error ? error.message : 'Unknown error',
        shouldRetry: true // Network errors might be temporary
      };
    }
  }
} 