import React, { useState, useEffect } from 'react';
import { format } from 'date-fns';

interface TradeHistoryProps {}

// Mock data for demonstration
const generateMockTradeHistory = () => {
  const trades = [];
  const now = new Date();
  const symbols = ['NVDA', 'SHOP', 'AAPL', 'MSFT', 'TD'];
  const groups = ['Tech Stocks', 'Financial Sector', 'Energy Plays'];
  
  for (let i = 0; i < 50; i++) {
    const date = new Date(now);
    date.setHours(date.getHours() - i * 2);
    
    const symbol = symbols[Math.floor(Math.random() * symbols.length)];
    const group = groups[Math.floor(Math.random() * groups.length)];
    const quantity = Math.floor(Math.random() * 20) + 5;
    const entryPrice = Math.random() * 500 + 50;
    const exitPrice = entryPrice * (1 + (Math.random() * 0.1 - 0.05));
    const pnl = (exitPrice - entryPrice) * quantity;
    const fees = quantity * 0.01;
    const type = Math.random() > 0.5 ? 'BUY' : 'SELL';
    
    trades.push({
      id: `trade-${i}`,
      timestamp: date,
      symbol,
      group,
      quantity,
      entryPrice,
      exitPrice,
      pnl,
      fees,
      type,
      reason: `${type === 'BUY' ? 'Bought' : 'Sold'} based on ${Math.random() > 0.5 ? 'RSI' : 'MA'} crossover and LLM recommendation.`
    });
  }
  
  return trades;
};

