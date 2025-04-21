import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

interface TradingGroupsProps {}

// Mock data for demonstration
const mockTradingGroups = [
  { 
    id: '1', 
    name: 'Tech Stocks', 
    capital: 5000, 
    maxLoss: 300, 
    allocated: 3500, 
    stocks: 4, 
    activeTrades: 2, 
    dailyPnL: 120.50,
    status: 'active'
  },
  { 
    id: '2', 
    name: 'Financial Sector', 
    capital: 3000, 
    maxLoss: 200, 
    allocated: 2000, 
    stocks: 3, 
    activeTrades: 1, 
    dailyPnL: -45.75,
    status: 'active'
  },
  { 
    id: '3', 
    name: 'Energy Plays', 
    capital: 2500, 
    maxLoss: 150, 
    allocated: 1500, 
    stocks: 2, 
    activeTrades: 0, 
    dailyPnL: 0,
    status: 'paused'
  },
];

const TradingGroups: React.FC<TradingGroupsProps> = () => {
  const [tradingGroups, setTradingGroups] = useState<any[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newGroup, setNewGroup] = useState({
    name: '',
    capital: 0,
    maxLoss: 0,
  });
  
  useEffect(() => {
    // In a real app, this would fetch data from the API
    setTradingGroups(mockTradingGroups);
  }, []);
  
  const handleCreateGroup = (e: React.FormEvent) => {
    e.preventDefault();
    
    // In a real app, this would send a request to the API
    const newId = (tradingGroups.length + 1).toString();
    const createdGroup = {
      id: newId,
      name: newGroup.name,
      capital: newGroup.capital,
      maxLoss: newGroup.maxLoss,
      allocated: 0,
      stocks: 0,
      activeTrades: 0,
      dailyPnL: 0,
      status: 'active'
    };
    
    setTradingGroups([...tradingGroups, createdGroup]);
    setShowCreateModal(false);
    setNewGroup({ name: '', capital: 0, maxLoss: 0 });
  };
  
  const handleToggleStatus = (id: string) => {
    setTradingGroups(tradingGroups.map(group => {
      if (group.id === id) {
        return {
          ...group,
          status: group.status === 'active' ? 'paused' : 'active'
        };
      }
      return group;
    }));
  };
  
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Trading Groups</h1>
        <button 
          className="btn btn-primary"
          onClick={() => setShowCreateModal(true)}
        >
          Create New Group
        </button>
      </div>
      
      {/* Trading Groups List */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {tradingGroups.map((group) => (
          <div key={group.id} className="card">
            <div className="flex justify-between items-start mb-4">
              <h2 className="text-xl font-semibold">{group.name}</h2>
              <span 
                className={`px-2 py-1 text-xs rounded-full ${
                  group.status === 'active' 
                    ? 'bg-success-100 text-success-800 dark:bg-success-900 dark:text-success-200' 
                    : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
                }`}
              >
                {group.status === 'active' ? 'Active' : 'Paused'}
              </span>
            </div>
            
            <div className="grid grid-cols-2 gap-4 mb-4">
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Capital</p>
                <p className="text-lg font-medium">${group.capital.toFixed(2)}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Max Loss</p>
                <p className="text-lg font-medium">${group.maxLoss.toFixed(2)}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Allocated</p>
                <p className="text-lg font-medium">${group.allocated.toFixed(2)}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Daily P&L</p>
                <p className={`text-lg font-medium ${group.dailyPnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                  ${group.dailyPnL.toFixed(2)}
                </p>
              </div>
            </div>
            
            <div className="flex justify-between items-center">
              <div>
                <span className="text-sm text-gray-500 dark:text-gray-400 mr-4">
                  {group.stocks} stocks
                </span>
                <span className="text-sm text-gray-500 dark:text-gray-400">
                  {group.activeTrades} active trades
                </span>
              </div>
              <div className="flex space-x-2">
                <button 
                  className={`btn ${group.status === 'active' ? 'btn-secondary' : 'btn-success'}`}
                  onClick={() => handleToggleStatus(group.id)}
                >
                  {group.status === 'active' ? 'Pause' : 'Resume'}
                </button>
                <Link to={`/groups/${group.id}`} className="btn btn-primary">
                  View
                </Link>
              </div>
            </div>
          </div>
        ))}
      </div>
      
      {/* Create Group Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
            <div className="p-6">
              <h2 className="text-xl font-semibold mb-4">Create New Trading Group</h2>
              <form onSubmit={handleCreateGroup}>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Group Name
                  </label>
                  <input
                    type="text"
                    className="input"
                    value={newGroup.name}
                    onChange={(e) => setNewGroup({ ...newGroup, name: e.target.value })}
                    required
                  />
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Capital Allocation ($)
                  </label>
                  <input
                    type="number"
                    className="input"
                    value={newGroup.capital || ''}
                    onChange={(e) => setNewGroup({ ...newGroup, capital: parseFloat(e.target.value) })}
                    min="0"
                    step="100"
                    required
                  />
                </div>
                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Max Loss Cap ($)
                  </label>
                  <input
                    type="number"
                    className="input"
                    value={newGroup.maxLoss || ''}
                    onChange={(e) => setNewGroup({ ...newGroup, maxLoss: parseFloat(e.target.value) })}
                    min="0"
                    step="10"
                    required
                  />
                </div>
                <div className="flex justify-end space-x-2">
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={() => setShowCreateModal(false)}
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="btn btn-primary"
                  >
                    Create Group
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TradingGroups;
