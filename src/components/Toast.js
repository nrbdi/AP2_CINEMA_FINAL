import { useState, useEffect, createContext, useContext, useCallback } from 'react';

const ToastContext = createContext(null);

let _addToast = null;

export function toast(msg, type = 'success') {
  if (_addToast) _addToast(msg, type);
}

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([]);

  const addToast = useCallback((msg, type = 'success') => {
    const id = Date.now();
    setToasts(t => [...t, { id, msg, type }]);
    setTimeout(() => setToasts(t => t.filter(x => x.id !== id)), 3500);
  }, []);

  useEffect(() => { _addToast = addToast; }, [addToast]);

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      <div className="toast-container">
        {toasts.map(t => (
          <div key={t.id} className={`toast toast-${t.type}`}>{t.msg}</div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export default function Toast() {
  // Rendered via ToastProvider in App; this is a placeholder if used standalone
  return null;
}