const TradeHistory: React.FC<TradeHistoryProps> = () => {
  const [trades, setTrades] = useState<any[]>([]);
  const [filteredTrades, setFilteredTrades] = useState<any[]>([]);
  const [filters, setFilters] = useState({
    dateRange: '7d',
    symbol: '',
    group: '',
    type: '',
    result: ''
  });
  
  useEffect(() => {
    // In a real app, this would fetch data from the API
    setTrades(generateMockTradeHistory());
  }, []);
  
  useEffect(() => {
    // Apply filters
    let filtered = [...trades];
    
    // Date range filter
    if (filters.dateRange) {
      const now = new Date();
      let cutoff = new Date();
      
      switch (filters.dateRange) {
        case '1d':
          cutoff.setDate(now.getDate() - 1);
          break;
        case '7d':
          cutoff.setDate(now.getDate() - 7);
          break;
        case '30d':
          cutoff.setDate(now.getDate() - 30);
          break;
        case '90d':
          cutoff.setDate(now.getDate() - 90);
          break;
        default:
          // No date filter
          break;
      }
      
      filtered = filtered.filter(trade => trade.timestamp >= cutoff);
    }
    
    // Symbol filter
    if (filters.symbol) {
      filtered = filtered.filter(trade => trade.symbol === filters.symbol);
    }
    
    // Group filter
    if (filters.group) {
      filtered = filtered.filter(trade => trade.group === filters.group);
    }
    
    // Type filter
    if (filters.type) {
      filtered = filtered.filter(trade => trade.type === filters.type);
    }
    
    // Result filter
    if (filters.result) {
      if (filters.result === 'profit') {
        filtered = filtered.filter(trade => trade.pnl > 0);
      } else if (filters.result === 'loss') {
        filtered = filtered.filter(trade => trade.pnl < 0);
      }
    }
    
    setFilteredTrades(filtered);
  }, [trades, filters]);
  
  const uniqueSymbols = Array.from(new Set(trades.map(trade => trade.symbol)));
  const uniqueGroups = Array.from(new Set(trades.map(trade => trade.group)));
  
  const totalPnL = filteredTrades.reduce((sum, trade) => sum + trade.pnl, 0);
  const totalFees = filteredTrades.reduce((sum, trade) => sum + trade.fees, 0);
  const netPnL = totalPnL - totalFees;
  const profitableTrades = filteredTrades.filter(trade => trade.pnl > 0).length;
  const winRate = filteredTrades.length > 0 ? (profitableTrades / filteredTrades.length) * 100 : 0;
  
  const handleExport = (format: 'csv' | 'json') => {
    // In a real app, this would generate and download the file
    alert(`Exporting trade history as ${format.toUpperCase()}`);
  };
  
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Trade History</h1>
        <div className="flex space-x-2">
          <button 
            className="btn btn-secondary"
            onClick={() => handleExport('csv')}
          >
            Export CSV
          </button>
          <button 
            className="btn btn-secondary"
            onClick={() => handleExport('json')}
          >
            Export JSON
          </button>
        </div>
      </div>
      
      {/* Filters */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">Filters</h2>
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Date Range
            </label>
            <select
              className="input"
              value={filters.dateRange}
              onChange={(e) => setFilters({ ...filters, dateRange: e.target.value })}
            >
              <option value="">All Time</option>
              <option value="1d">Last 24 Hours</option>
              <option value="7d">Last 7 Days</option>
              <option value="30d">Last 30 Days</option>
              <option value="90d">Last 90 Days</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Symbol
            </label>
            <select
              className="input"
              value={filters.symbol}
              onChange={(e) => setFilters({ ...filters, symbol: e.target.value })}
            >
              <option value="">All Symbols</option>
              {uniqueSymbols.map(symbol => (
                <option key={symbol} value={symbol}>{symbol}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Trading Group
            </label>
            <select
              className="input"
              value={filters.group}
              onChange={(e) => setFilters({ ...filters, group: e.target.value })}
            >
              <option value="">All Groups</option>
              {uniqueGroups.map(group => (
                <option key={group} value={group}>{group}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Trade Type
            </label>
            <select
              className="input"
              value={filters.type}
              onChange={(e) => setFilters({ ...filters, type: e.target.value })}
            >
              <option value="">All Types</option>
              <option value="BUY">Buy</option>
              <option value="SELL">Sell</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Result
            </label>
            <select
              className="input"
              value={filters.result}
              onChange={(e) => setFilters({ ...filters, result: e.target.value })}
            >
              <option value="">All Results</option>
              <option value="profit">Profit</option>
              <option value="loss">Loss</option>
            </select>
          </div>
        </div>
      </div>
      
      {/* Summary */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Total Trades</h2>
          <p className="text-2xl font-bold">{filteredTrades.length}</p>
        </div>
        <div className={`card ${netPnL >= 0 ? 'border-l-4 border-success-500' : 'border-l-4 border-danger-500'}`}>
          <h2 className="text-lg font-semibold mb-2">Net P&L</h2>
          <p className={`text-2xl font-bold ${netPnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
            ${netPnL.toFixed(2)}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            After ${totalFees.toFixed(2)} in fees
          </p>
        </div>
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Win Rate</h2>
          <p className="text-2xl font-bold">{winRate.toFixed(1)}%</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {profitableTrades} profitable trades
          </p>
        </div>
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Average P&L</h2>
          <p className={`text-2xl font-bold ${(netPnL / filteredTrades.length) >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
            ${filteredTrades.length > 0 ? (netPnL / filteredTrades.length).toFixed(2) : '0.00'}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Per trade
          </p>
        </div>
      </div>
      
      {/* Trade History Table */}
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Date & Time
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Symbol
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Group
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Type
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Quantity
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Entry Price
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Exit Price
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                P&L
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Fees
              </th>
              <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {filteredTrades.map((trade) => (
              <tr key={trade.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                  {format(trade.timestamp, 'MMM dd, yyyy HH:mm')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                  {trade.symbol}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                  {trade.group}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs rounded-full ${
                    trade.type === 'BUY' 
                      ? 'bg-success-100 text-success-800 dark:bg-success-900 dark:text-success-200' 
                      : 'bg-danger-100 text-danger-800 dark:bg-danger-900 dark:text-danger-200'
                  }`}>
                    {trade.type}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                  {trade.quantity}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                  ${trade.entryPrice.toFixed(2)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                  ${trade.exitPrice.toFixed(2)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`text-sm ${trade.pnl >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                    ${trade.pnl.toFixed(2)}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                  ${trade.fees.toFixed(2)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button 
                    className="text-primary-600 hover:text-primary-900 dark:hover:text-primary-400"
                    onClick={() => alert(`Trade details for ${trade.id}`)}
                  >
                    Details
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default TradeHistory;
