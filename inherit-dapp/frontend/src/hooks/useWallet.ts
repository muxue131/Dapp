import { useState, useCallback, useEffect } from 'react';
import type { WalletState } from '../types';
import { CHAIN_CONFIG, blockchainService } from '../services/blockchain';

const INITIAL_WALLET_STATE: WalletState = {
  connected: false,
  address: '',
  name: '',
  balance: '0',
};

export function useWallet() {
  const [wallet, setWallet] = useState<WalletState>(INITIAL_WALLET_STATE);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Check if Keplr is installed
  const isKeplrAvailable = useCallback(() => {
    return typeof window !== 'undefined' && !!window.keplr;
  }, []);

  // Connect to Keplr wallet
  const connect = useCallback(async () => {
    if (!isKeplrAvailable()) {
      setError('Keplr wallet not found. Please install the Keplr browser extension.');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Suggest chain to Keplr
      await window.keplr!.experimentalSuggestChain(CHAIN_CONFIG);

      // Enable the chain
      await window.keplr!.enable(CHAIN_CONFIG.chainId);

      // Get the key/account
      const key = await window.keplr!.getKey(CHAIN_CONFIG.chainId);

      // Get the offline signer
      const offlineSigner = window.keplr!.getOfflineSigner(CHAIN_CONFIG.chainId);

      // Connect blockchain service with signer
      await blockchainService.connectWithSigner(offlineSigner);

      // Get balance
      const balance = await blockchainService.getBalance(key.bech32Address);

      setWallet({
        connected: true,
        address: key.bech32Address,
        name: key.name,
        balance: balance,
      });

      console.log('Wallet connected:', key.bech32Address);
    } catch (err: any) {
      console.error('Failed to connect wallet:', err);
      setError(err.message || 'Failed to connect wallet');
    } finally {
      setLoading(false);
    }
  }, [isKeplrAvailable]);

  // Disconnect wallet
  const disconnect = useCallback(() => {
    blockchainService.disconnect();
    setWallet(INITIAL_WALLET_STATE);
  }, []);

  // Refresh balance
  const refreshBalance = useCallback(async () => {
    if (!wallet.connected || !wallet.address) return;

    try {
      const balance = await blockchainService.getBalance(wallet.address);
      setWallet((prev) => ({ ...prev, balance }));
    } catch (err) {
      console.error('Failed to refresh balance:', err);
    }
  }, [wallet.connected, wallet.address]);

  // Listen for Keplr account changes
  useEffect(() => {
    if (!isKeplrAvailable()) return;

    const handleAccountChange = () => {
      if (wallet.connected) {
        connect(); // Re-connect with new account
      }
    };

    window.addEventListener('keplr_keystorechange', handleAccountChange);
    return () => {
      window.removeEventListener('keplr_keystorechange', handleAccountChange);
    };
  }, [wallet.connected, connect, isKeplrAvailable]);

  return {
    wallet,
    loading,
    error,
    connect,
    disconnect,
    refreshBalance,
    isKeplrAvailable,
  };
}
