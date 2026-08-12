import React, { useState, useEffect, useCallback } from 'react';
import type { LegacyPlan, Asset } from '../types';
import { useWallet } from '../hooks/useWallet';
import { blockchainService } from '../services/blockchain';
import { getPlanAssets } from '../services/api';

interface Props {
  plan: LegacyPlan;
  onRefresh?: () => void;
}

export const PlanCard: React.FC<Props> = ({ plan, onRefresh }) => {
  const { wallet } = useWallet();
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(false);
  const [heartbeatLoading, setHeartbeatLoading] = useState(false);
  const [claimLoading, setClaimLoading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  const isCreator = wallet.address === plan.creator;
  const isBeneficiary = plan.beneficiaries.some((b) => b.address === wallet.address);

  const triggerDate = new Date(plan.trigger_time);
  const lastHeartbeat = new Date(plan.last_heartbeat);
  const now = new Date();
  const timeRemaining = triggerDate.getTime() - now.getTime();
  const daysRemaining = Math.max(0, Math.ceil(timeRemaining / (1000 * 60 * 60 * 24)));

  const loadAssets = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getPlanAssets(plan.plan_id);
      setAssets(data);
    } catch (err) {
      console.error('Failed to load assets:', err);
    } finally {
      setLoading(false);
    }
  }, [plan.plan_id]);

  useEffect(() => {
    loadAssets();
  }, [loadAssets]);

  const handleHeartbeat = async () => {
    if (!wallet.connected) return;
    setHeartbeatLoading(true);
    try {
      const txHash = await blockchainService.sendHeartbeat(wallet.address, plan.plan_id);
      setMessage(`✅ 心跳发送成功！交易: ${txHash}`);
      if (onRefresh) onRefresh();
    } catch (err: any) {
      setMessage(`❌ 心跳发送失败: ${err.message}`);
    } finally {
      setHeartbeatLoading(false);
    }
  };

  const handleClaim = async () => {
    if (!wallet.connected) return;
    setClaimLoading(true);
    try {
      const txHash = await blockchainService.claimInheritance(wallet.address, plan.plan_id);
      setMessage(`✅ 遗产认领成功！交易: ${txHash}`);
      if (onRefresh) onRefresh();
    } catch (err: any) {
      setMessage(`❌ 认领失败: ${err.message}`);
    } finally {
      setClaimLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return '#10b981';
      case 'triggered': return '#f59e0b';
      case 'claimed': return '#6b7280';
      default: return '#6b7280';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active': return '活跃';
      case 'triggered': return '已触发';
      case 'claimed': return '已认领';
      default: return status;
    }
  };

  return (
    <div className="plan-card">
      <div className="plan-header">
        <h3>遗产计划 #{plan.plan_id}</h3>
        <span
          className="status-badge"
          style={{ backgroundColor: getStatusColor(plan.status) }}
        >
          {getStatusText(plan.status)}
        </span>
      </div>

      {message && (
        <div className={message.startsWith('✅') ? 'success-message' : 'error-message'}>
          {message}
        </div>
      )}

      <div className="plan-details">
        <div className="detail-row">
          <span className="label">创建者:</span>
          <span className="value">{plan.creator}</span>
        </div>
        <div className="detail-row">
          <span className="label">心跳间隔:</span>
          <span className="value">{Math.floor(plan.heartbeat_interval / 86400)} 天</span>
        </div>
        <div className="detail-row">
          <span className="label">上次心跳:</span>
          <span className="value">{lastHeartbeat.toLocaleString('zh-CN')}</span>
        </div>
        <div className="detail-row">
          <span className="label">触发时间:</span>
          <span className="value">
            {triggerDate.toLocaleString('zh-CN')}
            {plan.status === 'active' && (
              <span className={daysRemaining <= 7 ? 'warning' : ''}>
                {' '}({daysRemaining} 天后)
              </span>
            )}
          </span>
        </div>
      </div>

      <div className="beneficiaries-section">
        <h4>受益人</h4>
        {plan.beneficiaries.map((b, index) => (
          <div key={index} className="beneficiary-item">
            <span className="address">{b.address}</span>
            <span className="share">{(b.share * 100).toFixed(1)}%</span>
          </div>
        ))}
      </div>

      {assets.length > 0 && (
        <div className="assets-section">
          <h4>资产 ({assets.length})</h4>
          {assets.map((asset, index) => (
            <div key={index} className="asset-item">
              <span className="asset-type">{asset.asset_type}</span>
              <span className="asset-denom">{asset.denom}</span>
              <span className="asset-amount">{asset.amount}</span>
            </div>
          ))}
        </div>
      )}

      <div className="plan-actions">
        {isCreator && plan.status === 'active' && (
          <button
            onClick={handleHeartbeat}
            disabled={heartbeatLoading}
            className="btn btn-primary"
          >
            {heartbeatLoading ? '发送中...' : '发送心跳'}
          </button>
        )}

        {isBeneficiary && plan.status === 'triggered' && (
          <button
            onClick={handleClaim}
            disabled={claimLoading}
            className="btn btn-success"
          >
            {claimLoading ? '认领中...' : '认领遗产'}
          </button>
        )}
      </div>
    </div>
  );
};
