const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export interface User {
  id: number;
  email: string;
  name?: string;
  role: 'user' | 'admin' | 'superadmin';
  created_at: string;
  updated_at: string;
}

export interface Restaurant {
  id: number;
  name: string;
  description?: string;
  address?: string;
  phone?: string;
  email?: string;
  image_url?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

class ApiClient {
  private getAuthToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('access_token');
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = this.getAuthToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const url = `${API_URL}${endpoint}`;
    console.log('API Request:', url, options.method || 'GET'); // Debug log

    const response = await fetch(url, {
      ...options,
      headers: headers as HeadersInit,
    });

    console.log('API Response:', response.status, response.statusText); // Debug log

    if (!response.ok) {
      const errorText = await response.text();
      console.error('API Error:', errorText); // Debug log
      let error;
      try {
        error = JSON.parse(errorText);
      } catch {
        error = { error: errorText || `HTTP ${response.status}: ${response.statusText}` };
      }
      throw new Error(error.error || 'Request failed');
    }

    return response.json();
  }

  // Auth
  async register(email: string, password: string, name?: string): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
    });
  }

  async login(email: string, password: string): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  async getMe(): Promise<{ user: User }> {
    return this.request<{ user: User }>('/me');
  }

  // Restaurants
  async getRestaurants(page = 1, pageSize = 10): Promise<{
    data: Restaurant[];
    total: number;
    page: number;
    page_size: number;
  }> {
    return this.request(`/restaurants?page=${page}&page_size=${pageSize}`);
  }

  async getRestaurant(id: number): Promise<Restaurant> {
    return this.request<Restaurant>(`/restaurants/${id}`);
  }

  async createRestaurant(data: Partial<Restaurant>): Promise<Restaurant> {
    return this.request<Restaurant>('/restaurants', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateRestaurant(id: number, data: Partial<Restaurant>): Promise<Restaurant> {
    return this.request<Restaurant>(`/restaurants/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteRestaurant(id: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/restaurants/${id}`, {
      method: 'DELETE',
    });
  }
}

export const api = new ApiClient();
