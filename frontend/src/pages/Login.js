import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api, setUserData, clearUserData } from '../config/api';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    
    // Limpiar datos del usuario anterior
    clearUserData();
    
    try {
      const response = await api.post('/auth/login', { email, password });
      
      // Guardar datos del usuario incluyendo user_id
      setUserData({
        user_id: response.data.user_id,
        email: response.data.email || email
      });
      
      navigate('/dashboard');
    } catch (error) {
      alert(error.response?.data?.message || 'Error al iniciar sesión');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <h1 className="auth-title">PayVue APP</h1>
      
      <div className="auth-card">
        <h2>Iniciar sesión</h2>
        <p>Ingrese su usuario y contraseña</p>
        
        <form className="auth-form" onSubmit={handleLogin}>
          <input
            type="email"
            className="form-input"
            placeholder="Usuario"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            className="form-input"
            placeholder="Contraseña"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Cargando...' : 'Iniciar Sesión'}
          </button>
        </form>

        <div className="divider">
          <span>O inicia sesión con</span>
        </div>

        <button className="btn-google">
          <svg viewBox="0 0 24 24" width="20" height="20">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          Google
        </button>

        <div className="auth-links">
          <Link to="/forgot-password">¿Olvidaste tu contraseña?</Link>
          <p>¿No tienes cuenta? <Link to="/register"><strong>Registrate aquí</strong></Link></p>
        </div>
      </div>
    </div>
  );
}

export default Login;
