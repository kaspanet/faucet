package main

import (
	"encoding/hex"

	"github.com/kaspanet/kaspad/domain/consensus/utils/utxo"

	"github.com/kaspanet/faucet/config"
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/domain/consensus/utils/consensushashing"
	"github.com/kaspanet/kaspad/domain/consensus/utils/constants"
	"github.com/kaspanet/kaspad/domain/consensus/utils/subnetworks"
	"github.com/kaspanet/kaspad/domain/consensus/utils/transactionid"
	"github.com/kaspanet/kaspad/domain/consensus/utils/txscript"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util"
	"github.com/pkg/errors"
)

const (
	sendAmountKaspa       = 1
	feeSompis             = 3000
	requiredConfirmations = 10
)

func sendToAddress(address util.Address) (string, error) {
	cfg, err := config.MainConfig()
	if err != nil {
		return "", err
	}
	client, err := rpcclient.NewRPCClient(cfg.RPCServer)
	if err != nil {
		return "", err
	}
	utxos, err := fetchSpendableUTXOs(client)
	if err != nil {
		return "", err
	}

	sendAmountSompi := uint64(sendAmountKaspa * util.SompiPerKaspa)
	totalToSend := sendAmountSompi + feeSompis
	selectedUTXOs, changeSompi, err := selectUTXOs(utxos, totalToSend)
	if err != nil {
		return "", err
	}

	rpcTransaction, err := generateTransaction(selectedUTXOs, sendAmountSompi, changeSompi, address)
	if err != nil {
		return "", err
	}

	return sendTransaction(client, rpcTransaction)
}

func fetchSpendableUTXOs(client *rpcclient.RPCClient) ([]*appmessage.UTXOsByAddressesEntry, error) {
	getUTXOsByAddressesResponse, err := client.GetUTXOsByAddresses([]string{faucetAddress.EncodeAddress()})
	if err != nil {
		return nil, err
	}
	virtualSelectedParentBlueScoreResponse, err := client.GetVirtualSelectedParentBlueScore()
	if err != nil {
		return nil, err
	}
	virtualSelectedParentBlueScore := virtualSelectedParentBlueScoreResponse.BlueScore

	spendableUTXOs := make([]*appmessage.UTXOsByAddressesEntry, 0)
	for _, entry := range getUTXOsByAddressesResponse.Entries {
		if !isUTXOSpendable(entry, virtualSelectedParentBlueScore) {
			continue
		}
		spendableUTXOs = append(spendableUTXOs, entry)
	}
	return spendableUTXOs, nil
}

func isUTXOSpendable(entry *appmessage.UTXOsByAddressesEntry, virtualSelectedParentBlueScore uint64) bool {
	blockDAAScore := entry.UTXOEntry.BlockDAAScore
	if !entry.UTXOEntry.IsCoinbase {
		return blockDAAScore+requiredConfirmations < virtualSelectedParentBlueScore
	}
	coinbaseMaturity := config.ActiveNetParams().BlockCoinbaseMaturity
	return blockDAAScore+coinbaseMaturity < virtualSelectedParentBlueScore
}

func selectUTXOs(utxos []*appmessage.UTXOsByAddressesEntry, totalToSpend uint64) (
	selectedUTXOs []*appmessage.UTXOsByAddressesEntry, changeSompi uint64, err error) {

	selectedUTXOs = []*appmessage.UTXOsByAddressesEntry{}
	totalValue := uint64(0)

	for _, utxo := range utxos {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalValue += utxo.UTXOEntry.Amount

		if totalValue >= totalToSpend {
			break
		}
	}

	if totalValue < totalToSpend {
		return nil, 0, errors.Errorf("Insufficient funds for send: %f required, while only %f available",
			float64(totalToSpend)/util.SompiPerKaspa, float64(totalValue)/util.SompiPerKaspa)
	}

	return selectedUTXOs, totalValue - totalToSpend, nil
}

func generateTransaction(selectedUTXOs []*appmessage.UTXOsByAddressesEntry,
	sompisToSend uint64, change uint64, toAddress util.Address) (*appmessage.RPCTransaction, error) {

	inputs := make([]*externalapi.DomainTransactionInput, len(selectedUTXOs))
	for i, selectedUTXO := range selectedUTXOs {
		outpointTransactionIDBytes, err := hex.DecodeString(selectedUTXO.Outpoint.TransactionID)
		if err != nil {
			return nil, err
		}
		outpointTransactionID, err := transactionid.FromBytes(outpointTransactionIDBytes)
		if err != nil {
			return nil, err
		}
		outpoint := externalapi.DomainOutpoint{
			TransactionID: *outpointTransactionID,
			Index:         selectedUTXO.Outpoint.Index,
		}

		utxoEntry, err := utxoEntryToDomain(selectedUTXO)
		if err != nil {
			return nil, err
		}

		inputs[i] = &externalapi.DomainTransactionInput{
			PreviousOutpoint: outpoint,
			SignatureScript:  nil,
			Sequence:         0,
			UTXOEntry:        utxoEntry,
		}
	}

	toScript, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		return nil, err
	}
	mainOutput := &externalapi.DomainTransactionOutput{
		Value:           sompisToSend,
		ScriptPublicKey: toScript,
	}
	fromScript, err := txscript.PayToAddrScript(faucetAddress)
	if err != nil {
		return nil, err
	}
	changeOutput := &externalapi.DomainTransactionOutput{
		Value:           change,
		ScriptPublicKey: fromScript,
	}
	outputs := []*externalapi.DomainTransactionOutput{mainOutput, changeOutput}

	domainTransaction := &externalapi.DomainTransaction{
		Version:      constants.MaxTransactionVersion,
		Inputs:       inputs,
		Outputs:      outputs,
		LockTime:     0,
		SubnetworkID: subnetworks.SubnetworkIDNative,
		Gas:          0,
		Payload:      nil,
	}

	sighashReusedValues := &consensushashing.SighashReusedValues{}
	for i, input := range domainTransaction.Inputs {
		signatureScript, err := txscript.SignatureScript(
			domainTransaction, i, consensushashing.SigHashAll, faucetPrivateKey, sighashReusedValues)
		if err != nil {
			return nil, err
		}
		input.SignatureScript = signatureScript
	}

	rpcTransaction := appmessage.DomainTransactionToRPCTransaction(domainTransaction)
	return rpcTransaction, nil
}

func utxoEntryToDomain(selectedUTXO *appmessage.UTXOsByAddressesEntry) (externalapi.UTXOEntry, error) {
	scriptPublicKey, err := hex.DecodeString(selectedUTXO.UTXOEntry.ScriptPublicKey.Script)
	if err != nil {
		return nil, err
	}
	return utxo.NewUTXOEntry(
		selectedUTXO.UTXOEntry.Amount,
		&externalapi.ScriptPublicKey{
			Script:  scriptPublicKey,
			Version: selectedUTXO.UTXOEntry.ScriptPublicKey.Version,
		},
		selectedUTXO.UTXOEntry.IsCoinbase,
		selectedUTXO.UTXOEntry.BlockDAAScore), nil
}

func sendTransaction(client *rpcclient.RPCClient, rpcTransaction *appmessage.RPCTransaction) (string, error) {
	submitTransactionResponse, err := client.SubmitTransaction(rpcTransaction, false)
	if err != nil {
		return "", errors.Wrapf(err, "error submitting transaction")
	}
	return submitTransactionResponse.TransactionID, nil
}
