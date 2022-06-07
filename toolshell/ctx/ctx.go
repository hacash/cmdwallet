package ctx

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

type Context interface {
	NotLoadedYetAccountAddress(string) bool         // Check whether the account has been logged in
	IsInvalidAccountAddress(string) *fields.Address // Check whether it is a legal account name
	IsInvalidAmountString(string) *fields.Amount    // Check whether it is a legal amount quantity
	GetAllPrivateKeyBytes() map[string][]byte       // Get all private keys for filling in signatures
	SetTxToRecord([]byte, interfaces.Transaction)   // Record transactions
	GetTxFromRecord([]byte) interfaces.Transaction  // Acquire transaction
	UseTimestamp() uint64                           // Currently used timestamp

	//////////////////////////////////////////////////////////////////

	Println(...interface{})
	Print(...interface{})
}
