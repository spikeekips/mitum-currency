package digest

import (
	"encoding/json"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/localtime"
	"github.com/spikeekips/mitum/util/valuehash"
)

type OperationValueJSONPacker struct {
	jsonenc.HintedHead
	HS valuehash.Hash      `json:"hash"`
	OP operation.Operation `json:"operation"`
	HT base.Height         `json:"height"`
	CF localtime.Time      `json:"confirmed_at"`
	IN bool                `json:"in_state"`
	ID uint64              `json:"index"`
}

func (va OperationValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(OperationValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(va.Hint()),
		HS:         va.op.Fact().Hash(),
		OP:         va.op,
		HT:         va.height,
		CF:         localtime.NewTime(va.confirmedAt),
		IN:         va.inStates,
		ID:         va.index,
	})
}

type OperationValueJSONUnpacker struct {
	OP json.RawMessage `json:"operation"`
	HT base.Height     `json:"height"`
	CF localtime.Time  `json:"confirmed_at"`
	IN bool            `json:"in_state"`
	ID uint64          `json:"index"`
}

func (va *OperationValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva OperationValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	if op, err := operation.DecodeOperation(enc, uva.OP); err != nil {
		return err
	} else {
		va.op = op
	}

	va.height = uva.HT
	va.confirmedAt = uva.CF.Time
	va.inStates = uva.IN
	va.index = uva.ID

	return nil
}
