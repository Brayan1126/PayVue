import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Toast from '../components/Toast';
import { api } from '../config/api';

function AddDebt() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    total_amount: '',
    remaining_amount: '',
    installment_amount: '',
    payment_day: '',
    due_date: '',
    num_installments: '',
    interest_rate: '0'
  });
  const [toast, setToast] = useState({ show: false, message: '', type: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post('/finances/debt', {
        name: formData.name,
        total_amount: parseFloat(formData.total_amount),
        remaining_amount: parseFloat(formData.remaining_amount || formData.total_amount),
        installment_amount: parseFloat(formData.installment_amount),
        payment_day: parseInt(formData.payment_day),
        due_date: formData.due_date,
        num_installments: parseInt(formData.num_installments),
        interest_rate: parseFloat(formData.interest_rate || 0)
      });
      
      setToast({ show: true, message: '¡Deuda guardada con éxito!', type: 'success' });
      setFormData({ 
        name: '', total_amount: '', remaining_amount: '', 
        installment_amount: '', payment_day: '', due_date: '',
        num_installments: '', interest_rate: '0'
      });
      
      setTimeout(() => {
        setToast({ show: false, message: '', type: '' });
        navigate('/dashboard');
      }, 2000);
    } catch (error) {
      console.error('Error:', error.response?.data);
      setToast({ show: true, message: error.response?.data?.message || 'Error al guardar la deuda', type: 'error' });
      setTimeout(() => setToast({ show: false, message: '', type: '' }), 3000);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => navigate('/dashboard');

  const handleTotalAmountChange = (value) => {
    setFormData(prev => {
      const newData = { ...prev, total_amount: value };
      if (value && prev.installment_amount) {
        newData.num_installments = Math.ceil(parseFloat(value) / parseFloat(prev.installment_amount)).toString();
      }
      if (!prev.remaining_amount) {
        newData.remaining_amount = value;
      }
      return newData;
    });
  };

  const handleInstallmentChange = (value) => {
    setFormData(prev => {
      const newData = { ...prev, installment_amount: value };
      if (value && prev.total_amount) {
        newData.num_installments = Math.ceil(parseFloat(prev.total_amount) / parseFloat(value)).toString();
      }
      return newData;
    });
  };

  return (
    <div className="modal-page">
      <div className="modal-card" style={{ position: 'relative', maxWidth: '520px' }}>
        <button onClick={handleClose} style={{ position: 'absolute', top: '20px', right: '25px', background: 'none', border: 'none', fontSize: '1.8rem', cursor: 'pointer', color: '#666', lineHeight: 1 }}>×</button>
        <h2>Agregar Deuda</h2>
        <form className="modal-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Nombre / Entidad</label>
            <input type="text" placeholder="Ej: Tarjeta Visa" value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} required />
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px' }}>
            <div className="form-group">
              <label>Monto Total</label>
              <input type="number" placeholder="Valor total" value={formData.total_amount} onChange={(e) => handleTotalAmountChange(e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Monto Restante</label>
              <input type="number" placeholder="Lo que falta" value={formData.remaining_amount} onChange={(e) => setFormData({ ...formData, remaining_amount: e.target.value })} />
            </div>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px' }}>
            <div className="form-group">
              <label>Cuota Mensual</label>
              <input type="number" placeholder="Valor cuota" value={formData.installment_amount} onChange={(e) => handleInstallmentChange(e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Número de Cuotas</label>
              <input type="number" placeholder="Total cuotas" value={formData.num_installments} onChange={(e) => setFormData({ ...formData, num_installments: e.target.value })} required />
            </div>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px' }}>
            <div className="form-group">
              <label>Día de Pago (1-31)</label>
              <input type="number" min="1" max="31" placeholder="Día del mes" value={formData.payment_day} onChange={(e) => setFormData({ ...formData, payment_day: e.target.value })} required />
            </div>
            <div className="form-group">
              <label>Fecha Límite</label>
              <input type="date" value={formData.due_date} onChange={(e) => setFormData({ ...formData, due_date: e.target.value })} required />
            </div>
          </div>
          <div className="form-group">
            <label>Tasa de Interés (%) - Opcional</label>
            <input type="number" step="0.01" placeholder="0" value={formData.interest_rate} onChange={(e) => setFormData({ ...formData, interest_rate: e.target.value })} />
          </div>
          <button type="submit" className="btn-modal" disabled={loading}>{loading ? 'Guardando...' : 'Guardar'}</button>
        </form>
      </div>
      {toast.show && <Toast message={toast.message} type={toast.type} />}
    </div>
  );
}

export default AddDebt;
