import axios from 'axios';
import { mockAdapter } from './mockApi';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Crear instancia de axios con configuración base
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  },
  // Sin backend: REACT_APP_MOCK_API=true usa datos falsos en memoria
  ...(process.env.REACT_APP_MOCK_API === 'true' ? { adapter: mockAdapter } : {})
});

// Interceptor para añadir el user_id a todas las peticiones
api.interceptors.request.use(
  (config) => {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    if (user.user_id) {
      config.headers['X-User-ID'] = user.user_id;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Función helper para obtener el user_id actual
export const getCurrentUserId = () => {
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  return user.user_id || 0;
};

// Función helper para limpiar datos del usuario
export const clearUserData = () => {
  localStorage.removeItem('user');
};

// Función helper para guardar datos del usuario
export const setUserData = (userData) => {
  localStorage.setItem('user', JSON.stringify(userData));
};

// Función helper para obtener datos del usuario
export const getUserData = () => {
  return JSON.parse(localStorage.getItem('user') || '{}');
};

export { api };
export default API_URL;
