// === Blockchain Types ===

export interface LegacyPlan {
  plan_id: number;
  creator: string;
  beneficiaries: Beneficiary[];
  heartbeat_interval: number;
  last_heartbeat: string;
  trigger_time: string;
  status: PlanStatus;
  created_at: string;
  encrypted_key?: string;
  ipfs_cid?: string;
}

export interface Beneficiary {
  address: string;
  share: number;
  key_share?: string;
}

export type PlanStatus = 'active' | 'triggered' | 'claimed';

export interface Asset {
  asset_id: number;
  plan_id: number;
  owner: string;
  asset_type: AssetType;
  denom: string;
  amount: string;
  ipfs_cid?: string;
  metadata?: string;
  encrypted_data?: string;
}

export type AssetType = 'native' | 'cw20' | 'nft' | 'ipfs';

// === API Types ===

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface CreatePlanRequest {
  beneficiaries: { address: string; share: number }[];
  heartbeat_interval: number;
}

export interface AddAssetRequest {
  plan_id: number;
  asset_type: AssetType;
  denom: string;
  amount: string;
  ipfs_cid?: string;
  metadata?: string;
}

export interface HeartbeatRequest {
  sender_address: string;
  tx_hash?: string;
}

export interface ClaimRequest {
  beneficiary_address: string;
  key_shares?: number[];
}

// === Wallet Types ===

export interface WalletState {
  connected: boolean;
  address: string;
  name: string;
  balance: string;
}

// === Keplr Types ===

export interface KeplrWindow {
  keplr?: {
    enable(chainId: string): Promise<void>;
    getKey(chainId: string): Promise<{
      name: string;
      algo: string;
      pubKey: Uint8Array;
      address: string;
      bech32Address: string;
    }>;
    getOfflineSigner(chainId: string): any;
    experimentalSuggestChain(chainInfo: any): Promise<void>;
  };
}

declare global {
  interface Window extends KeplrWindow {}
}
