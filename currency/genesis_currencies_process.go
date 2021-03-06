package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (op GenesisCurrencies) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := op.Fact().(GenesisCurrenciesFact)

	var newAddress base.Address
	if a, err := fact.Address(); err != nil {
		return util.IgnoreError.Wrap(err)
	} else {
		newAddress = a
	}

	var ns state.State
	if st, err := notExistsState(StateKeyAccount(newAddress), "key of genesis", getState); err != nil {
		return err
	} else {
		ns = st
	}

	gas := map[CurrencyID]state.State{}
	sts := map[CurrencyID]state.State{}
	for i := range fact.cs {
		c := fact.cs[i]

		if st, err := notExistsState(StateKeyCurrencyDesign(c.Currency()), "currency", getState); err != nil {
			return err
		} else {
			sts[c.Currency()] = st
		}

		if st, err := notExistsState(StateKeyBalance(newAddress, c.Currency()), "balance of genesis", getState); err != nil {
			return err
		} else {
			gas[c.Currency()] = NewAmountState(st, c.Currency())
		}
	}

	var states []state.State
	if ac, err := NewAccountFromKeys(fact.keys); err != nil {
		return err
	} else if st, err := SetStateAccountValue(ns, ac); err != nil {
		return util.IgnoreError.Wrap(err)
	} else {
		states = append(states, st)
	}

	for i := range fact.cs {
		c := fact.cs[i]
		am := NewAmount(c.Big(), c.Currency())
		if gst, err := SetStateBalanceValue(gas[c.Currency()], am); err != nil {
			return err
		} else if dst, err := SetStateCurrencyDesignValue(sts[c.Currency()], c); err != nil {
			return err
		} else {
			states = append(states, gst, dst)
		}
	}

	return setState(fact.Hash(), states...)
}
