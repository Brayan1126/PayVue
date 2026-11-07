import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
import ForgotPassword from './pages/ForgotPassword';
import Dashboard from './pages/Dashboard';
import AddDebt from './pages/AddDebt';
import AddPayment from './pages/AddPayment';
import AddIncome from './pages/AddIncome';
import History from './pages/History';
import ChangePassword from './pages/ChangePassword';

const PrivateRoute = ({ children }) => {
  const user = localStorage.getItem('user');
  return user ? children : <Navigate to="/" />;
};

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/dashboard" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
        <Route path="/add-debt" element={<PrivateRoute><AddDebt /></PrivateRoute>} />
        <Route path="/add-payment" element={<PrivateRoute><AddPayment /></PrivateRoute>} />
        <Route path="/add-income" element={<PrivateRoute><AddIncome /></PrivateRoute>} />
        <Route path="/history" element={<PrivateRoute><History /></PrivateRoute>} />
        <Route path="/change-password" element={<PrivateRoute><ChangePassword /></PrivateRoute>} />
      </Routes>
    </Router>
  );
}

export default App;
