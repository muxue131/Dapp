import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { WalletConnect } from './components/WalletConnect';
import { Dashboard } from './components/Dashboard';
import { CreatePlan } from './components/CreatePlan';
import { useWallet } from './hooks/useWallet';

const Header: React.FC = () => {
  const { wallet } = useWallet();

  return (
    <header className="header">
      <div className="header-content">
        <Link to="/" className="logo">
          <h1>🏛️ Legacy DApp</h1>
          <span className="subtitle">去中心化遗产管理</span>
        </Link>
        <nav className="nav">
          <Link to="/" className="nav-link">仪表盘</Link>
          <Link to="/create" className="nav-link">创建计划</Link>
          {wallet.connected && (
            <Link to="/beneficiary" className="nav-link">受益人视图</Link>
          )}
        </nav>
        <WalletConnect />
      </div>
    </header>
  );
};

const HomePage: React.FC = () => {
  const { wallet } = useWallet();

  return (
    <div className="home-page">
      {!wallet.connected ? (
        <div className="hero">
          <h2>安全、去中心化的遗产管理</h2>
          <p>
            Legacy DApp 使用区块链技术确保您的数字资产能够安全地传递给您的受益人。
            通过心跳机制和 Shamir 秘密共享，您的遗产将在满足条件时自动触发分配。
          </p>
          <div className="features">
            <div className="feature">
              <h3>🔐 加密安全</h3>
              <p>AES-256-GCM 加密保护您的资产信息，Shamir 秘密共享确保密钥安全分发。</p>
            </div>
            <div className="feature">
              <h3>💓 心跳监控</h3>
              <p>定期发送心跳证明存活。心跳过期后，遗产将自动触发分配给受益人。</p>
            </div>
            <div className="feature">
              <h3>👥 多受益人支持</h3>
              <p>支持多个受益人，可自定义每个受益人的份额比例。</p>
            </div>
            <div className="feature">
              <h3>📦 去中心化存储</h3>
              <p>加密文档存储在 IPFS 上，确保数据的持久性和去中心化。</p>
            </div>
          </div>
        </div>
      ) : (
        <Dashboard />
      )}
    </div>
  );
};

const CreatePlanPage: React.FC = () => {
  return (
    <div className="page">
      <CreatePlan />
    </div>
  );
};

const BeneficiaryPage: React.FC = () => {
  const { wallet } = useWallet();

  return (
    <div className="page">
      <h2>受益人视图</h2>
      <p>地址: {wallet.address}</p>
      <p className="help-text">
        在此页面中，您可以查看您作为受益人的所有遗产计划。
        当计划被触发时，您可以认领属于您的遗产份额。
      </p>
      {/* In production: query plans where current user is a beneficiary */}
      <div className="empty-state">
        <p>暂无可认领的遗产。</p>
      </div>
    </div>
  );
};

const App: React.FC = () => {
  return (
    <Router>
      <div className="app">
        <Header />
        <main className="main">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/create" element={<CreatePlanPage />} />
            <Route path="/beneficiary" element={<BeneficiaryPage />} />
          </Routes>
        </main>
        <footer className="footer">
          <p>Legacy DApp v0.1.0 | 基于 Cosmos SDK 构建</p>
        </footer>
      </div>
    </Router>
  );
};

export default App;
