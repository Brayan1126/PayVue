import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import Toast from '../components/Toast';
import { api, getCurrentUserId } from '../config/api';

function AddPayment() {
  const navigate = useNavigate();
  const [debts, setDebts] = useState([]);
  const [formData, setFormData] = useState({ amount: '', debt_id: '', receipt: null });
  const [toast, setToast] = useState({ show: false, message: '', type: '' });
  const [loading, setLoading] = useState(false);

  const fetchDebts = useCallback(async () => {
    try {
      const res = await api.get('/finances/debt');
      setDebts(res.data || []);
    } catch (error) {
      console.error('Error fetching debts:', error);
    }
  }, []);

  useEffect(() => { fetchDebts(); }, [fetchDebts]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const data = new FormData();
      data.append('amount', formData.amount);
      data.append('debt_id', formData.debt_id);
      data.append('user_id', getCurrentUserId());
      if (formData.receipt) data.append('receipt', formData.receipt);

      await api.post('/finances/payment', data, { headers: { 'Content-Type': 'multipart/form-data' } });
      setToast({ show: true, message: '¡Pago guardado con éxito!', type: 'success' });
      setFormData({ amount: '', debt_id: '', receipt: null });
      fetchDebts();
      setTimeout(() => { setToast({ show: false, message: '', type: '' }); navigate('/dashboard'); }, 2000);
    } catch (error) {
      setToast({ show: true, message: 'Error al guardar el pago', type: 'error' });
      setTimeout(() => setToast({ show: false, message: '', type: '' }), 3000);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-page">
      <div className="modal-card" style={{ position: 'relative' }}>
        <button onClick={() => navigate('/dashboard')} style={{ position: 'absolute', top: '20px', right: '25px', background: 'none', border: 'none', fontSize: '1.8rem', cursor: 'pointer', color: '#666', lineHeight: 1 }}>×</button>
        <h2>Agregar Pago</h2>
        <form className="modal-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Monto de pago</label>
            <input type="number" placeholder="Ingresa monto de pago" value={formData.amount} onChange={(e) => setFormData({ ...formData, amount: e.target.value })} required />
          </div>
          <div className="form-group">
            <label>Recibo</label>
            <label className="file-upload">
              <span>{formData.receipt ? formData.receipt.name : 'Sube tu Recibo'}</span>
              <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              <input type="file" accept="image/*,.pdf" style={{ display: 'none' }} onChange={(e) => setFormData({ ...formData, receipt: e.target.files[0] })} />
            </label>
          </div>
          <div className="form-group">
            <label>Selección de la deuda</label>
            <select value={formData.debt_id} onChange={(e) => setFormData({ ...formData, debt_id: e.target.value })} required>
              <option value="">Selecciona deuda</option>
              {debts.map((debt) => (<option key={debt.id} value={debt.id}>{debt.name} - ${debt.remaining_amount?.toLocaleString('es-CO')}</option>))}
            </select>
          </div>
          <button type="submit" className="btn-modal" disabled={loading}>{loading ? 'Guardando...' : 'Guardar'}</button>
        </form>
      </div>
      {toast.show && <Toast message={toast.message} type={toast.type} />}
    </div>
  );
}

export default AddPayment;
