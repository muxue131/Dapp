import React, { useState } from 'react';
import { useWallet } from '../hooks/useWallet';
import { blockchainService } from '../services/blockchain';
import type { Beneficiary } from '../types';

interface Props {
  onPlanCreated?: (planId: number) => void;
}

export const CreatePlan: React.FC<Props> = ({ onPlanCreated }) => {
  const { wallet } = useWallet();
  const [beneficiaries, setBeneficiaries] = useState<Beneficiary[]>([
    { address: '', share: 1 },
  ]);
  const [heartbeatDays, setHeartbeatDays] = useState(30);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [txHash, setTxHash] = useState<string | null>(null);

  const addBeneficiary = () => {
    setBeneficiaries([...beneficiaries, { address: '', share: 0 }]);
  };

  const removeBeneficiary = (index: number) => {
    if (beneficiaries.length > 1) {
      setBeneficiaries(beneficiaries.filter((_, i) => i !== index));
    }
  };

  const updateBeneficiary = (index: number, field: keyof Beneficiary, value: string | number) => {
    const updated = [...beneficiaries];
    if (field === 'share') {
      updated[index] = { ...updated[index], [field]: Number(value) };
    } else {
      updated[index] = { ...updated[index], [field]: value };
    }
    setBeneficiaries(updated);
  };

  const totalShare = beneficiaries.reduce((sum, b) => sum + b.share, 0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet.connected) {
      setError('请先连接钱包');
      return;
    }

    if (Math.abs(totalShare - 1) > 0.001) {
      setError('受益人份额之和必须等于 1');
      return;
    }

    if (beneficiaries.some((b) => !b.address)) {
      setError('请填写所有受益人地址');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const heartbeatInterval = heartbeatDays * 24 * 60 * 60; // Convert days to seconds
      const txHash = await blockchainService.createLegacyPlan(
        wallet.address,
        beneficiaries.map((b) => ({
          address: b.address,
          share: b.share.toString(),
        })),
        heartbeatInterval
      );

      setTxHash(txHash);
      console.log('Plan created, tx:', txHash);

      if (onPlanCreated) {
        onPlanCreated(0); // In production, parse plan ID from tx events
      }
    } catch (err: any) {
      setError(err.message || '创建计划失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="create-plan">
      <h2>创建遗产计划</h2>

      {txHash && (
        <div className="success-message">
          <p>✅ 计划创建成功！</p>
          <p>交易哈希: {txHash}</p>
        </div>
      )}

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-section">
          <h3>受益人设置</h3>
          {beneficiaries.map((b, index) => (
            <div key={index} className="beneficiary-row">
              <input
                type="text"
                placeholder="受益人地址 (legacy1...)"
                value={b.address}
                onChange={(e) => updateBeneficiary(index, 'address', e.target.value)}
                className="input"
              />
              <input
                type="number"
                placeholder="份额 (0-1)"
                value={b.share}
                step="0.01"
                min="0"
                max="1"
                onChange={(e) => updateBeneficiary(index, 'share', e.target.value)}
                className="input input-small"
              />
              {beneficiaries.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeBeneficiary(index)}
                  className="btn btn-danger btn-sm"
                >
                  删除
                </button>
              )}
            </div>
          ))}

          <div className="share-total">
            总份额: {totalShare.toFixed(2)}{' '}
            {Math.abs(totalShare - 1) < 0.001 ? '✅' : '⚠️ 必须等于 1'}
          </div>

          <button type="button" onClick={addBeneficiary} className="btn btn-secondary">
            + 添加受益人
          </button>
        </div>

        <div className="form-section">
          <h3>心跳设置</h3>
          <label>
            心跳间隔（天）:
            <input
              type="number"
              value={heartbeatDays}
              min="1"
              max="3650"
              onChange={(e) => setHeartbeatDays(Number(e.target.value))}
              className="input input-small"
            />
          </label>
          <p className="help-text">
            每 {heartbeatDays} 天需要发送一次心跳。如果心跳过期，遗产将自动触发分配给受益人。
          </p>
        </div>

        <button
          type="submit"
          disabled={loading || !wallet.connected || Math.abs(totalShare - 1) > 0.001}
          className="btn btn-primary btn-lg"
        >
          {loading ? '创建中...' : '创建遗产计划'}
        </button>
      </form>
    </div>
  );
};
