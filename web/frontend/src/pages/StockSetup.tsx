import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface StockSetupProps {}

// Mock data for demonstration
const mockStockData = {
  'NVDA': {
    symbol: 'NVDA',
    name: 'NVIDIA Corporation',
    currentPrice: 462.75,
    dailyChange: 2.5,
    dailyHigh: 465.20,
    dailyLow: 450.10,
    volume: 25000000,
    indicators: {
      rsi: 65.4,
      sma20: 445.30,
      ema50: 430.25,
      volumeSurge: 15.2
    },
    historicalPrices: Array.from({ length: 30 }, (_, i) => ({
      date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      price: 400 + Math.random() * 100,
      volume: 15000000 + Math.random() * 20000000
    }))
  },
  'SHOP': {
    symbol: 'SHOP',
    name: 'Shopify Inc.',
    currentPrice: 73.25,
    dailyChange: -3.0,
    dailyHigh: 76.50,
    dailyLow: 72.80,
    volume: 12000000,
    indicators: {
      rsi: 42.8,
      sma20: 75.40,
      ema50: 72.15,
      volumeSurge: -5.3
    },
    historicalPrices: Array.from({ length: 30 }, (_, i) => ({
      date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      price: 65 + Math.random() * 20,
      volume: 8000000 + Math.random() * 10000000
    }))
  },
  'AAPL': {
    symbol: 'AAPL',
    name: 'Apple Inc.',
    currentPrice: 175.50,
    dailyChange: 0.8,
    dailyHigh: 176.20,
    dailyLow: 174.30,
    volume: 35000000,
    indicators: {
      rsi: 55.2,
      sma20: 172.40,
      ema50: 170.35,
      volumeSurge: 3.7
    },
    historicalPrices: Array.from({ length: 30 }, (_, i) => ({
      date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      price: 165 + Math.random() * 15,
      volume: 30000000 + Math.random() * 15000000
    }))
  },
  'MSFT': {
    symbol: 'MSFT',
    name: 'Microsoft Corporation',
    currentPrice: 410.25,
    dailyChange: 1.2,
    dailyHigh: 412.50,
    dailyLow: 407.80,
    volume: 28000000,
    indicators: {
      rsi: 58.7,
      sma20: 405.60,
      ema50: 400.45,
      volumeSurge: 7.2
    },
    historicalPrices: Array.from({ length: 30 }, (_, i) => ({
      date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      price: 390 + Math.random() * 30,
      volume: 25000000 + Math.random() * 12000000
    }))
  }
};

// Mock LLM insights
const mockLLMInsights = {
  'NVDA': {
    signal: 'BUY',
    confidence: 0.85,
    rationale: 'NVIDIA shows strong momentum with RSI at 65.4, indicating bullish sentiment without being overbought. The stock is trading above both its 20-day SMA and 50-day EMA, confirming an uptrend. Volume surge of 15.2% suggests increasing buying interest. The semiconductor sector remains strong with AI demand driving growth. Recommend buying with a tight stop loss at $450.'
  },
  'SHOP': {
    signal: 'HOLD',
    confidence: 0.65,
    rationale: 'Shopify is showing mixed signals. RSI at 42.8 indicates neither overbought nor oversold conditions. The stock is trading below its 20-day SMA but above its 50-day EMA, suggesting consolidation. Volume is declining (-5.3%), which could indicate waning selling pressure. E-commerce sector faces headwinds from consumer spending concerns, but Shopify\'s recent product innovations provide potential upside. Recommend holding existing positions but not adding new exposure.'
  },
  'AAPL': {
    signal: 'HOLD',
    confidence: 0.70,
    rationale: 'Apple is showing neutral technical indicators with RSI at 55.2. The stock is trading above both its 20-day SMA and 50-day EMA, indicating a modest uptrend. Volume is slightly elevated (+3.7%) but not significantly. Recent product announcements have been incremental rather than revolutionary. The stock appears fairly valued at current levels. Recommend holding existing positions and monitoring upcoming earnings for potential catalysts.'
  },
  'MSFT': {
    signal: 'BUY',
    confidence: 0.80,
    rationale: 'Microsoft exhibits positive momentum with RSI at 58.7, indicating bullish sentiment with room to run. The stock is trading above both its 20-day SMA and 50-day EMA, confirming an uptrend. Volume increase of 7.2% suggests growing interest. Cloud business continues to show strong growth, and AI integration across product lines positions the company well. Recommend buying with a stop loss at $400.'
  }
};

