import React, { useState, useEffect, useCallback } from 'react';
import Sidebar from '../components/Sidebar';
import Toast from '../components/Toast';
import { api } from '../config/api';
import API_URL from '../config/api';

function History() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [incomes, setIncomes] = useState([]);
  const [debts, setDebts] = useState([]);
  const [payments, setPayments] = useState([]);
  const [activeTab, setActiveTab] = useState('debts');
  const [searchTerm, setSearchTerm] = useState('');
  const [dateFilter, setDateFilter] = useState({ start: '', end: '' });
  const [debtFilter, setDebtFilter] = useState('');
  const [toast, setToast] = useState({ show: false, message: '', type: '' });
  const [editModal, setEditModal] = useState({ show: false, type: '', data: null });
  const [receiptModal, setReceiptModal] = useState({ show: false, url: '' });
  const [deleteConfirm, setDeleteConfirm] = useState({ show: false, type: '', id: null });

  const fetchData = useCallback(async () => {
    try {
      const [incomeRes, debtRes, paymentRes] = await Promise.all([
        api.get('/finances/income'),
        api.get('/finances/debt'),
        api.get('/finances/payment')
      ]);
      setIncomes(incomeRes.data || []);
      setDebts(debtRes.data || []);
      setPayments(paymentRes.data || []);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, [fetchData]);

  const showToast = (message, type) => {
    setToast({ show: true, message, type });
    setTimeout(() => setToast({ show: false, message: '', type: '' }), 3000);
  };

  const filterByDate = (items, dateField) => {
    return items.filter(item => {
      const itemDate = new Date(item[dateField]);
      const startOk = !dateFilter.start || itemDate >= new Date(dateFilter.start);
      const endOk = !dateFilter.end || itemDate <= new Date(dateFilter.end);
      return startOk && endOk;
    });
  };

  const filteredIncomes = filterByDate(incomes, 'date').filter(income =>
    income.source?.toLowerCase().includes(searchTerm.toLowerCase()) || income.amount?.toString().includes(searchTerm)
  );

  const filteredDebts = debts.filter(debt => debt.name?.toLowerCase().includes(searchTerm.toLowerCase()));

  const filteredPayments = payments.filter(payment => {
    // Filtrar por búsqueda
    const matchSearch = !searchTerm || 
      payment.amount?.toString().includes(searchTerm) ||
      payment.debt_name?.toLowerCase().includes(searchTerm.toLowerCase());
    
    // Filtrar por deuda
    const matchDebt = !debtFilter || payment.debt_id?.toString() === debtFilter;
    
    // Filtrar por fecha
    const paymentDate = new Date(payment.date || payment.created_at);
    const startOk = !dateFilter.start || paymentDate >= new Date(dateFilter.start);
    const endOk = !dateFilter.end || paymentDate <= new Date(dateFilter.end);
    
    return matchSearch && matchDebt && startOk && endOk;
  });

  const handleEdit = (type, item) => setEditModal({ show: true, type, data: { ...item } });

  const handleSaveEdit = async () => {
    const { type, data } = editModal;
    try {
      if (type === 'income') {
        await api.put(`/finances/income/${data.id}`, { amount: parseFloat(data.amount), source: data.source, date: data.date });
      } else if (type === 'debt') {
        await api.put(`/finances/debt/${data.id}`, {
          name: data.name, total_amount: parseFloat(data.total_amount), remaining_amount: parseFloat(data.remaining_amount),
          installment_amount: parseFloat(data.installment_amount), payment_day: parseInt(data.payment_day),
          due_date: data.due_date?.split('T')[0] || data.due_date, num_installments: parseInt(data.num_installments),
          interest_rate: parseFloat(data.interest_rate || 0), paid: data.paid || false
        });
      }
      showToast('¡Actualizado con éxito!', 'success');
      setEditModal({ show: false, type: '', data: null });
      fetchData();
    } catch (error) {
      console.error('Error updating:', error.response?.data);
      showToast('Error al actualizar', 'error');
    }
  };

  const handleDelete = async () => {
    const { type, id } = deleteConfirm;
    try {
      if (type === 'income') await api.delete(`/finances/income/${id}`);
      else if (type === 'debt') await api.delete(`/finances/debt/${id}`);
      else if (type === 'payment') await api.delete(`/finances/payment/${id}`);
      showToast('¡Eliminado con éxito!', 'success');
      setDeleteConfirm({ show: false, type: '', id: null });
      fetchData();
    } catch (error) {
      showToast('Error al eliminar', 'error');
    }
  };

  const getDebtName = (debtId) => debts.find(d => d.id === debtId)?.name || 'Desconocido';
  const clearFilters = () => { setSearchTerm(''); setDateFilter({ start: '', end: '' }); setDebtFilter(''); };

  const actionBtnStyle = { background: 'none', border: '1px solid #ddd', borderRadius: '6px', padding: '6px', cursor: 'pointer', display: 'flex', alignItems: 'center', justifyContent: 'center' };

  return (
    <div className="dashboard-layout">
      <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      <main className="main-content">
        <header className="header">
          <div className="header-left">
            <button className="menu-toggle" onClick={() => setSidebarOpen(true)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <h1 style={{ fontSize: '1.5rem', fontWeight: '600' }}>Historial</h1>
          </div>
        </header>

        <div style={{ display: 'flex', gap: '10px', marginBottom: '20px', flexWrap: 'wrap' }}>
          {['debts', 'incomes', 'payments'].map(tab => (
            <button key={tab} onClick={() => { setActiveTab(tab); clearFilters(); }} style={{ padding: '12px 25px', background: activeTab === tab ? '#1a1a1a' : '#fff', color: activeTab === tab ? '#fff' : '#1a1a1a', border: '1px solid #1a1a1a', borderRadius: '8px', cursor: 'pointer', fontFamily: 'Poppins, sans-serif', fontWeight: '500' }}>
              {tab === 'debts' ? 'Deudas' : tab === 'incomes' ? 'Ingresos' : 'Pagos'}
            </button>
          ))}
        </div>

        <div className="chart-card" style={{ marginBottom: '20px', padding: '20px' }}>
          <div style={{ display: 'flex', gap: '15px', flexWrap: 'wrap', alignItems: 'flex-end' }}>
            <div className="form-group" style={{ flex: '1', minWidth: '200px', marginBottom: 0 }}>
              <label style={{ fontSize: '0.85rem', color: '#666' }}>Buscar</label>
              <input type="text" placeholder="Buscar..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ddd' }} />
            </div>
            {(activeTab === 'incomes' || activeTab === 'payments') && (
              <>
                <div className="form-group" style={{ minWidth: '150px', marginBottom: 0 }}>
                  <label style={{ fontSize: '0.85rem', color: '#666' }}>Desde</label>
                  <input type="date" value={dateFilter.start} onChange={(e) => setDateFilter({ ...dateFilter, start: e.target.value })} style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ddd' }} />
                </div>
                <div className="form-group" style={{ minWidth: '150px', marginBottom: 0 }}>
                  <label style={{ fontSize: '0.85rem', color: '#666' }}>Hasta</label>
                  <input type="date" value={dateFilter.end} onChange={(e) => setDateFilter({ ...dateFilter, end: e.target.value })} style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ddd' }} />
                </div>
              </>
            )}
            {activeTab === 'payments' && (
              <div className="form-group" style={{ minWidth: '180px', marginBottom: 0 }}>
                <label style={{ fontSize: '0.85rem', color: '#666' }}>Deuda</label>
                <select value={debtFilter} onChange={(e) => setDebtFilter(e.target.value)} style={{ padding: '12px', borderRadius: '8px', border: '1px solid #ddd', width: '100%' }}>
                  <option value="">Todas las deudas</option>
                  {debts.map(debt => (<option key={debt.id} value={debt.id}>{debt.name}</option>))}
                </select>
              </div>
            )}
            <button onClick={clearFilters} style={{ padding: '12px 20px', background: '#f5f5f5', border: '1px solid #ddd', borderRadius: '8px', cursor: 'pointer', fontFamily: 'Poppins' }}>Limpiar</button>
          </div>
        </div>

        {activeTab === 'debts' && (
          <div className="chart-card">
            <h3 style={{ marginBottom: '20px' }}>Lista de Deudas ({filteredDebts.length})</h3>
            {filteredDebts.length > 0 ? (
              <div style={{ overflowX: 'auto' }}>
                <table className="data-table">
                  <thead><tr><th>Nombre</th><th>Total</th><th>Restante</th><th>Cuota</th><th>Día Pago</th><th>Estado</th><th>Acciones</th></tr></thead>
                  <tbody>
                    {filteredDebts.map(debt => (
                      <tr key={debt.id}>
                        <td style={{ fontWeight: '500' }}>{debt.name}</td>
                        <td>${debt.total_amount?.toLocaleString('es-CO')}</td>
                        <td>${debt.remaining_amount?.toLocaleString('es-CO')}</td>
                        <td>${debt.installment_amount?.toLocaleString('es-CO')}</td>
                        <td>Día {debt.payment_day}</td>
                        <td><span style={{ padding: '4px 12px', borderRadius: '20px', fontSize: '0.8rem', background: debt.paid ? '#dcfce7' : '#fef3c7', color: debt.paid ? '#166534' : '#92400e' }}>{debt.paid ? 'Pagada' : 'Pendiente'}</span></td>
                        <td>
                          <div style={{ display: 'flex', gap: '8px' }}>
                            <button onClick={() => handleEdit('debt', debt)} style={actionBtnStyle}><EditIcon /></button>
                            <button onClick={() => setDeleteConfirm({ show: true, type: 'debt', id: debt.id })} style={{ ...actionBtnStyle, color: '#ef4444' }}><DeleteIcon /></button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (<p style={{ color: '#888', textAlign: 'center', padding: '40px' }}>No hay deudas registradas</p>)}
          </div>
        )}

        {activeTab === 'incomes' && (
          <div className="chart-card">
            <h3 style={{ marginBottom: '20px' }}>Historial de Ingresos ({filteredIncomes.length})</h3>
            {filteredIncomes.length > 0 ? (
              <table className="data-table">
                <thead><tr><th>Fuente</th><th>Monto</th><th>Fecha</th><th>Acciones</th></tr></thead>
                <tbody>
                  {filteredIncomes.map(income => (
                    <tr key={income.id}>
                      <td style={{ fontWeight: '500' }}>{income.source}</td>
                      <td style={{ color: '#22c55e', fontWeight: '500' }}>+${income.amount?.toLocaleString('es-CO')}</td>
                      <td>{income.date}</td>
                      <td>
                        <div style={{ display: 'flex', gap: '8px' }}>
                          <button onClick={() => handleEdit('income', income)} style={actionBtnStyle}><EditIcon /></button>
                          <button onClick={() => setDeleteConfirm({ show: true, type: 'income', id: income.id })} style={{ ...actionBtnStyle, color: '#ef4444' }}><DeleteIcon /></button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (<p style={{ color: '#888', textAlign: 'center', padding: '40px' }}>No hay ingresos registrados</p>)}
          </div>
        )}

        {activeTab === 'payments' && (
          <div className="chart-card">
            <h3 style={{ marginBottom: '20px' }}>Historial de Pagos ({filteredPayments.length})</h3>
            {filteredPayments.length > 0 ? (
              <table className="data-table">
                <thead><tr><th>Deuda</th><th>Monto</th><th>Fecha</th><th>Recibo</th><th>Acciones</th></tr></thead>
                <tbody>
                  {filteredPayments.map(payment => (
                    <tr key={payment.id}>
                      <td style={{ fontWeight: '500' }}>{payment.debt_name || getDebtName(payment.debt_id)}</td>
                      <td style={{ color: '#ef4444', fontWeight: '500' }}>-${payment.amount?.toLocaleString('es-CO')}</td>
                      <td>{payment.date || payment.created_at}</td>
                      <td>
                        {payment.receipt_url ? (
                          <button onClick={() => setReceiptModal({ show: true, url: `${API_URL}${payment.receipt_url}` })} style={{ padding: '6px 12px', background: '#1a1a1a', color: '#fff', border: 'none', borderRadius: '6px', cursor: 'pointer', fontSize: '0.8rem' }}>Ver Recibo</button>
                        ) : (<span style={{ color: '#888', fontSize: '0.85rem' }}>Sin recibo</span>)}
                      </td>
                      <td><button onClick={() => setDeleteConfirm({ show: true, type: 'payment', id: payment.id })} style={{ ...actionBtnStyle, color: '#ef4444' }}><DeleteIcon /></button></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (<p style={{ color: '#888', textAlign: 'center', padding: '40px' }}>No hay pagos registrados</p>)}
          </div>
        )}
      </main>

      {editModal.show && (
        <div className="modal-overlay" onClick={() => setEditModal({ show: false, type: '', data: null })}>
          <div className="modal" onClick={e => e.stopPropagation()} style={{ maxWidth: '500px' }}>
            <h2 style={{ marginBottom: '25px' }}>Editar {editModal.type === 'income' ? 'Ingreso' : 'Deuda'}</h2>
            {editModal.type === 'income' && (
              <div className="modal-form">
                <div className="form-group"><label>Fuente</label><input type="text" value={editModal.data.source || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, source: e.target.value } })} /></div>
                <div className="form-group"><label>Monto</label><input type="number" value={editModal.data.amount || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, amount: e.target.value } })} /></div>
                <div className="form-group"><label>Fecha</label><input type="date" value={editModal.data.date || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, date: e.target.value } })} /></div>
              </div>
            )}
            {editModal.type === 'debt' && (
              <div className="modal-form">
                <div className="form-group"><label>Nombre</label><input type="text" value={editModal.data.name || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, name: e.target.value } })} /></div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px' }}>
                  <div className="form-group"><label>Monto Total</label><input type="number" value={editModal.data.total_amount || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, total_amount: e.target.value } })} /></div>
                  <div className="form-group"><label>Monto Restante</label><input type="number" value={editModal.data.remaining_amount || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, remaining_amount: e.target.value } })} /></div>
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px' }}>
                  <div className="form-group"><label>Cuota Mensual</label><input type="number" value={editModal.data.installment_amount || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, installment_amount: e.target.value } })} /></div>
                  <div className="form-group"><label>Día de Pago</label><input type="number" min="1" max="31" value={editModal.data.payment_day || ''} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, payment_day: e.target.value } })} /></div>
                </div>
                <div className="form-group"><label style={{ display: 'flex', alignItems: 'center', gap: '10px', cursor: 'pointer' }}><input type="checkbox" checked={editModal.data.paid || false} onChange={(e) => setEditModal({ ...editModal, data: { ...editModal.data, paid: e.target.checked } })} style={{ width: '20px', height: '20px' }} />Marcar como pagada</label></div>
              </div>
            )}
            <div style={{ display: 'flex', gap: '10px', marginTop: '25px' }}>
              <button onClick={() => setEditModal({ show: false, type: '', data: null })} style={{ flex: 1, padding: '12px', background: '#f5f5f5', border: '1px solid #ddd', borderRadius: '8px', cursor: 'pointer' }}>Cancelar</button>
              <button onClick={handleSaveEdit} className="btn-modal" style={{ flex: 1, margin: 0 }}>Guardar</button>
            </div>
          </div>
        </div>
      )}

      {deleteConfirm.show && (
        <div className="modal-overlay" onClick={() => setDeleteConfirm({ show: false, type: '', id: null })}>
          <div className="modal" onClick={e => e.stopPropagation()} style={{ maxWidth: '400px', textAlign: 'center' }}>
            <div style={{ fontSize: '3rem', marginBottom: '15px' }}>⚠️</div>
            <h2 style={{ marginBottom: '10px' }}>¿Estás seguro?</h2>
            <p style={{ color: '#666', marginBottom: '25px' }}>Esta acción no se puede deshacer</p>
            <div style={{ display: 'flex', gap: '10px' }}>
              <button onClick={() => setDeleteConfirm({ show: false, type: '', id: null })} style={{ flex: 1, padding: '12px', background: '#f5f5f5', border: '1px solid #ddd', borderRadius: '8px', cursor: 'pointer' }}>Cancelar</button>
              <button onClick={handleDelete} style={{ flex: 1, padding: '12px', background: '#ef4444', color: '#fff', border: 'none', borderRadius: '8px', cursor: 'pointer' }}>Eliminar</button>
            </div>
          </div>
        </div>
      )}

      {receiptModal.show && (
        <div className="modal-overlay" onClick={() => setReceiptModal({ show: false, url: '' })}>
          <div className="modal" onClick={e => e.stopPropagation()} style={{ maxWidth: '600px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
              <h2 style={{ margin: 0 }}>Recibo</h2>
              <button onClick={() => setReceiptModal({ show: false, url: '' })} style={{ background: 'none', border: 'none', fontSize: '1.5rem', cursor: 'pointer' }}>×</button>
            </div>
            <div style={{ textAlign: 'center' }}>
              <img src={receiptModal.url} alt="Recibo" style={{ maxWidth: '100%', maxHeight: '70vh', borderRadius: '8px' }} onError={(e) => { e.target.style.display = 'none'; e.target.nextSibling.style.display = 'block'; }} />
              <div style={{ display: 'none', padding: '40px', background: '#f5f5f5', borderRadius: '8px' }}><p>No se pudo cargar la imagen</p><a href={receiptModal.url} target="_blank" rel="noopener noreferrer" style={{ color: '#1a1a1a' }}>Abrir en nueva pestaña</a></div>
            </div>
          </div>
        </div>
      )}

      {toast.show && <Toast message={toast.message} type={toast.type} />}
    </div>
  );
}

const EditIcon = () => (<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>);
const DeleteIcon = () => (<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>);

export default History;
