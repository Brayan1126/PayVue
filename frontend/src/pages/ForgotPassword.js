import React, { useState } from 'react';
import { Link } from 'react-router-dom';

function ForgotPassword() {
  const [email, setEmail] = useState('');
  const [sent, setSent] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    setSent(true);
  };

  return (
    <div className="auth-container">
      <h1 className="auth-title">PayVue APP</h1>
      
      <div className="auth-card">
        <h2>Recuperar contraseña</h2>
        <p>Ingrese su Email de registro</p>
        
        {!sent ? (
          <form className="auth-form" onSubmit={handleSubmit}>
            <input
              type="email"
              className="form-input"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <button type="submit" className="btn-primary">
              Recuperar Contraseña
            </button>
          </form>
        ) : (
          <div style={{ textAlign: 'center', padding: '20px' }}>
            <div style={{ marginBottom: '20px', color: '#22c55e' }}>
              <svg viewBox="0 0 24 24" width="60" height="60" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
            </div>
            <p>Se ha enviado un enlace de recuperación a tu correo electrónico.</p>
          </div>
        )}

        <div className="auth-links">
          <Link to="/">Volver al inicio de sesión</Link>
        </div>
      </div>
    </div>
  );
}

export default ForgotPassword;

