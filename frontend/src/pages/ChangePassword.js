import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Toast from '../components/Toast';

function ChangePassword() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [toast, setToast] = useState({ show: false, message: '', type: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    setLoading(true);

    if (formData.newPassword !== formData.confirmPassword) {
      setToast({ show: true, message: 'Las contraseñas no coinciden', type: 'error' });
      setTimeout(() => setToast({ show: false, message: '', type: '' }), 3000);
      setLoading(false);
      return;
    }

    if (formData.newPassword.length < 6) {
      setToast({ show: true, message: 'La contraseña debe tener al menos 6 caracteres', type: 'error' });
      setTimeout(() => setToast({ show: false, message: '', type: '' }), 3000);
      setLoading(false);
      return;
    }

    // Simular guardado
    setTimeout(() => {
      setToast({ show: true, message: '¡Contraseña actualizada con éxito!', type: 'success' });
      setFormData({ currentPassword: '', newPassword: '', confirmPassword: '' });
      setLoading(false);
      
      setTimeout(() => {
        setToast({ show: false, message: '', type: '' });
        navigate('/dashboard');
      }, 2000);
    }, 1000);
  };

  const handleClose = () => {
    navigate('/dashboard');
  };

  return (
    <div className="modal-page">
      <div className="modal-card" style={{ position: 'relative' }}>
        <button 
          onClick={handleClose}
          style={{
            position: 'absolute',
            top: '20px',
            right: '25px',
            background: 'none',
            border: 'none',
            fontSize: '1.8rem',
            cursor: 'pointer',
            color: '#666',
            lineHeight: 1
          }}
        >×</button>
        
        <h2>Cambiar Contraseña</h2>
        
        <form className="modal-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Contraseña Actual</label>
            <input
              type="password"
              placeholder="Ingresa tu contraseña actual"
              value={formData.currentPassword}
              onChange={(e) => setFormData({ ...formData, currentPassword: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label>Nueva Contraseña</label>
            <input
              type="password"
              placeholder="Mínimo 6 caracteres"
              value={formData.newPassword}
              onChange={(e) => setFormData({ ...formData, newPassword: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label>Confirmar Contraseña</label>
            <input
              type="password"
              placeholder="Repite la nueva contraseña"
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
              required
            />
          </div>

          <button type="submit" className="btn-modal" disabled={loading}>
            {loading ? 'Guardando...' : 'Guardar'}
          </button>
        </form>
      </div>

      {toast.show && <Toast message={toast.message} type={toast.type} />}
    </div>
  );
}

export default ChangePassword;
