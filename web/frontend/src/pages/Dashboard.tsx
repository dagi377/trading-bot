import React, { useState, useEffect } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { format } from 'date-fns';

interface DashboardProps {}

// Mock data for demonstration
const generateMockData = () => {
  const data = [];
  const now = new Date();
  for (let i = 30; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    
    // Generate some random PnL data
    const dailyPnL = Math.random() * 200 - 100; // Between -100 and 100
    const cumulativePnL = Math.random() * 500 - 100; // Between -100 and 400
    
    data.push({
      date: format(date, 'MMM dd'),
      dailyPnL,
      cumulativePnL,
    });
  }
  return data;
};

const Dashboard: React.FC<DashboardProps> = () => {
  const [pnlData, setPnlData] = useState<any[]>([]);
  const [activeGroups, setActiveGroups] = useState<any[]>([]);
  const [openPositions, setOpenPositions] = useState<any[]>([]);
  
  useEffect(() => {
    // In a real app, this would fetch data from the API
    setPnlData(generateMockData());
    
    // Mock active groups
    setActiveGroups([
      { id: '1', name: 'Tech Stocks', capital: 5000, allocated: 3500, stocks: 4, activeTrades: 2, dailyPnL: 120.50 },
      { id: '2', name: 'Financial Sector', capital: 3000, allocated: 2000, stocks: 3, activeTrades: 1, dailyPnL: -45.75 },
      { id: '3', name: 'Energy Plays', capital: 2500, allocated: 1500, stocks: 2, activeTrades: 0, dailyPnL: 0 },
    ]);
    
    // Mock open positions
    setOpenPositions([
      { id: '1', symbol: 'NVDA', quantity: 10, entryPrice: 450.25, currentPrice: 462.75, pnL: 125.00, group: 'Tech Stocks' },
      { id: '2', name: 'SHOP', quantity: 15, entryPrice: 75.50, currentPrice: 73.25, pnL: -33.75, group: 'Tech Stocks' },
      { id: '3', name: 'TD', quantity: 20, entryPrice: 82.30, currentPrice: 80.00, pnL: -46.00, group: 'Financial Sector' },
    ]);
  }, []);
  
  const totalDailyPnL = activeGroups.reduce((sum, group) => sum + group.dailyPnL, 0);
  const totalAllocated = activeGroups.reduce((sum, group) => sum + group.allocated, 0);
  const totalCapital = activeGroups.reduce((sum, group) => sum + group.capital, 0);
  
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="flex space-x-2">
          <button className="btn btn-primary">New Trading Group</button>
          <button className="btn btn-secondary">Add Stock</button>
        </div>
      </div>
      
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className={`card ${totalDailyPnL >= 0 ? 'border-l-4 border-success-500' : 'border-l-4 border-danger-500'}`}>
          <h2 className="text-lg font-semibold mb-2">Daily P&L</h2>
          <p className={`text-2xl font-bold ${totalDailyPnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
            ${totalDailyPnL.toFixed(2)}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Across all trading groups
          </p>
        </div>
        
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Open Positions</h2>
          <p className="text-2xl font-bold">{openPositions.length}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Active trades
          </p>
        </div>
        
        <div className="card">
          <h2 className="text-lg font-semibold mb-2">Capital Allocation</h2>
          <p className="text-2xl font-bold">${totalAllocated.toFixed(2)}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            of ${totalCapital.toFixed(2)} total capital
          </p>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5 mt-2">
            <div 
              className="bg-primary-600 h-2.5 rounded-full" 
              style={{ width: `${(totalAllocated / totalCapital) * 100}%` }}
            ></div>
          </div>
        </div>
      </div>
      
      {/* P&L Chart */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">P&L Performance</h2>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart
              data={pnlData}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis />
              <Tooltip />
              <Area 
                type="monotone" 
                dataKey="dailyPnL" 
                stroke="#0ea5e9" 
                fill="#0ea5e9" 
                fillOpacity={0.3} 
                name="Daily P&L"
                unit="$"
              />
              <Area 
                type="monotone" 
                dataKey="cumulativePnL" 
                stroke="#22c55e" 
                fill="#22c55e" 
                fillOpacity={0.3}
                name="Cumulative P&L"
                unit="$"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>
      
      {/* Active Trading Groups */}
      <div>
        <h2 className="text-lg font-semibold mb-4">Active Trading Groups</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Name
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Capital
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Stocks
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Active Trades
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Daily P&L
                </th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {activeGroups.map((group) => (
                <tr key={group.id}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900 dark:text-white">{group.name}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900 dark:text-white">${group.capital.toFixed(2)}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">${group.allocated.toFixed(2)} allocated</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    {group.stocks}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    {group.activeTrades}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className={`text-sm ${group.dailyPnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                      ${group.dailyPnL.toFixed(2)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button className="text-primary-600 hover:text-primary-900 dark:hover:text-primary-400 mr-3">View</button>
                    <button className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-200">Edit</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
      
      {/* Open Positions */}
      <div>
        <h2 className="text-lg font-semibold mb-4">Open Positions</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Symbol
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Group
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Quantity
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Entry Price
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Current Price
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
              {openPositions.map((position) => (
                <tr key={position.id}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900 dark:text-white">{position.symbol}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    {position.group}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    {position.quantity}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    ${position.entryPrice.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                    ${position.currentPrice.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className={`text-sm ${position.pnL >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                      ${position.pnL.toFixed(2)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button className="text-primary-600 hover:text-primary-900 dark:hover:text-primary-400 mr-3">Details</button>
                    <button className="text-danger-600 hover:text-danger-900 dark:hover:text-danger-400">Close</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
