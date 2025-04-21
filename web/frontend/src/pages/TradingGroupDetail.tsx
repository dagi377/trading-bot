import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';

interface TradingGroupDetailProps {}

// Mock data for demonstration
const mockStocks = [
  { 
    id: '1', 
    symbol: 'NVDA', 
    capital: 1000, 
    maxLoss: 50, 
    allocated: 950, 
    currentPrice: 462.75, 
    dailyChange: 2.5, 
    dailyPnL: 125.00,
    status: 'active',
    position: { quantity: 10, entryPrice: 450.25 }
  },
  { 
    id: '2', 
    symbol: 'SHOP', 
    capital: 1500, 
    maxLoss: 75, 
    allocated: 1132.50, 
    currentPrice: 73.25, 
    dailyChange: -3.0, 
    dailyPnL: -33.75,
    status: 'active',
    position: { quantity: 15, entryPrice: 75.50 }
  },
  { 
    id: '3', 
    symbol: 'AAPL', 
    capital: 1000, 
    maxLoss: 50, 
    allocated: 0, 
    currentPrice: 175.50, 
    dailyChange: 0.8, 
    dailyPnL: 0,
    status: 'paused',
    position: null
  },
  { 
    id: '4', 
    symbol: 'MSFT', 
    capital: 1500, 
    maxLoss: 75, 
    allocated: 0, 
    currentPrice: 410.25, 
    dailyChange: 1.2, 
    dailyPnL: 0,
    status: 'active',
    position: null
  },
];

const mockTradingGroups = [
  { 
    id: '1', 
    name: 'Tech Stocks', 
    capital: 5000, 
    maxLoss: 300, 
    allocated: 3500, 
    stocks: mockStocks,
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
    stocks: [],
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
    stocks: [],
    activeTrades: 0, 
    dailyPnL: 0,
    status: 'paused'
  },
];

