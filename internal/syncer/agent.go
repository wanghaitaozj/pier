package syncer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/pier/internal/repo"
)

var ErrIBTPNotFound = fmt.Errorf("receipt from bitxhub failed")

func (syncer *WrapperSyncer) GetAssetExchangeSigns(id string) ([]byte, error) {
	resp, err := syncer.client.GetMultiSigns(id, pb.GetMultiSignsRequest_ASSET_EXCHANGE)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Sign == nil {
		return nil, fmt.Errorf("get empty signatures for asset exchange id %s", id)
	}

	var signs []byte
	for _, sign := range resp.Sign {
		signs = append(signs, sign...)
	}

	return signs, nil
}

func (syncer *WrapperSyncer) GetIBTPSigns(ibtp *pb.IBTP) ([]byte, error) {
	hash := ibtp.Hash()
	resp, err := syncer.client.GetMultiSigns(hash.String(), pb.GetMultiSignsRequest_IBTP)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Sign == nil {
		return nil, fmt.Errorf("get empty signatures for ibtp %s", ibtp.ID())
	}
	signs, err := resp.Marshal()
	if err != nil {
		return nil, err
	}

	return signs, nil
}

func (syncer *WrapperSyncer) GetAppchains() ([]*rpcx.Appchain, error) {
	tx, err := syncer.client.GenerateContractTx(pb.TransactionData_BVM, constant.AppchainMgrContractAddr.Address(), "Appchains")
	if err != nil {
		return nil, err
	}
	var receipt *pb.Receipt
	if err := retryFunc(func(attempt uint) error {
		receipt, err = syncer.client.SendView(tx)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	ret := make([]*rpcx.Appchain, 0)
	if receipt.Ret == nil {
		return ret, nil
	}
	if err := json.Unmarshal(receipt.Ret, &ret); err != nil {
		return nil, err
	}
	appchains := make([]*rpcx.Appchain, 0)
	for _, appchain := range ret {
		if appchain.ChainType != repo.BitxhubType {
			appchains = append(appchains, appchain)
		}
	}
	return appchains, nil
}

func (syncer *WrapperSyncer) GetInterchainById(from string) *pb.Interchain {
	ic := &pb.Interchain{}
	tx, err := syncer.client.GenerateContractTx(pb.TransactionData_BVM, constant.InterchainContractAddr.Address(), "GetInterchain", rpcx.String(from))
	if err != nil {
		return ic
	}
	receipt, err := syncer.client.SendView(tx)
	if err != nil {
		return ic
	}
	var interchain pb.Interchain
	if err := interchain.Unmarshal(receipt.Ret); err != nil {
		return ic
	}
	return &interchain
}

func (syncer *WrapperSyncer) QueryInterchainMeta() map[string]uint64 {
	interchainCounter := map[string]uint64{}
	if err := retryFunc(func(attempt uint) error {
		receipt, err := syncer.client.InvokeBVMContract(constant.InterchainContractAddr.Address(), "Interchain", nil)
		if err != nil {
			return err
		}
		if !receipt.IsSuccess() {
			return fmt.Errorf("receipt: %s", receipt.Ret)
		}
		ret := &pb.Interchain{}
		if err := ret.Unmarshal(receipt.Ret); err != nil {
			return fmt.Errorf("unmarshal interchain meta from bitxhub: %w", err)
		}
		return nil
	}); err != nil {
		logger.Panicf("query interchain meta: %s", err.Error())
	}

	return interchainCounter
}

func (syncer *WrapperSyncer) QueryIBTP(ibtpID string) (*pb.IBTP, error) {
	receipt, err := syncer.client.InvokeContract(pb.TransactionData_BVM, constant.InterchainContractAddr.Address(),
		"GetIBTPByID", nil, rpcx.String(ibtpID))
	if err != nil {
		return nil, err
	}

	if !receipt.IsSuccess() {
		return nil, fmt.Errorf("%w: %s", ErrIBTPNotFound, string(receipt.Ret))
	}
	hash := types.NewHash(receipt.Ret)

	response, err := syncer.client.GetTransaction(hash.String())
	if err != nil {
		return nil, err
	}
	ibtp := response.Tx.GetIBTP()
	if ibtp == nil {
		return nil, fmt.Errorf("empty ibtp from bitxhub")
	}
	return ibtp, nil
}

func (syncer *WrapperSyncer) ListenIBTP() <-chan *pb.IBTP {
	return syncer.ibtpC
}

func (syncer *WrapperSyncer) SendIBTP(ibtp *pb.IBTP) error {
	proof := ibtp.GetProof()
	proofHash := sha256.Sum256(proof)
	ibtp.Proof = proofHash[:]

	tx, err := syncer.client.GenerateIBTPTx(ibtp)
	if err != nil {
		return fmt.Errorf("generate ibtp tx error:%v", err)
	}
	tx.Extra = proof
	retryFunc(func(attempt uint) error {
		receipt, err := syncer.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
			From:      fmt.Sprintf("%s-%s-%d", ibtp.From, ibtp.To, ibtp.Category()),
			IBTPNonce: ibtp.Index,
		})
		if err != nil {
			return err
		}
		if receipt.IsSuccess() {
			return nil
		}
		// query if this ibtp is on chain
		_, err = syncer.QueryIBTP(ibtp.ID())
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func retryFunc(handle func(uint) error) error {
	return retry.Retry(func(attempt uint) error {
		if err := handle(attempt); err != nil {
			logger.Errorf("retry failed for reason: %s", err.Error())
		}
		return nil
	}, strategy.Wait(500*time.Millisecond))
}