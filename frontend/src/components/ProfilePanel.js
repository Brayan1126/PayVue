import React, { useState, useEffect } from 'react';

function ProfilePanel({ isOpen, onClose }) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');

  useEffect(() => {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    setUsername(user.username || user.email?.split('@')[0] || '');
    setEmail(user.email || '');
  }, [isOpen]);

  const handleSave = () => {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    user.username = username;
    user.email = email;
    localStorage.setItem('user', JSON.stringify(user));
    onClose();
  };

  if (!isOpen) return null;

  return (
    <>
      <div className={`profile-overlay ${isOpen ? 'open' : ''}`} onClick={onClose} />
      <div className={`profile-panel ${isOpen ? 'open' : ''}`}>
        <button onClick={onClose} style={{
          position: 'absolute',
          top: '20px',
          right: '20px',
          background: 'none',
          border: 'none',
          fontSize: '1.8rem',
          cursor: 'pointer',
          color: '#666',
          lineHeight: 1
        }}>×</button>
        
        <h2 style={{ textAlign: 'center', marginBottom: '30px' }}>Editar Perfil</h2>
        
        <img 
          src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150" 
          alt="Profile" 
          className="profile-avatar"
        />
        
        <p className="update-photo">Actualizar Foto</p>
        
        <div className="modal-form">
          <div className="form-group">
            <label>Nombre de Usuario</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Tu nombre de usuario"
            />
          </div>
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="tu@email.com"
            />
          </div>
          <button className="btn-modal" onClick={handleSave}>Guardar</button>
        </div>
      </div>
    </>
  );
}

export default ProfilePanel;
