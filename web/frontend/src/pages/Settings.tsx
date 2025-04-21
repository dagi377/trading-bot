import React, { useState } from 'react';
import { useAuth } from '../../hooks/useAuth';

interface SettingsProps {}

const Settings: React.FC<SettingsProps> = () => {
  const { user } = useAuth();
  const [questradeSettings, setQuestradeSettings] = useState({
    clientId: '',
    refreshToken: '',
    isConnected: false
  });
  const [llmSettings, setLLMSettings] = useState({
    provider: 'openai',
    apiKey: '',
    isConnected: false
  });
  const [defaultRiskSettings, setDefaultRiskSettings] = useState({
    defaultMaxLossPerTrade: 30,
    defaultMaxDailyLoss: 400,
    defaultCapitalPerStock: 300
  });
  const [notificationSettings, setNotificationSettings] = useState({
    email: true,
    slack: false,
    slackWebhook: '',
    telegram: false,
    telegramChatId: ''
  });
  const [tradingHours, setTradingHours] = useState({
    start: '09:30',
    end: '16:00',
    timezone: 'America/New_York'
  });
  
  const handleConnectQuestrade = () => {
    // In a real app, this would initiate OAuth flow
    alert('Connecting to Questrade...');
    setQuestradeSettings({
      ...questradeSettings,
      isConnected: true
    });
  };
  
  const handleConnectLLM = () => {
    // In a real app, this would validate the API key
    alert(`Connecting to ${llmSettings.provider}...`);
    setLLMSettings({
      ...llmSettings,
      isConnected: true
    });
  };
  
  const handleSaveSettings = () => {
    // In a real app, this would save settings to the backend
    alert('Settings saved successfully!');
  };
  
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Settings</h1>
        <button 
          className="btn btn-primary"
          onClick={handleSaveSettings}
        >
          Save All Settings
        </button>
      </div>
      
      {/* User Profile */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">User Profile</h2>
        <div className="flex items-center mb-4">
          <div className="w-12 h-12 rounded-full bg-primary-600 flex items-center justify-center text-white text-xl font-bold">
            {user?.name?.charAt(0) || 'U'}
          </div>
          <div className="ml-4">
            <p className="text-lg font-medium">{user?.name || 'User'}</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">{user?.email || 'user@example.com'}</p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Role: {user?.role || 'Trader'}
            </p>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Name
            </label>
            <input
              type="text"
              className="input"
              value={user?.name || ''}
              disabled
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Email
            </label>
            <input
              type="email"
              className="input"
              value={user?.email || ''}
              disabled
            />
          </div>
        </div>
      </div>
      
      {/* Questrade API Integration */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">Questrade API Integration</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Client ID
            </label>
            <input
              type="text"
              className="input"
              value={questradeSettings.clientId}
              onChange={(e) => setQuestradeSettings({ ...questradeSettings, clientId: e.target.value })}
              placeholder="Enter Questrade Client ID"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Refresh Token
            </label>
            <input
              type="password"
              className="input"
              value={questradeSettings.refreshToken}
              onChange={(e) => setQuestradeSettings({ ...questradeSettings, refreshToken: e.target.value })}
              placeholder="Enter Questrade Refresh Token"
            />
          </div>
        </div>
        <div className="flex items-center justify-between">
          <div>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
              questradeSettings.isConnected 
                ? 'bg-success-100 text-success-800 dark:bg-success-900 dark:text-success-200' 
                : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
            }`}>
              {questradeSettings.isConnected ? 'Connected' : 'Not Connected'}
            </span>
          </div>
          <button 
            className="btn btn-primary"
            onClick={handleConnectQuestrade}
          >
            {questradeSettings.isConnected ? 'Reconnect' : 'Connect to Questrade'}
          </button>
        </div>
      </div>
      
      {/* LLM Integration */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">LLM Integration</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              LLM Provider
            </label>
            <select
              className="input"
              value={llmSettings.provider}
              onChange={(e) => setLLMSettings({ ...llmSettings, provider: e.target.value })}
            >
              <option value="openai">OpenAI</option>
              <option value="anthropic">Anthropic</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              API Key
            </label>
            <input
              type="password"
              className="input"
              value={llmSettings.apiKey}
              onChange={(e) => setLLMSettings({ ...llmSettings, apiKey: e.target.value })}
              placeholder={`Enter ${llmSettings.provider === 'openai' ? 'OpenAI' : 'Anthropic'} API Key`}
            />
          </div>
        </div>
        <div className="flex items-center justify-between">
          <div>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
              llmSettings.isConnected 
                ? 'bg-success-100 text-success-800 dark:bg-success-900 dark:text-success-200' 
                : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
            }`}>
              {llmSettings.isConnected ? 'Connected' : 'Not Connected'}
            </span>
          </div>
          <button 
            className="btn btn-primary"
            onClick={handleConnectLLM}
          >
            {llmSettings.isConnected ? 'Reconnect' : 'Connect to LLM'}
          </button>
        </div>
      </div>
      
      {/* Default Risk Parameters */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">Default Risk Parameters</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Default Max Loss Per Trade ($)
            </label>
            <input
              type="number"
              className="input"
              value={defaultRiskSettings.defaultMaxLossPerTrade}
              onChange={(e) => setDefaultRiskSettings({ ...defaultRiskSettings, defaultMaxLossPerTrade: parseFloat(e.target.value) })}
              min="0"
              step="5"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Default Max Daily Loss ($)
            </label>
            <input
              type="number"
              className="input"
              value={defaultRiskSettings.defaultMaxDailyLoss}
              onChange={(e) => setDefaultRiskSettings({ ...defaultRiskSettings, defaultMaxDailyLoss: parseFloat(e.target.value) })}
              min="0"
              step="50"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Default Capital Per Stock ($)
            </label>
            <input
              type="number"
              className="input"
              value={defaultRiskSettings.defaultCapitalPerStock}
              onChange={(e) => setDefaultRiskSettings({ ...defaultRiskSettings, defaultCapitalPerStock: parseFloat(e.target.value) })}
              min="0"
              step="50"
            />
          </div>
        </div>
      </div>
      
      {/* Trading Hours */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">Trading Hours</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Start Time
            </label>
            <input
              type="time"
              className="input"
              value={tradingHours.start}
              onChange={(e) => setTradingHours({ ...tradingHours, start: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              End Time
            </label>
            <input
              type="time"
              className="input"
              value={tradingHours.end}
              onChange={(e) => setTradingHours({ ...tradingHours, end: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Timezone
            </label>
            <select
              className="input"
              value={tradingHours.timezone}
              onChange={(e) => setTradingHours({ ...tradingHours, timezone: e.target.value })}
            >
              <option value="America/New_York">Eastern Time (ET)</option>
              <option value="America/Chicago">Central Time (CT)</option>
              <option value="America/Denver">Mountain Time (MT)</option>
              <option value="America/Los_Angeles">Pacific Time (PT)</option>
            </select>
          </div>
        </div>
      </div>
      
      {/* Notification Settings */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4">Notification Settings</h2>
        <div className="space-y-4">
          <div className="flex items-center">
            <input
              id="email-notifications"
              type="checkbox"
              className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              checked={notificationSettings.email}
              onChange={(e) => setNotificationSettings({ ...notificationSettings, email: e.target.checked })}
            />
            <label htmlFor="email-notifications" className="ml-2 block text-sm text-gray-900 dark:text-gray-100">
              Email Notifications
            </label>
          </div>
          
          <div>
            <div className="flex items-center mb-2">
              <input
                id="slack-notifications"
                type="checkbox"
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                checked={notificationSettings.slack}
                onChange={(e) => setNotificationSettings({ ...notificationSettings, slack: e.target.checked })}
              />
              <label htmlFor="slack-notifications" className="ml-2 block text-sm text-gray-900 dark:text-gray-100">
                Slack Notifications
              </label>
            </div>
            {notificationSettings.slack && (
              <div className="ml-6">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Slack Webhook URL
                </label>
                <input
                  type="text"
                  className="input"
                  value={notificationSettings.slackWebhook}
                  onChange={(e) => setNotificationSettings({ ...notificationSettings, slackWebhook: e.target.value })}
                  placeholder="https://hooks.slack.com/services/..."
                />
              </div>
            )}
          </div>
          
          <div>
            <div className="flex items-center mb-2">
              <input
                id="telegram-notifications"
                type="checkbox"
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                checked={notificationSettings.telegram}
                onChange={(e) => setNotificationSettings({ ...notificationSettings, telegram: e.target.checked })}
              />
              <label htmlFor="telegram-notifications" className="ml-2 block text-sm text-gray-900 dark:text-gray-100">
                Telegram Notifications
              </label>
            </div>
            {notificationSettings.telegram && (
              <div className="ml-6">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Telegram Chat ID
                </label>
                <input
                  type="text"
                  className="input"
                  value={notificationSettings.telegramChatId}
                  onChange={(e) => setNotificationSettings({ ...notificationSettings, telegramChatId: e.target.value })}
                  placeholder="Enter Telegram Chat ID"
                />
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings;
