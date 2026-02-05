const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export function getAuthHeaders(): HeadersInit {
  const auth = localStorage.getItem('auth');
  if (auth) {
    return {
      'Authorization': `Basic ${auth}`,
      'Content-Type': 'application/json',
    };
  }
  return {
    'Content-Type': 'application/json',
  };
}

export default API_URL;
