// Mock en memoria de la API para correr el frontend sin backend.
// Se activa con REACT_APP_MOCK_API=true (ver frontend/.env.local).

const iso = (daysAgo) => {
  const d = new Date();
  d.setDate(d.getDate() - daysAgo);
  return d.toISOString().slice(0, 10);
};

const db = {
  nextId: 100,
  incomes: [
    { id: 1, amount: 3500000, source: 'Salario', date: iso(60) },
    { id: 2, amount: 450000, source: 'Freelance', date: iso(45) },
    { id: 3, amount: 3500000, source: 'Salario', date: iso(30) },
    { id: 4, amount: 280000, source: 'Venta', date: iso(15) },
    { id: 5, amount: 3500000, source: 'Salario', date: iso(1) },
  ],
  debts: [
    { id: 1, name: 'Tarjeta de crédito', total_amount: 2400000, remaining_amount: 1600000, due_date: iso(-180), interest_rate: 2.1, num_installments: 12, installment_amount: 200000, payment_day: 15, remaining_payments: 8, paid: false },
    { id: 2, name: 'Préstamo estudiantil', total_amount: 6000000, remaining_amount: 4500000, due_date: iso(-540), interest_rate: 1.2, num_installments: 24, installment_amount: 250000, payment_day: 5, remaining_payments: 18, paid: false },
    { id: 3, name: 'Celular', total_amount: 1200000, remaining_amount: 300000, due_date: iso(-60), interest_rate: 0, num_installments: 6, installment_amount: 200000, payment_day: 20, remaining_payments: 2, paid: false },
  ],
  payments: [
    { id: 1, debt_id: 1, amount: 200000, date: iso(40), created_at: iso(40), debt_name: 'Tarjeta de crédito', remaining_installments: 9, remaining_amount: 1800000, receipt_url: '' },
    { id: 2, debt_id: 2, amount: 250000, date: iso(35), created_at: iso(35), debt_name: 'Préstamo estudiantil', remaining_installments: 18, remaining_amount: 4500000, receipt_url: '' },
    { id: 3, debt_id: 1, amount: 200000, date: iso(10), created_at: iso(10), debt_name: 'Tarjeta de crédito', remaining_installments: 8, remaining_amount: 1600000, receipt_url: '' },
    { id: 4, debt_id: 3, amount: 200000, date: iso(5), created_at: iso(5), debt_name: 'Celular', remaining_installments: 2, remaining_amount: 300000, receipt_url: '' },
  ],
};

const parseBody = (data) => {
  if (!data) return {};
  if (typeof FormData !== 'undefined' && data instanceof FormData) {
    return Object.fromEntries(data.entries());
  }
  if (typeof data === 'string') {
    try { return JSON.parse(data); } catch { return {}; }
  }
  return data;
};

const respond = (config, status, data) =>
  new Promise((resolve, reject) => {
    const response = { data, status, statusText: String(status), headers: {}, config, request: {} };
    setTimeout(() => {
      if (status >= 400) {
        const err = new Error(data?.message || 'Mock error');
        err.response = response;
        reject(err);
      } else {
        resolve(response);
      }
    }, 150);
  });

const collections = { income: 'incomes', debt: 'debts', payment: 'payments' };

export function mockAdapter(config) {
  const method = (config.method || 'get').toLowerCase();
  const path = (config.url || '').replace(/^https?:\/\/[^/]+/, '').replace(/\/+$/, '');
  const body = parseBody(config.data);

  if (path === '/auth/login' || path === '/auth/register') {
    return respond(config, 200, { user_id: 1, email: body.email || 'demo@payvue.com', message: 'OK' });
  }
  if (path === '/auth/logout') {
    return respond(config, 200, { message: 'Sesión cerrada exitosamente' });
  }

  const m = path.match(/^\/finances\/(income|debt|payment)(?:\/(\d+))?$/);
  if (!m) return respond(config, 404, { message: `Mock: ruta no encontrada ${path}` });

  const key = collections[m[1]];
  const id = m[2] ? Number(m[2]) : null;
  const list = db[key];

  if (method === 'get') {
    if (id == null) return respond(config, 200, list);
    const item = list.find((x) => x.id === id);
    return item ? respond(config, 200, item) : respond(config, 404, { message: 'No encontrado' });
  }

  if (method === 'post') {
    const item = { ...body, id: db.nextId++ };
    if (item.amount != null) item.amount = Number(item.amount);
    if (key === 'debts') {
      item.remaining_payments = Number(item.num_installments) || 0;
      item.paid = false;
    }
    if (key === 'payments') {
      const debt = db.debts.find((d) => d.id === Number(item.debt_id));
      item.debt_id = Number(item.debt_id);
      item.date = item.date || iso(0);
      item.created_at = iso(0);
      item.debt_name = debt?.name || '';
      item.receipt_url = '';
      if (debt) {
        debt.remaining_amount = Math.max(0, debt.remaining_amount - item.amount);
        debt.remaining_payments = Math.max(0, debt.remaining_payments - 1);
        item.remaining_amount = debt.remaining_amount;
        item.remaining_installments = debt.remaining_payments;
      }
      delete item.receipt;
    }
    list.push(item);
    return respond(config, 201, item);
  }

  if (method === 'put' && id != null) {
    const idx = list.findIndex((x) => x.id === id);
    if (idx === -1) return respond(config, 404, { message: 'No encontrado' });
    list[idx] = { ...list[idx], ...body, id };
    return respond(config, 200, list[idx]);
  }

  if (method === 'delete' && id != null) {
    db[key] = list.filter((x) => x.id !== id);
    return respond(config, 200, { message: 'Eliminado' });
  }

  return respond(config, 405, { message: 'Método no soportado en mock' });
}
