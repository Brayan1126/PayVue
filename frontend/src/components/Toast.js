import React from 'react';

function Toast({ message, type }) {
  return (
    <div className={`toast ${type}`}>
      {type === 'success' && (
        <div className="icon">
          <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </div>
      )}
      {type === 'error' && (
        <div className="icon" style={{ color: '#ef4444' }}>
          <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </div>
      )}
      <div>
        <div className="message">{message}</div>
      </div>
    </div>
  );
}

export default Toast;

