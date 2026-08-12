import axios from 'axios';
import type {
  LegacyPlan,
  Asset,
  CreatePlanRequest,
  AddAssetRequest,
  HeartbeatRequest,
  ClaimRequest,
} from '../types';

const API_BASE = '/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error.response?.data?.error || error.message || 'Unknown error';
    console.error(`API Error [${error.config?.method?.toUpperCase()} ${error.config?.url}]:`, message);
    return Promise.reject(new Error(message));
  }
);

// === Plan API ===

export async function createPlan(data: CreatePlanRequest): Promise<LegacyPlan> {
  const response = await apiClient.post('/plans', data);
  return response.data.plan;
}

export async function getPlan(planId: number): Promise<LegacyPlan> {
  const response = await apiClient.get(`/plans/${planId}`);
  return response.data.plan;
}

export async function listPlans(creator?: string): Promise<LegacyPlan[]> {
  const params = creator ? { creator } : {};
  const response = await apiClient.get('/plans', { params });
  return response.data.plans;
}

export async function getPlanAssets(planId: number): Promise<Asset[]> {
  const response = await apiClient.get(`/plans/${planId}/assets`);
  return response.data.assets;
}

export async function getPlansByCreator(address: string): Promise<LegacyPlan[]> {
  const response = await apiClient.get(`/creators/${address}/plans`);
  return response.data.plans;
}

// === Asset API ===

export async function addAsset(data: AddAssetRequest): Promise<Asset> {
  const response = await apiClient.post('/assets', data);
  return response.data.asset;
}

export async function getAsset(assetId: number): Promise<Asset> {
  const response = await apiClient.get(`/assets/${assetId}`);
  return response.data;
}

// === Heartbeat API ===

export async function sendHeartbeat(
  planId: number,
  data: HeartbeatRequest
): Promise<{ message: string; trigger_time: string }> {
  const response = await apiClient.post(`/plans/${planId}/heartbeat`, data);
  return response.data;
}

export async function getHeartbeatLogs(planId: number): Promise<any[]> {
  const response = await apiClient.get(`/plans/${planId}/heartbeat-logs`);
  return response.data.logs;
}

// === Claim API ===

export async function claimInheritance(
  planId: number,
  data: ClaimRequest
): Promise<{ message: string }> {
  const response = await apiClient.post(`/plans/${planId}/claim`, data);
  return response.data;
}

// === Status API ===

export async function getStatus(): Promise<any> {
  const response = await apiClient.get('/status');
  return response.data;
}

export async function getHealth(): Promise<any> {
  const response = await apiClient.get('/health', { baseURL: '' });
  return response.data;
}
