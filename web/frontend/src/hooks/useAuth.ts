import { useState, useEffect } from 'react';

interface AuthState {
  user: User | null;
  loading: boolean;
  error: string | null;
}

interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'trader' | 'viewer';
}

// This is a mock implementation - would be replaced with actual Firebase/Auth0 implementation
export const useAuth = (): AuthState & {
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  logout: () => Promise<void>;
} => {
  const [state, setState] = useState<AuthState>({
    user: null,
    loading: true,
    error: null,
  });

  useEffect(() => {
    // Check if user is already logged in
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setState({
        user: JSON.parse(storedUser),
        loading: false,
        error: null,
      });
    } else {
      setState(prev => ({ ...prev, loading: false }));
    }
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      // Mock API call - would be replaced with actual authentication
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Mock successful login
      const user: User = {
        id: '1',
        email,
        name: email.split('@')[0],
        role: 'trader',
      };
      
      localStorage.setItem('user', JSON.stringify(user));
      setState({ user, loading: false, error: null });
    } catch (error) {
      setState(prev => ({ 
        ...prev, 
        loading: false, 
        error: error instanceof Error ? error.message : 'Failed to login' 
      }));
    }
  };

  const register = async (email: string, password: string, name: string) => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      // Mock API call - would be replaced with actual registration
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Mock successful registration
      const user: User = {
        id: '1',
        email,
        name,
        role: 'trader',
      };
      
      localStorage.setItem('user', JSON.stringify(user));
      setState({ user, loading: false, error: null });
    } catch (error) {
      setState(prev => ({ 
        ...prev, 
        loading: false, 
        error: error instanceof Error ? error.message : 'Failed to register' 
      }));
    }
  };

  const logout = async () => {
    try {
      setState(prev => ({ ...prev, loading: true }));
      
      // Mock API call - would be replaced with actual logout
      await new Promise(resolve => setTimeout(resolve, 500));
      
      localStorage.removeItem('user');
      setState({ user: null, loading: false, error: null });
    } catch (error) {
      setState(prev => ({ 
        ...prev, 
        loading: false, 
        error: error instanceof Error ? error.message : 'Failed to logout' 
      }));
    }
  };

  return {
    ...state,
    login,
    register,
    logout,
  };
};
