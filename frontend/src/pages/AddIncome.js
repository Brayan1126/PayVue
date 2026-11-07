import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Toast from '../components/Toast';
import { api } from '../config/api';

function AddIncome() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({ amount: '', source: '', date: new Date().toISOString().split('T')[0] });
  const [toast, setToast] = useState({ show: false, message: '', type: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post('/finances/income', {
        amount: parseFloat(formData.amount),
        source: formData.source,
        date: formData.date
      });
      setToast({ show: true, message: '¡Ingreso guardado con éxito!', type: 'success' });
      setFormData({ amount: '', source: '', date: new Date().toISOString().split('T')[0] });
      setTimeout(() => { setToast({ show: false, message: '', type: '' }); navigate('/dashboard'); }, 2000);
    } catch (error) {
      setToast({ show: true, message: 'Error al guardar el ingreso', type: 'error' });
      setTimeout(() => setToast({ show: false, message: '', type: '' }), 3000);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-page">
      <div className="modal-card" style={{ position: 'relative' }}>
        <button onClick={() => navigate('/dashboard')} style={{ position: 'absolute', top: '20px', right: '25px', background: 'none', border: 'none', fontSize: '1.8rem', cursor: 'pointer', color: '#666', lineHeight: 1 }}>×</button>
        <h2>Agregar Ingreso</h2>
        <form className="modal-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Monto</label>
            <input type="number" placeholder="Ingresa el monto" value={formData.amount} onChange={(e) => setFormData({ ...formData, amount: e.target.value })} required />
          </div>
          <div className="form-group">
            <label>Fuente</label>
            <input type="text" placeholder="Ej: Salario, Freelance" value={formData.source} onChange={(e) => setFormData({ ...formData, source: e.target.value })} required />
          </div>
          <div className="form-group">
            <label>Fecha</label>
            <input type="date" value={formData.date} onChange={(e) => setFormData({ ...formData, date: e.target.value })} required />
          </div>
          <button type="submit" className="btn-modal" disabled={loading}>{loading ? 'Guardando...' : 'Guardar'}</button>
        </form>
      </div>
      {toast.show && <Toast message={toast.message} type={toast.type} />}
    </div>
  );
}

export default AddIncome;
