// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package deployutils

import (
	"math/big"

	"github.com/ChainSafe/sygma-relayer/chains/evm/calls/contracts/feeHandler"

	"github.com/ChainSafe/sygma-relayer/chains/evm/calls/contracts/bridge"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/common"
)

type FeeHandlerDeployResults struct {
	FeeRouter                   *feeHandler.FeeRouter
	BasicFeeHandlerAddress      common.Address
	FeeHandlerWithOracleAddress common.Address
	FeeRouterAddress            common.Address
}

func SetupFeeRouter(ethClient EVMClient, t transactor.Transactor, bridgeContract *bridge.BridgeContract) (*feeHandler.FeeRouter, error) {
	fr, err := DeployFeeRouter(ethClient, t, *bridgeContract.ContractAddress())
	if err != nil {
		return nil, err
	}
	_, err = bridgeContract.AdminChangeFeeHandler(*fr.ContractAddress(), transactor.TransactOptions{GasLimit: 2000000})
	if err != nil {
		return nil, err
	}
	return fr, nil
}

func SetupFeeHandlerWithOracle(ethClient EVMClient, t transactor.Transactor, bridgeContractAddress, feeRouterAddress, oracleAddress common.Address, feeGas uint32, feePercent uint16) (*feeHandler.FeeHandlerWithOracleContract, error) {
	fh, err := DeployFeeHandlerWithOracle(ethClient, t, bridgeContractAddress, feeRouterAddress)
	if err != nil {
		return nil, err
	}
	// Set FeeOracle address  for FeeHandlers (if required)
	_, err = fh.SetFeeOracle(oracleAddress, transactor.TransactOptions{GasLimit: 2000000})
	if err != nil {
		return nil, err
	}
	// Set fee properties (percentage, gasUsed)
	_, err = fh.SetFeeProperties(feeGas, feePercent, transactor.TransactOptions{GasLimit: 2000000})
	if err != nil {
		return nil, err
	}
	return fh, nil
}

func SetupFeeBasicHandler(ethClient EVMClient, t transactor.Transactor, bridgeContractAddress, feeRouterAddress common.Address, feeAmount *big.Int) (*feeHandler.BasicFeeHandlerContract, error) {
	fh, err := DeployBasicFeeHandler(ethClient, t, bridgeContractAddress, feeRouterAddress)
	if err != nil {
		return nil, err
	}
	return fh, err
}

func DeployFeeRouter(
	ethClient EVMClient, t transactor.Transactor, bridgeContractAddress common.Address,
) (*feeHandler.FeeRouter, error) {
	feeRouterContract := feeHandler.NewFeeRouter(ethClient, common.Address{}, t)
	_, err := feeRouterContract.DeployContract(bridgeContractAddress)
	if err != nil {
		return nil, err
	}
	return feeRouterContract, nil
}

func DeployFeeHandlerWithOracle(
	ethClient EVMClient, t transactor.Transactor, bridgeContractAddress, feeRouterAddress common.Address,
) (*feeHandler.FeeHandlerWithOracleContract, error) {
	feeHandlerContract := feeHandler.NewFeeHandlerWithOracleContract(ethClient, common.Address{}, t)
	_, err := feeHandlerContract.DeployContract(bridgeContractAddress, feeRouterAddress)
	if err != nil {
		return nil, err
	}

	return feeHandlerContract, nil
}

func DeployBasicFeeHandler(
	ethClient EVMClient, t transactor.Transactor, bridgeContractAddress, feeRouterAddress common.Address,
) (*feeHandler.BasicFeeHandlerContract, error) {
	feeHandlerContract := feeHandler.NewBasicFeeHandlerContract(ethClient, common.Address{}, t)
	_, err := feeHandlerContract.DeployContract(bridgeContractAddress, feeRouterAddress)
	if err != nil {
		return nil, err
	}

	return feeHandlerContract, nil
}
