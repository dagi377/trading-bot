import React, { useEffect, useState } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import Navbar from './components/layout/Navbar';
import Sidebar from './components/layout/Sidebar';
import { useAuth } from './hooks/useAuth';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Settings from './pages/Settings';
import StockSetup from './pages/StockSetup';
import TradeHistory from './pages/TradeHistory';
import TradingGroupDetail from './pages/TradingGroupDetail';
import TradingGroups from './pages/TradingGroups';
import Register from './pages/auth/Register';

const App: React.FC = () => {
  const [darkMode, setDarkMode] = useState(false);
  const { user, loading } = useAuth();
  
  useEffect(() => {
    // Check user preference for dark mode
    const isDarkMode = localStorage.getItem('darkMode') === 'true';
    setDarkMode(isDarkMode);
    
    if (isDarkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, []);
  
  const toggleDarkMode = () => {
    const newDarkMode = !darkMode;
    setDarkMode(newDarkMode);
    localStorage.setItem('darkMode', newDarkMode.toString());
    
    if (newDarkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  };
  
  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100 dark:bg-gray-900">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100`}>
      {user && <Navbar darkMode={darkMode} toggleDarkMode={toggleDarkMode} />}
      <div className="flex">
        {user && <Sidebar />}
        <main className={`flex-1 p-6 overflow-auto ${!user ? 'w-full' : ''}`}>
          <Routes>
            {!user ? (
              <>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="*" element={<Navigate to="/login" replace />} />
              </>
            ) : (
              <>
                <Route path="/" element={<Dashboard />} />
                <Route path="/groups" element={<TradingGroups />} />
                <Route path="/groups/:id" element={<TradingGroupDetail />} />
                <Route path="/stock-setup/:symbol?" element={<StockSetup />} />
                <Route path="/history" element={<TradeHistory />} />
                <Route path="/settings" element={<Settings />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </>
            )}
          </Routes>
        </main>
      </div>
    </div>
  );
};

export default App;
