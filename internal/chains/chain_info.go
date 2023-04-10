package chains

// Generated from gojsonstruct of https://raw.githubusercontent.com/cosmos/chain-registry/master/juno/chain.json
type ChainInfo struct {
	Schema string `json:"$schema"`
	Apis   struct {
		Grpc []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"grpc"`
		Rest []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"rest"`
		Rpc []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"rpc"`
	} `json:"apis"`
	Bech32Prefix string `json:"bech32_prefix"`
	ChainID      string `json:"chain_id"`
	ChainName    string `json:"chain_name"`
	Codebase     struct {
		CompatibleVersions []string `json:"compatible_versions"`
		Consensus          struct {
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"consensus"`
		CosmosSdkVersion string `json:"cosmos_sdk_version"`
		CosmwasmEnabled  bool   `json:"cosmwasm_enabled"`
		CosmwasmVersion  string `json:"cosmwasm_version"`
		Genesis          struct {
			GenesisURL string `json:"genesis_url"`
		} `json:"genesis"`
		GitRepo            string `json:"git_repo"`
		RecommendedVersion string `json:"recommended_version"`
		Versions           []struct {
			CompatibleVersions []string `json:"compatible_versions"`
			Consensus          struct {
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"consensus"`
			CosmosSdkVersion   string `json:"cosmos_sdk_version"`
			CosmwasmEnabled    bool   `json:"cosmwasm_enabled"`
			CosmwasmVersion    string `json:"cosmwasm_version"`
			Name               string `json:"name"`
			RecommendedVersion string `json:"recommended_version"`
		} `json:"versions"`
	} `json:"codebase"`
	DaemonName string `json:"daemon_name"`
	Explorers  []struct {
		AccountPage string `json:"account_page,omitempty"`
		Kind        string `json:"kind"`
		TxPage      string `json:"tx_page"`
		URL         string `json:"url"`
	} `json:"explorers"`
	Fees struct {
		FeeTokens []struct {
			AverageGasPrice  float64 `json:"average_gas_price"`
			Denom            string  `json:"denom"`
			FixedMinGasPrice float64 `json:"fixed_min_gas_price"`
			HighGasPrice     float64 `json:"high_gas_price"`
			LowGasPrice      float64 `json:"low_gas_price"`
		} `json:"fee_tokens"`
	} `json:"fees"`
	KeyAlgos    []string `json:"key_algos"`
	NetworkType string   `json:"network_type"`
	NodeHome    string   `json:"node_home"`
	Peers       struct {
		PersistentPeers []struct {
			Address  string `json:"address"`
			ID       string `json:"id"`
			Provider string `json:"provider,omitempty"`
		} `json:"persistent_peers"`
		Seeds []struct {
			Address  string `json:"address"`
			ID       string `json:"id"`
			Provider string `json:"provider,omitempty"`
		} `json:"seeds"`
	} `json:"peers"`
	PrettyName string `json:"pretty_name"`
	Slip44     int    `json:"slip44"`
	Staking    struct {
		StakingTokens []struct {
			Denom string `json:"denom"`
		} `json:"staking_tokens"`
	} `json:"staking"`
	Status  string `json:"status"`
	Website string `json:"website"`
}
