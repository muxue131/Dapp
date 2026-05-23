import React from 'react';
import { useWallet } from '../hooks/useWallet';

interface WalletConnectProps {
  onConnect?: (address: string) => void;
}

export const WalletConnect: React.FC<WalletConnectProps> = ({ onConnect }) => {
  const { wallet, loading, error, connect, disconnect, isKeplrAvailable } = useWallet();

  const handleConnect = async () => {
    await connect();
    if (wallet.connected && onConnect) {
      onConnect(wallet.address);
    }
  };

  if (!isKeplrAvailable()) {
    return (
      <div className="wallet-install">
        <p>请安装 Keplr 钱包扩展</p>
        <a
          href="https://chrome.google.com/webstore/detail/keplr/dmkamcknogkgcdfhhbddcghachkejeap"
          target="_blank"
          rel="noopener noreferrer"
          className="btn btn-secondary"
        >
          安装 Keplr
        </a>
      </div>
    );
  }

  if (wallet.connected) {
    return (
      <div className="wallet-connected">
        <div className="wallet-info">
          <span className="wallet-name">{wallet.name}</span>
          <span className="wallet-address">
            {wallet.address.slice(0, 10)}...{wallet.address.slice(-6)}
          </span>
          <span className="wallet-balance">{wallet.balance} uleg</span>
        </div>
        <button onClick={disconnect} className="btn btn-outline">
          断开连接
        </button>
      </div>
    );
  }

  return (
    <div className="wallet-connect">
      {error && <div className="error-message">{error}</div>}
      <button onClick={handleConnect} disabled={loading} className="btn btn-primary">
        {loading ? '连接中...' : '连接 Keplr 钱包'}
      </button>
    </div>
  );
};
