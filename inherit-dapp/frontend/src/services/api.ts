import axios from 'axios';
import type {
  LegacyPlan,
  Asset,
  CreatePlanRequest,
  AddAssetRequest,
  HeartbeatRequest,
  ClaimRequest,
  ApiResponse,
} from '../types';

const API_BASE = '/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// === Plan API ===

export async function createPlan(data: CreatePlanRequest): Promise<ApiResponse<LegacyPlan>> {
  const response = await apiClient.post('/plans', data);
  return response.data;
}

export async function getPlan(planId: number): Promise<ApiResponse<LegacyPlan>> {
  const response = await apiClient.get(`/plans/${planId}`);
  return response.data;
}

export async function listPlans(creator?: string): Promise<ApiResponse<LegacyPlan[]>> {
  const params = creator ? { creator } : {};
  const response = await apiClient.get('/plans', { params });
  return response.data;
}

export async function getPlanAssets(planId: number): Promise<ApiResponse<Asset[]>> {
  const response = await apiClient.get(`/plans/${planId}/assets`);
  return response.data;
}

export async function getPlansByCreator(address: string): Promise<ApiResponse<LegacyPlan[]>> {
  const response = await apiClient.get(`/creators/${address}/plans`);
  return response.data;
}

// === Asset API ===

export async function addAsset(data: AddAssetRequest): Promise<ApiResponse<Asset>> {
  const response = await apiClient.post('/assets', data);
  return response.data;
}

export async function getAsset(assetId: number): Promise<ApiResponse<Asset>> {
  const response = await apiClient.get(`/assets/${assetId}`);
  return response.data;
}

// === Heartbeat API ===

export async function sendHeartbeat(
  planId: number,
  data: HeartbeatRequest
): Promise<ApiResponse<{ message: string; trigger_time: string }>> {
  const response = await apiClient.post(`/plans/${planId}/heartbeat`, data);
  return response.data;
}

export async function getHeartbeatLogs(planId: number): Promise<ApiResponse<any[]>> {
  const response = await apiClient.get(`/plans/${planId}/heartbeat-logs`);
  return response.data;
}

// === Claim API ===

export async function claimInheritance(
  planId: number,
  data: ClaimRequest
): Promise<ApiResponse<{ message: string }>> {
  const response = await apiClient.post(`/plans/${planId}/claim`, data);
  return response.data;
}

// === Status API ===

export async function getStatus(): Promise<ApiResponse<any>> {
  const response = await apiClient.get('/status');
  return response.data;
}

export async function getHealth(): Promise<ApiResponse<any>> {
  const response = await apiClient.get('/health');
  return response.data;
}
