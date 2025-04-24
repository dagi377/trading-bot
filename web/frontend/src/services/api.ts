import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export const api = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await axios.post(`${API_URL}/login`, credentials);
    return response.data;
  },

  // Add token to all subsequent requests
  setAuthToken: (token: string) => {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  },

  // Remove token
  clearAuthToken: () => {
    delete axios.defaults.headers.common['Authorization'];
  }
}; 
