import { SigningStargateClient, StargateClient } from '@cosmjs/stargate';
import { OfflineSigner } from '@cosmjs/proto-signing';

// Chain configuration for Legacy Chain
export const CHAIN_CONFIG = {
  chainId: 'legacy-testnet-1',
  chainName: 'Legacy Chain',
  rpc: 'http://localhost:26657',
  rest: 'http://localhost:1317',
  bip44: {
    coinType: 118,
  },
  bech32Config: {
    bech32PrefixAccAddr: 'legacy',
    bech32PrefixAccPub: 'legacy' + 'pub',
    bech32PrefixValAddr: 'legacy' + 'valoper',
    bech32PrefixValPub: 'legacy' + 'valoperpub',
    bech32PrefixConsAddr: 'legacy' + 'valcons',
    bech32PrefixConsPub: 'legacy' + 'valconspub',
  },
  currencies: [
    {
      coinDenom: 'LEG',
      coinMinimalDenom: 'uleg',
      coinDecimals: 6,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: 'LEG',
      coinMinimalDenom: 'uleg',
      coinDecimals: 6,
      gasPriceStep: {
        low: 0.01,
        average: 0.025,
        high: 0.04,
      },
    },
  ],
  stakeCurrency: {
    coinDenom: 'LEG',
    coinMinimalDenom: 'uleg',
    coinDecimals: 6,
  },
};

// Message types for the legacy module
export const LEGACY_MSG_TYPES = {
  CREATE_PLAN: '/legacy.MsgCreateLegacyPlan',
  ADD_ASSET: '/legacy.MsgAddAsset',
  HEARTBEAT: '/legacy.MsgHeartbeat',
  CLAIM: '/legacy.MsgClaimInheritance',
  UPDATE_BENEFICIARIES: '/legacy.MsgUpdateBeneficiaries',
};

export class BlockchainService {
  private client: StargateClient | null = null;
  private signingClient: SigningStargateClient | null = null;

  async connect(rpcEndpoint?: string): Promise<void> {
    const rpc = rpcEndpoint || CHAIN_CONFIG.rpc;
    this.client = await StargateClient.connect(rpc);
    console.log('Connected to Legacy Chain at', rpc);
  }

  async connectWithSigner(signer: OfflineSigner, rpcEndpoint?: string): Promise<void> {
    const rpc = rpcEndpoint || CHAIN_CONFIG.rpc;
    this.signingClient = await SigningStargateClient.connectWithSigner(rpc, signer);
    this.client = this.signingClient;
    console.log('Connected with signer to Legacy Chain at', rpc);
  }

  async getBalance(address: string, denom = 'uleg'): Promise<string> {
    if (!this.client) throw new Error('Client not connected');
    const balance = await this.client.getBalance(address, denom);
    return balance.amount;
  }

  async getChainId(): Promise<string> {
    if (!this.client) throw new Error('Client not connected');
    return this.client.getChainId();
  }

  async getHeight(): Promise<number> {
    if (!this.client) throw new Error('Client not connected');
    return this.client.getHeight();
  }

  // === Transaction Methods ===

  async createLegacyPlan(
    creatorAddress: string,
    beneficiaries: { address: string; share: string }[],
    heartbeatInterval: number
  ): Promise<string> {
    if (!this.signingClient) throw new Error('Signing client not connected');

    const msg = {
      typeUrl: LEGACY_MSG_TYPES.CREATE_PLAN,
      value: {
        creator: creatorAddress,
        beneficiaries: beneficiaries.map((b) => ({
          address: b.address,
          share: b.share,
        })),
        heartbeat_interval: heartbeatInterval,
      },
    };

    const result = await this.signingClient.signAndBroadcast(
      creatorAddress,
      [msg],
      {
        amount: [{ denom: 'uleg', amount: '5000' }],
        gas: '200000',
      }
    );

    return result.transactionHash;
  }

  async addAsset(
    ownerAddress: string,
    planId: number,
    assetType: string,
    denom: string,
    amount: string,
    ipfsCid?: string
  ): Promise<string> {
    if (!this.signingClient) throw new Error('Signing client not connected');

    const msg = {
      typeUrl: LEGACY_MSG_TYPES.ADD_ASSET,
      value: {
        owner: ownerAddress,
        plan_id: planId,
        asset_type: assetType,
        denom: denom,
        amount: amount,
        ipfs_cid: ipfsCid || '',
      },
    };

    const result = await this.signingClient.signAndBroadcast(
      ownerAddress,
      [msg],
      {
        amount: [{ denom: 'uleg', amount: '5000' }],
        gas: '200000',
      }
    );

    return result.transactionHash;
  }

  async sendHeartbeat(creatorAddress: string, planId: number): Promise<string> {
    if (!this.signingClient) throw new Error('Signing client not connected');

    const msg = {
      typeUrl: LEGACY_MSG_TYPES.HEARTBEAT,
      value: {
        creator: creatorAddress,
        plan_id: planId,
      },
    };

    const result = await this.signingClient.signAndBroadcast(
      creatorAddress,
      [msg],
      {
        amount: [{ denom: 'uleg', amount: '2000' }],
        gas: '100000',
      }
    );

    return result.transactionHash;
  }

  async claimInheritance(beneficiaryAddress: string, planId: number): Promise<string> {
    if (!this.signingClient) throw new Error('Signing client not connected');

    const msg = {
      typeUrl: LEGACY_MSG_TYPES.CLAIM,
      value: {
        beneficiary: beneficiaryAddress,
        plan_id: planId,
      },
    };

    const result = await this.signingClient.signAndBroadcast(
      beneficiaryAddress,
      [msg],
      {
        amount: [{ denom: 'uleg', amount: '5000' }],
        gas: '200000',
      }
    );

    return result.transactionHash;
  }

  disconnect(): void {
    this.client?.disconnect();
    this.client = null;
    this.signingClient = null;
  }
}

export const blockchainService = new BlockchainService();
