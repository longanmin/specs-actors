package states

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/specs-actors/actors/states"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/specs-actors/v2/actors/builtin"
	init_ "github.com/filecoin-project/specs-actors/v2/actors/builtin/init"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/verifreg"
)

func CheckStateInvariants(tree states.Tree, expectedBalanceTotal abi.TokenAmount) (*builtin.MessageAccumulator, error) {
	msgs := &builtin.MessageAccumulator{}
	totalFIl := big.Zero()
	var initSummary *init_.StateSummary
	var verifregSummary *verifreg.StateSummary

	if err := tree.ForEach(func(key addr.Address, actor *states.Actor) error {
		if key.Protocol() != addr.ID {
			msgs.Addf("unexpected address protocol in state tree root: %v", key)
		}
		totalFIl = big.Add(totalFIl, actor.Balance)

		switch actor.Code {
		case builtin.SystemActorCodeID:

		case builtin.InitActorCodeID:
			var st init_.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			if summary, msgs, err := init_.CheckStateInvariants(&st, tree.Store); err != nil {
				return err
			} else {
				msgs.AddAll(msgs)
				initSummary = summary
			}
		case builtin.CronActorCodeID:

		case builtin.AccountActorCodeID:

		case builtin.StoragePowerActorCodeID:

		case builtin.StorageMinerActorCodeID:

		case builtin.StorageMarketActorCodeID:

		case builtin.PaymentChannelActorCodeID:

		case builtin.MultisigActorCodeID:

		case builtin.RewardActorCodeID:

		case builtin.VerifiedRegistryActorCodeID:
			var st verifreg.State
			if err := tree.Store.Get(tree.Store.Context(), actor.Head, &st); err != nil {
				return err
			}
			if summary, msgs, err := verifreg.CheckStateInvariants(&st, tree.Store); err != nil {
				return err
			} else {
				msgs.AddAll(msgs)
				verifregSummary = summary
			}
		default:
			return xerrors.Errorf("unexpected actor code CID %v for address %v", actor.Code, key)

		}
		return nil
	}); err != nil {
		return nil, err
	}

	//
	// Perform cross-actor checks from state summaries here.
	//
	_ = initSummary
	_ = verifregSummary

	if !totalFIl.Equals(expectedBalanceTotal) {
		msgs.Addf("total token balance is %v, expected %v", totalFIl, expectedBalanceTotal)
	}

	return msgs, nil
}
