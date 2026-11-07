import React, { useEffect, useState, useCallback } from 'react';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend } from 'chart.js';
import { Line, Bar } from 'react-chartjs-2';
import Sidebar from '../components/Sidebar';
import ProfilePanel from '../components/ProfilePanel';
import { api } from '../config/api';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend);

function Dashboard() {
  const [incomes, setIncomes] = useState([]);
  const [debts, setDebts] = useState([]);
  const [payments, setPayments] = useState([]);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [profileOpen, setProfileOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

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

  const totalIncomes = incomes.reduce((sum, i) => sum + (i.amount || 0), 0);
  const totalExpenses = debts.reduce((sum, d) => sum + (d.installment_amount || 0), 0);
  const totalPaid = payments.reduce((sum, p) => sum + (p.amount || 0), 0);
  const savings = totalIncomes - totalExpenses;

  const filteredDebts = debts.filter(debt => 
    debt.name?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const sortedIncomes = [...incomes].sort((a, b) => new Date(a.date) - new Date(b.date));
  const lineChartData = {
    labels: sortedIncomes.length > 0 
      ? sortedIncomes.map(i => {
          const date = new Date(i.date);
          return `${date.getDate()}/${date.getMonth() + 1}`;
        })
      : [],
    datasets: [{
      label: 'Ingresos',
      data: sortedIncomes.map(i => i.amount),
      borderColor: '#1a1a1a',
      backgroundColor: '#1a1a1a',
      tension: 0.4,
      fill: false,
      pointRadius: 4,
    }]
  };

  const expensesByMonth = new Array(12).fill(0);
  debts.forEach(debt => {
    const month = new Date().getMonth();
    expensesByMonth[month] += debt.installment_amount || 0;
  });
  
  const barChartData = {
    labels: ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'],
    datasets: [{
      label: 'Gastos',
      data: expensesByMonth,
      backgroundColor: '#1a1a1a',
    }]
  };

  const chartOptions = {
    responsive: true,
    plugins: { legend: { display: false } },
    scales: { y: { beginAtZero: true } }
  };

  const calculateDaysUntilPayment = (paymentDay) => {
    const today = new Date();
    const currentMonth = today.getMonth();
    const currentYear = today.getFullYear();
    const paymentDate = new Date(currentYear, currentMonth, paymentDay);
    if (paymentDate < today) paymentDate.setMonth(currentMonth + 1);
    return Math.ceil((paymentDate - today) / (1000 * 60 * 60 * 24));
  };

  const hasData = incomes.length > 0 || debts.length > 0;

  return (
    <div className="dashboard-layout">
      <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      
      <main className="main-content">
        <header className="header">
          <div className="header-left">
            <button className="menu-toggle" onClick={() => setSidebarOpen(true)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <line x1="3" y1="12" x2="21" y2="12"/>
                <line x1="3" y1="6" x2="21" y2="6"/>
                <line x1="3" y1="18" x2="21" y2="18"/>
              </svg>
            </button>
            <div className="search-box">
              <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="#888" strokeWidth="2">
                <circle cx="11" cy="11" r="8"/>
                <path d="m21 21-4.35-4.35"/>
              </svg>
              <input 
                type="text" 
                placeholder="Buscar deudas..." 
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>
          <div className="header-right">
            <img 
              src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150" 
              alt="User" 
              className="user-avatar"
              onClick={() => setProfileOpen(true)}
            />
          </div>
        </header>

        <div className="stats-grid">
          <div className="stat-card">
            <div className="label">Ingresos mensuales</div>
            <div className="value">${totalIncomes.toLocaleString('es-CO', { minimumFractionDigits: 0 })}</div>
          </div>
          <div className="stat-card">
            <div className="label">Ahorro Disponible</div>
            <div className="value" style={{ color: savings >= 0 ? '#22c55e' : '#ef4444' }}>
              ${savings.toLocaleString('es-CO', { minimumFractionDigits: 0 })}
            </div>
          </div>
          <div className="stat-card">
            <div className="label">Gastos Mensuales</div>
            <div className="value">${totalExpenses.toLocaleString('es-CO', { minimumFractionDigits: 0 })}</div>
          </div>
        </div>

        {!hasData ? (
          <div className="chart-card" style={{ textAlign: 'center', padding: '60px 20px' }}>
            <svg viewBox="0 0 24 24" width="64" height="64" fill="none" stroke="#ccc" strokeWidth="1.5" style={{ marginBottom: '20px' }}>
              <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              <path d="M9 10h.01M15 10h.01M9.5 15.5s1.5 1 2.5 1 2.5-1 2.5-1"/>
            </svg>
            <h3 style={{ color: '#666', marginBottom: '10px' }}>No hay datos todavía</h3>
            <p style={{ color: '#888' }}>Agrega ingresos y deudas para ver tus estadísticas</p>
          </div>
        ) : (
          <>
            <div className="charts-grid">
              <div className="chart-card">
                <h3>Historial de Ingresos</h3>
                {incomes.length > 0 ? (
                  <Line data={lineChartData} options={chartOptions} />
                ) : (
                  <p style={{ color: '#888', textAlign: 'center', padding: '40px' }}>No hay ingresos registrados</p>
                )}
              </div>
              <div className="chart-card">
                <h3>Deudas Restantes</h3>
                <div className="debts-list">
                  {filteredDebts.length > 0 ? filteredDebts.map((debt) => (
                    <div key={debt.id} className="debt-item">
                      <span className="name">{debt.name}</span>
                      <span className="installments">
                        {debt.remaining_amount ? Math.ceil(debt.remaining_amount / debt.installment_amount) : '-'} cuotas
                      </span>
                    </div>
                  )) : (
                    <p style={{ color: '#888', textAlign: 'center', padding: '20px' }}>
                      {searchTerm ? 'No se encontraron deudas' : 'No hay deudas registradas'}
                    </p>
                  )}
                </div>
              </div>
            </div>

            <div className="bottom-grid">
              <div className="chart-card">
                <h3>Deudas Próximas a Vencer</h3>
                {debts.length > 0 ? (
                  <table className="data-table">
                    <thead>
                      <tr>
                        <th>Entidad</th>
                        <th>Valor</th>
                        <th>Días Restantes</th>
                      </tr>
                    </thead>
                    <tbody>
                      {filteredDebts.slice(0, 6).map((debt) => (
                        <tr key={debt.id}>
                          <td>{debt.name}</td>
                          <td>${debt.installment_amount?.toLocaleString('es-CO')}</td>
                          <td>{calculateDaysUntilPayment(debt.payment_day)} Días</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                ) : (
                  <p style={{ color: '#888', textAlign: 'center', padding: '40px' }}>No hay deudas registradas</p>
                )}
              </div>
              <div className="chart-card">
                <h3>Gastos por Mes</h3>
                {debts.length > 0 ? (
                  <Bar data={barChartData} options={chartOptions} />
                ) : (
                  <p style={{ color: '#888', textAlign: 'center', padding: '40px' }}>No hay gastos registrados</p>
                )}
              </div>
            </div>
          </>
        )}
      </main>

      <ProfilePanel isOpen={profileOpen} onClose={() => setProfileOpen(false)} />
    </div>
  );
}

export default Dashboard;
