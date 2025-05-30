export abstract class Api {
  protected baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  protected async fetchWithError(url: string, options?: RequestInit) {
    console.log('Fetching:', url, options);
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
      });

      console.log('Response status:', response.status);
      const responseText = await response.text();
      console.log('Response body:', responseText);

      if (!response.ok) {
        let errorMessage = 'An error occurred';
        try {
          const error = JSON.parse(responseText);
          errorMessage = error.message || errorMessage;
        } catch (e) {
          errorMessage = responseText || errorMessage;
        }
        console.error('API Error:', errorMessage);
        throw new Error(errorMessage);
      }

      return new Response(responseText, {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
      });
    } catch (error) {
      console.error('Fetch Error:', error);
      throw error;
    }
  }

  protected async get<T>(endpoint: string): Promise<T> {
    console.log('GET request to:', `${this.baseUrl}${endpoint}`);
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`);
    const data = await response.json();
    console.log('GET response:', data);
    return data;
  }

  protected async post<T>(endpoint: string, data: unknown): Promise<T> {
    console.log('POST request to:', `${this.baseUrl}${endpoint}`, 'with data:', data);
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
    const responseData = await response.json();
    console.log('POST response:', responseData);
    return responseData;
  }

  protected async put<T>(endpoint: string, data: unknown): Promise<T> {
    console.log('PUT request to:', `${this.baseUrl}${endpoint}`, 'with data:', data);
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    if (response.headers.get('content-type')?.includes('application/json')) {
      const responseData = await response.json();
      console.log('PUT response:', responseData);
      return responseData;
    }
    return undefined as T;
  }

  protected async delete<T>(endpoint: string): Promise<T> {
    console.log('DELETE request to:', `${this.baseUrl}${endpoint}`);
    const response = await this.fetchWithError(`${this.baseUrl}${endpoint}`, {
      method: 'DELETE',
    });
    if (response.headers.get('content-type')?.includes('application/json')) {
      const responseData = await response.json();
      console.log('DELETE response:', responseData);
      return responseData;
    }
    return undefined as T;
  }
}