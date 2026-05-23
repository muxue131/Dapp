import React, { useState } from 'react';
import { useWallet } from '../hooks/useWallet';
import { blockchainService } from '../services/blockchain';
import type { AssetType } from '../types';

interface Props {
  planId: number;
  onAssetAdded?: () => void;
}

export const AddAsset: React.FC<Props> = ({ planId, onAssetAdded }) => {
  const { wallet } = useWallet();
  const [assetType, setAssetType] = useState<AssetType>('native');
  const [denom, setDenom] = useState('uleg');
  const [amount, setAmount] = useState('');
  const [ipfsCid, setIpfsCid] = useState('');
  const [metadata, setMetadata] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [txHash, setTxHash] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet.connected) {
      setError('请先连接钱包');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const hash = await blockchainService.addAsset(
        wallet.address,
        planId,
        assetType,
        denom,
        amount,
        ipfsCid || undefined
      );

      setTxHash(hash);
      if (onAssetAdded) onAssetAdded();
    } catch (err: any) {
      setError(err.message || '添加资产失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="add-asset">
      <h3>添加资产到计划 #{planId}</h3>

      {txHash && (
        <div className="success-message">
          ✅ 资产添加成功！交易哈希: {txHash}
        </div>
      )}

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>资产类型</label>
          <select
            value={assetType}
            onChange={(e) => setAssetType(e.target.value as AssetType)}
            className="input"
          >
            <option value="native">原生代币</option>
            <option value="cw20">CW20 代币</option>
            <option value="nft">NFT</option>
            <option value="ipfs">IPFS 文档</option>
          </select>
        </div>

        <div className="form-group">
          <label>
            {assetType === 'native' ? '代币面额' :
             assetType === 'cw20' ? '合约地址' :
             assetType === 'nft' ? 'NFT 合约' : '文档名称'}
          </label>
          <input
            type="text"
            value={denom}
            onChange={(e) => setDenom(e.target.value)}
            placeholder="uleg"
            className="input"
          />
        </div>

        <div className="form-group">
          <label>数量</label>
          <input
            type="text"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            placeholder="1000000"
            className="input"
          />
        </div>

        {(assetType === 'ipfs' || assetType === 'nft') && (
          <div className="form-group">
            <label>IPFS CID</label>
            <input
              type="text"
              value={ipfsCid}
              onChange={(e) => setIpfsCid(e.target.value)}
              placeholder="Qm..."
              className="input"
            />
          </div>
        )}

        <div className="form-group">
          <label>元数据 (JSON)</label>
          <textarea
            value={metadata}
            onChange={(e) => setMetadata(e.target.value)}
            placeholder='{"name": "My Asset", "description": "..."}'
            className="input"
            rows={3}
          />
        </div>

        <button type="submit" disabled={loading} className="btn btn-primary">
          {loading ? '添加中...' : '添加资产'}
        </button>
      </form>
    </div>
  );
};
