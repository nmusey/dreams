export abstract class Api {
  protected baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  protected async fetchWithError(url: string, options?: RequestInit) {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
      });

      // Handle 204 No Content as a special case - return the original response
      if (response.status === 204) {
        return response;
      }

      const responseText = await response.text();

      if (!response.ok) {
        let errorMessage = 'An error occurred';
        try {
          const error = JSON.parse(responseText);
          errorMessage = error.message || errorMessage;
        } catch (e) {
          errorMessage = responseText || errorMessage;
        }
        console.error('API Error:', { url, status: response.status, message: errorMessage });
        throw new Error(errorMessage);
      }

      return new Response(responseText, {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
      });
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';
      console.error('Network Error:', { url, error: errorMessage });
      throw error;
    }
  }

  protected async get<T>(endpoint: string): Promise<T> {
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`);
    return response.json();
  }

  protected async post<T>(endpoint: string, data: unknown): Promise<T> {
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.json();
  }

  protected async put<T>(endpoint: string, data: unknown): Promise<T> {
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    if (response.status === 204) {
      return undefined as T;
    }
    if (response.headers.get('content-type')?.includes('application/json')) {
      return response.json();
    }
    return undefined as T;
  }

  protected async delete<T>(endpoint: string): Promise<T> {
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`, {
      method: 'DELETE',
    });
    if (response.status === 204) {
      return undefined as T;
    }
    if (response.headers.get('content-type')?.includes('application/json')) {
      return response.json();
    }
    return undefined as T;
  }
}