const TradingGroupDetail: React.FC<TradingGroupDetailProps> = () => {
  const { id } = useParams<{ id: string }>();
  const [group, setGroup] = useState<any | null>(null);
  const [showAddStockModal, setShowAddStockModal] = useState(false);
  const [newStock, setNewStock] = useState({
    symbol: '',
    capital: 0,
    maxLoss: 0,
  });
  
  useEffect(() => {
    // In a real app, this would fetch data from the API
    const foundGroup = mockTradingGroups.find(g => g.id === id);
    setGroup(foundGroup || null);
  }, [id]);
  
  if (!group) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-500 dark:text-gray-400">Trading group not found</div>
      </div>
    );
  }
  
  const handleAddStock = (e: React.FormEvent) => {
    e.preventDefault();
    
    // In a real app, this would send a request to the API
    const newId = (group.stocks.length + 1).toString();
    const addedStock = {
      id: newId,
      symbol: newStock.symbol,
      capital: newStock.capital,
      maxLoss: newStock.maxLoss,
      allocated: 0,
      currentPrice: 0,
      dailyChange: 0,
      dailyPnL: 0,
      status: 'active',
      position: null
    };
    
    const updatedStocks = [...group.stocks, addedStock];
    setGroup({
      ...group,
      stocks: updatedStocks,
    });
    
    setShowAddStockModal(false);
    setNewStock({ symbol: '', capital: 0, maxLoss: 0 });
  };
  
  const handleToggleStockStatus = (stockId: string) => {
    const updatedStocks = group.stocks.map((stock: any) => {
      if (stock.id === stockId) {
        return {
          ...stock,
          status: stock.status === 'active' ? 'paused' : 'active'
        };
      }
      return stock;
    });
    
    setGroup({
      ...group,
      stocks: updatedStocks,
    });
  };
  
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <Link to="/groups" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 flex items-center mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clipRule="evenodd" />
            </svg>
            Back to Trading Groups
          </Link>
          <h1 className="text-2xl font-bold">{group.name}</h1>
        </div>
        <div className="flex space-x-2">
          <button 
            className={`btn ${group.status === 'active' ? 'btn-secondary' : 'btn-success'}`}
          >
            {group.status === 'active' ? 'Pause Group' : 'Resume Group'}
          </button>
          <button 
            className="btn btn-primary"
            onClick={() => setShowAddStockModal(true)}
          >
            Add Stock
          </button>
        </div>
      </div>
      
      {/* Group Summary */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Capital</h2>
          <p className="text-2xl font-bold">${group.capital.toFixed(2)}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            ${group.allocated.toFixed(2)} allocated
          </p>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5 mt-2">
            <div 
              className="bg-primary-600 h-2.5 rounded-full" 
              style={{ width: `${(group.allocated / group.capital) * 100}%` }}
            ></div>
          </div>
        </div>
        
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Max Loss</h2>
          <p className="text-2xl font-bold">${group.maxLoss.toFixed(2)}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Group loss limit
          </p>
        </div>
        
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Active Trades</h2>
          <p className="text-2xl font-bold">{group.activeTrades}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Out of {group.stocks.length} stocks
          </p>
        </div>
        
        <div className={`card ${group.dailyPnL >= 0 ? 'border-l-4 border-success-500' : 'border-l-4 border-danger-500'}`}>
          <h2 className="text-lg font-semibold mb-2">Daily P&L</h2>
          <p className={`text-2xl font-bold ${group.dailyPnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
            ${group.dailyPnL.toFixed(2)}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Group performance today
          </p>
        </div>
      </div>
      
      {/* Stocks Table */}
      <div>
        <h2 className="text-lg font-semibold mb-4">Stocks</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Symbol
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Status
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Capital
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Max Loss
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Current Price
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Position
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  P&L
                </th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {group.stocks.map((stock: any) => (
                <tr key={stock.id}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900 dark:text-white">{stock.symbol}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span 
                      className={`px-2 py-1 text-xs rounded-full ${
                        stock.status === 'active' 
                          ? 'bg-success-100 text-success-800 dark:bg-success-900 dark:text-success-200' 
                          : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
                      }`}
                    >
                      {stock.status === 'active' ? 'Active' : 'Paused'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900 dark:text-white">${stock.capital.toFixed(2)}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      {stock.allocated > 0 ? `$${stock.allocated.toFixed(2)} allocated` : 'No allocation'}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    ${stock.maxLoss.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900 dark:text-white">${stock.currentPrice.toFixed(2)}</div>
                    <div className={`text-xs ${stock.dailyChange >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                      {stock.dailyChange >= 0 ? '+' : ''}{stock.dailyChange.toFixed(2)}%
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    {stock.position ? (
                      <div>
                        {stock.position.quantity} shares @ ${stock.position.entryPrice.toFixed(2)}
                      </div>
                    ) : (
                      <div className="text-gray-500 dark:text-gray-400">No position</div>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className={`text-sm ${stock.dailyPnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                      ${stock.dailyPnL.toFixed(2)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button 
                      className={`text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-200 mr-3`}
                      onClick={() => handleToggleStockStatus(stock.id)}
                    >
                      {stock.status === 'active' ? 'Pause' : 'Resume'}
                    </button>
                    <Link 
                      to={`/stock-setup/${stock.symbol}`} 
                      className="text-primary-600 hover:text-primary-900 dark:hover:text-primary-400"
                    >
                      Setup
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
      
      {/* Add Stock Modal */}
      {showAddStockModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
            <div className="p-6">
              <h2 className="text-xl font-semibold mb-4">Add Stock to Group</h2>
              <form onSubmit={handleAddStock}>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Stock Symbol
                  </label>
                  <input
                    type="text"
                    className="input"
                    value={newStock.symbol}
                    onChange={(e) => setNewStock({ ...newStock, symbol: e.target.value.toUpperCase() })}
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
                    value={newStock.capital || ''}
                    onChange={(e) => setNewStock({ ...newStock, capital: parseFloat(e.target.value) })}
                    min="0"
                    step="100"
                    required
                  />
                </div>
                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Max Loss ($)
                  </label>
                  <input
                    type="number"
                    className="input"
                    value={newStock.maxLoss || ''}
                    onChange={(e) => setNewStock({ ...newStock, maxLoss: parseFloat(e.target.value) })}
                    min="0"
                    step="10"
                    required
                  />
                </div>
                <div className="flex justify-end space-x-2">
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={() => setShowAddStockModal(false)}
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="btn btn-primary"
                  >
                    Add Stock
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

export default TradingGroupDetail;
