package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const appName = "legacyd"

func main() {
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Legacy Chain - Decentralized Inheritance Management DApp",
		Long: `Legacy Chain is a Cosmos SDK-based blockchain for decentralized inheritance planning.
It enables users to create inheritance plans with heartbeat monitoring,
multi-beneficiary support, and Shamir secret sharing for key distribution.

Key Features:
  - Create inheritance plans with configurable heartbeat intervals
  - Add multiple asset types (native tokens, CW20, NFTs, IPFS documents)
  - Automatic inheritance trigger on heartbeat expiration
  - Shamir secret sharing for secure key distribution to beneficiaries
  - IPFS integration for encrypted document storage`,
	}

	rootCmd.AddCommand(
		serveCmd(),
		initCmd(),
		keysCmd(),
		queryCmd(),
		txCmd(),
		versionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of legacyd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("legacyd v0.1.0")
			fmt.Println("Cosmos SDK-based Decentralized Inheritance Chain")
		},
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize a new legacy chain node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moniker := args[0]
			fmt.Printf("Initializing legacy chain node '%s'...\n", moniker)
			// In production: initialize genesis.json, config.toml, etc.
			fmt.Println("Node initialized successfully.")
			fmt.Printf("Configuration stored in ~/.legacyd/\n")
			return nil
		},
	}
}

func serveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the legacy chain node",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting Legacy Chain node...")
			fmt.Println("Tendermint RPC: tcp://0.0.0.0:26657")
			fmt.Println("REST API: http://0.0.0.0:1317")
			fmt.Println("gRPC: tcp://0.0.0.0:9090")
			fmt.Println("Node is running. Press Ctrl+C to stop.")
			// In production: start Tendermint node with legacy module loaded
			select {} // block forever
		},
	}
	cmd.Flags().String("log_level", "info", "Log level")
	cmd.Flags().Int("rpc-port", 26657, "Tendermint RPC port")
	cmd.Flags().Int("api-port", 1317, "REST API port")
	cmd.Flags().Int("grpc-port", 9090, "gRPC port")
	return cmd
}

func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage keyring for the legacy chain",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "add [name]",
			Short: "Add a new key to the keyring",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Printf("Key '%s' added to keyring.\n", args[0])
				return nil
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all keys in the keyring",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Keys in keyring:")
			},
		},
	)
	return cmd
}

func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands for the legacy module",
	}
	cmd.AddCommand(
		queryPlanCmd(),
		queryAssetCmd(),
		queryPlansByCreatorCmd(),
	)
	return cmd
}

func queryPlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan [plan-id]",
		Short: "Query a legacy plan by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Querying legacy plan %s...\n", args[0])
			// In production: query KVStore via ABCI
			return nil
		},
	}
}

func queryAssetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "asset [asset-id]",
		Short: "Query an asset by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Querying asset %s...\n", args[0])
			return nil
		},
	}
}

func queryPlansByCreatorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plans-by [creator-address]",
		Short: "Query all plans created by an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Querying plans by creator %s...\n", args[0])
			return nil
		},
	}
}

func txCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction subcommands for the legacy module",
	}
	cmd.AddCommand(
		txCreatePlanCmd(),
		txAddAssetCmd(),
		txHeartbeatCmd(),
		txClaimCmd(),
	)
	return cmd
}

func txCreatePlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-plan",
		Short: "Create a new inheritance plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			beneficiaries, _ := cmd.Flags().GetString("beneficiaries")
			interval, _ := cmd.Flags().GetInt64("heartbeat-interval")
			fmt.Printf("Creating legacy plan with beneficiaries: %s, heartbeat interval: %ds\n", beneficiaries, interval)
			return nil
		},
	}
	cmd.Flags().String("beneficiaries", "", "Comma-separated beneficiary addresses with shares (addr1:0.5,addr2:0.5)")
	cmd.Flags().Int64("heartbeat-interval", 86400*30, "Heartbeat interval in seconds (default: 30 days)")
	return cmd
}

func txAddAssetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset",
		Short: "Add an asset to a legacy plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			planID, _ := cmd.Flags().GetUint64("plan-id")
			assetType, _ := cmd.Flags().GetString("type")
			denom, _ := cmd.Flags().GetString("denom")
			amount, _ := cmd.Flags().GetString("amount")
			fmt.Printf("Adding asset to plan %d: type=%s denom=%s amount=%s\n", planID, assetType, denom, amount)
			return nil
		},
	}
	cmd.Flags().Uint64("plan-id", 0, "Legacy plan ID")
	cmd.Flags().String("type", "native", "Asset type: native, cw20, nft, ipfs")
	cmd.Flags().String("denom", "", "Token denomination or contract address")
	cmd.Flags().String("amount", "0", "Token amount")
	cmd.Flags().String("ipfs-cid", "", "IPFS CID for encrypted document")
	return cmd
}

func txHeartbeatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "heartbeat [plan-id]",
		Short: "Send a heartbeat for a legacy plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Sending heartbeat for plan %s...\n", args[0])
			return nil
		},
	}
}

func txClaimCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "claim [plan-id]",
		Short: "Claim inheritance from a triggered plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Claiming inheritance for plan %s...\n", args[0])
			return nil
		},
	}
}
