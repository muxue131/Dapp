import React, { useState, useEffect } from 'react';
import { useWallet } from '../hooks/useWallet';
import { getPlansByCreator } from '../services/api';
import { PlanCard } from './PlanCard';
import type { LegacyPlan } from '../types';

export const Dashboard: React.FC = () => {
  const { wallet } = useWallet();
  const [plans, setPlans] = useState<LegacyPlan[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<'all' | 'active' | 'triggered' | 'claimed'>('all');

  const loadPlans = async () => {
    if (!wallet.connected) return;

    setLoading(true);
    setError(null);

    try {
      const response = await getPlansByCreator(wallet.address);
      if (response.plans) {
        setPlans(response.plans);
      }
    } catch (err: any) {
      setError(err.message || '加载计划失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPlans();
  }, [wallet.connected, wallet.address]);

  const filteredPlans = filter === 'all'
    ? plans
    : plans.filter((p) => p.status === filter);

  if (!wallet.connected) {
    return (
      <div className="dashboard">
        <div className="empty-state">
          <h2>欢迎使用 Legacy DApp</h2>
          <p>请连接 Keplr 钱包开始管理您的遗产计划</p>
        </div>
      </div>
    );
  }

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h2>我的遗产计划</h2>
        <div className="filter-tabs">
          {(['all', 'active', 'triggered', 'claimed'] as const).map((f) => (
            <button
              key={f}
              className={`tab ${filter === f ? 'active' : ''}`}
              onClick={() => setFilter(f)}
            >
              {f === 'all' ? '全部' :
               f === 'active' ? '活跃' :
               f === 'triggered' ? '已触发' : '已认领'}
            </button>
          ))}
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <div className="loading">加载中...</div>
      ) : filteredPlans.length === 0 ? (
        <div className="empty-state">
          <p>
            {filter === 'all'
              ? '暂无遗产计划。点击"创建计划"开始。'
              : `暂无${filter === 'active' ? '活跃' : filter === 'triggered' ? '已触发' : '已认领'}的计划。`}
          </p>
        </div>
      ) : (
        <div className="plans-grid">
          {filteredPlans.map((plan) => (
            <PlanCard key={plan.plan_id} plan={plan} onRefresh={loadPlans} />
          ))}
        </div>
      )}
    </div>
  );
};