const StockSetup: React.FC<StockSetupProps> = () => {
  const { symbol } = useParams<{ symbol?: string }>();
  const [stockData, setStockData] = useState<any | null>(null);
  const [llmInsight, setLLMInsight] = useState<any | null>(null);
  const [searchSymbol, setSearchSymbol] = useState('');
  const [stockSettings, setStockSettings] = useState({
    maxCapital: 0,
    maxLoss: 0,
    tradingHoursStart: '09:30',
    tradingHoursEnd: '16:00',
  });
  
  useEffect(() => {
    if (symbol && mockStockData[symbol as keyof typeof mockStockData]) {
      setStockData(mockStockData[symbol as keyof typeof mockStockData]);
      setLLMInsight(mockLLMInsights[symbol as keyof typeof mockLLMInsights]);
      
      // Set default settings
      setStockSettings({
        maxCapital: 1000,
        maxLoss: 50,
        tradingHoursStart: '09:30',
        tradingHoursEnd: '16:00',
      });
    } else if (symbol) {
      // Handle unknown symbol
      setStockData(null);
      setLLMInsight(null);
    }
  }, [symbol]);
  
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchSymbol && mockStockData[searchSymbol as keyof typeof mockStockData]) {
      setStockData(mockStockData[searchSymbol as keyof typeof mockStockData]);
      setLLMInsight(mockLLMInsights[searchSymbol as keyof typeof mockLLMInsights]);
      
      // Set default settings
      setStockSettings({
        maxCapital: 1000,
        maxLoss: 50,
        tradingHoursStart: '09:30',
        tradingHoursEnd: '16:00',
      });
    } else {
      // Handle unknown symbol
      setStockData(null);
      setLLMInsight(null);
    }
  };
  
  const handleSaveSettings = () => {
    // In a real app, this would send a request to the API
    alert(`Settings saved for ${stockData.symbol}`);
  };
  
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Stock Setup</h1>
        {stockData && (
          <button 
            className="btn btn-primary"
            onClick={handleSaveSettings}
          >
            Save Settings
          </button>
        )}
      </div>
      
      {/* Stock Search */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">Search Stock</h2>
        <form onSubmit={handleSearch} className="flex space-x-2">
          <input
            type="text"
            className="input flex-1"
            placeholder="Enter stock symbol (e.g., NVDA, SHOP, AAPL, MSFT)"
            value={searchSymbol}
            onChange={(e) => setSearchSymbol(e.target.value.toUpperCase())}
          />
          <button type="submit" className="btn btn-primary">
            Search
          </button>
        </form>
      </div>
      
      {stockData ? (
        <>
          {/* Stock Overview */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="card col-span-2">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h2 className="text-xl font-semibold">{stockData.symbol}</h2>
                  <p className="text-gray-500 dark:text-gray-400">{stockData.name}</p>
                </div>
                <div className="text-right">
                  <p className="text-2xl font-bold">${stockData.currentPrice.toFixed(2)}</p>
                  <p className={`text-sm ${stockData.dailyChange >= 0 ? 'text-success-600' : 'text-danger-600'}`}>
                    {stockData.dailyChange >= 0 ? '+' : ''}{stockData.dailyChange.toFixed(2)}%
                  </p>
                </div>
              </div>
              
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-4">
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Daily High</p>
                  <p className="text-lg font-medium">${stockData.dailyHigh.toFixed(2)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Daily Low</p>
                  <p className="text-lg font-medium">${stockData.dailyLow.toFixed(2)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Volume</p>
                  <p className="text-lg font-medium">{(stockData.volume / 1000000).toFixed(1)}M</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400">RSI</p>
                  <p className={`text-lg font-medium ${
                    stockData.indicators.rsi > 70 ? 'text-danger-600' : 
                    stockData.indicators.rsi < 30 ? 'text-success-600' : 
                    'text-gray-900 dark:text-white'
                  }`}>
                    {stockData.indicators.rsi.toFixed(1)}
                  </p>
                </div>
              </div>
              
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart
                    data={stockData.historicalPrices}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="date" 
                      tickFormatter={(date) => {
                        const d = new Date(date);
                        return `${d.getMonth() + 1}/${d.getDate()}`;
                      }}
                    />
                    <YAxis domain={['auto', 'auto']} />
                    <Tooltip 
                      formatter={(value: any) => [`$${parseFloat(value).toFixed(2)}`, 'Price']}
                      labelFormatter={(label) => `Date: ${label}`}
                    />
                    <Line 
                      type="monotone" 
                      dataKey="price" 
                      stroke="#0ea5e9" 
                      dot={false}
                      activeDot={{ r: 8 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
            
            <div className="card">
              <h2 className="text-lg font-semibold mb-4">Trading Settings</h2>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Max Capital ($)
                  </label>
                  <input
                    type="number"
                    className="input"
                    value={stockSettings.maxCapital}
                    onChange={(e) => setStockSettings({ ...stockSettings, maxCapital: parseFloat(e.target.value) })}
                    min="0"
                    step="100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Max Loss ($)
                  </label>
                  <input
                    type="number"
                    className="input"
                    value={stockSettings.maxLoss}
                    onChange={(e) => setStockSettings({ ...stockSettings, maxLoss: parseFloat(e.target.value) })}
                    min="0"
                    step="10"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Trading Hours (EST)
                  </label>
                  <div className="grid grid-cols-2 gap-2">
                    <input
                      type="time"
                      className="input"
                      value={stockSettings.tradingHoursStart}
                      onChange={(e) => setStockSettings({ ...stockSettings, tradingHoursStart: e.target.value })}
                    />
                    <input
                      type="time"
                      className="input"
                      value={stockSettings.tradingHoursEnd}
                      onChange={(e) => setStockSettings({ ...stockSettings, tradingHoursEnd: e.target.value })}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          {/* Technical Indicators */}
          <div className="card">
            <h2 className="text-lg font-semibold mb-4">Technical Indicators</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">RSI (14)</p>
                <p className={`text-xl font-medium ${
                  stockData.indicators.rsi > 70 ? 'text-danger-600' : 
                  stockData.indicators.rsi < 30 ? 'text-success-600' : 
                  'text-gray-900 dark:text-white'
                }`}>
                  {stockData.indicators.rsi.toFixed(1)}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {stockData.indicators.rsi > 70 ? 'Overbought' : 
                   stockData.indicators.rsi < 30 ? 'Oversold' : 
                   'Neutral'}
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">SMA (20)</p>
                <p className="text-xl font-medium">${stockData.indicators.sma20.toFixed(2)}</p>
                <p className={`text-xs ${stockData.currentPrice > stockData.indicators.sma20 ? 'text-success-600' : 'text-danger-600'}`}>
                  {stockData.currentPrice > stockData.indicators.sma20 ? 'Above' : 'Below'} SMA
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">EMA (50)</p>
                <p className="text-xl font-medium">${stockData.indicators.ema50.toFixed(2)}</p>
                <p className={`text-xs ${stockData.currentPrice > stockData.indicators.ema50 ? 'text-success-600' : 'text-danger-600'}`}>
                  {stockData.currentPrice > stockData.indicators.ema50 ? 'Above' : 'Below'} EMA
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">Volume Surge</p>
                <p className={`text-xl font-medium ${stockData.indicators.volumeSurge > 0 ? 'text-success-600' : 'text-danger-600'}`}>
                  {stockData.indicators.volumeSurge > 0 ? '+' : ''}{stockData.indicators.volumeSurge.toFixed(1)}%
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {Math.abs(stockData.indicators.volumeSurge) > 10 ? 'Significant' : 'Normal'} volume
                </p>
              </div>
            </div>
          </div>
          
          {/* LLM Insights */}
          {llmInsight && (
            <div className={`card border-l-4 ${
              llmInsight.signal === 'BUY' ? 'border-success-500' : 
              llmInsight.signal === 'SELL' ? 'border-danger-500' : 
              'border-gray-500'
            }`}>
              <div className="flex justify-between items-start mb-4">
                <h2 className="text-lg font-semibold">LLM Trading Insights</h2>
                <div className="flex items-center">
                  <span className={`px-3 py-1 text-sm font-medium rounded-full ${
                    llmInsight.signal === 'BUY' ? 'bg-success-100 text-success-800 dark:bg-success-900 dark:text-success-200' : 
                    llmInsight.signal === 'SELL' ? 'bg-danger-100 text-danger-800 dark:bg-danger-900 dark:text-danger-200' : 
                    'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
                  }`}>
                    {llmInsight.signal}
                  </span>
                  <span className="ml-2 text-sm text-gray-500 dark:text-gray-400">
                    {(llmInsight.confidence * 100).toFixed(0)}% confidence
                  </span>
                </div>
              </div>
              <p className="text-gray-700 dark:text-gray-300">
                {llmInsight.rationale}
              </p>
            </div>
          )}
        </>
      ) : (
        <div className="card">
          <div className="text-center py-8">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 mx-auto text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h2 className="mt-2 text-lg font-medium text-gray-900 dark:text-white">No Stock Selected</h2>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Search for a stock symbol or select a stock from a trading group.
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

export default StockSetup;
