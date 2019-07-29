package atom

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	app1Name      = "App1"
	bankCodespace = "BANK"
)

//------------------------------------------------------------------
// Msg

// MsgSend implements sdk.Msg
var _ sdk.Msg = MsgSend{}

// MsgSend to send coins from Input to Output
type MsgSend struct {
	From   sdk.AccAddress `json:"from_address"`
	To     sdk.AccAddress `json:"to_address"`
	Amount sdk.Coins      `json:"amount"`
}

// NewMsgSend
func NewMsgSend(from, to sdk.AccAddress, amt sdk.Coins) MsgSend {
	return MsgSend{from, to, amt}
}

// Implements Msg.
// nolint
func (msg MsgSend) Route() string { return "send" }
func (msg MsgSend) Type() string  { return "send" }

// Implements Msg. Ensure the addresses are good and the
// amount is positive.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if len(msg.From) == 0 {
		return sdk.ErrInvalidAddress("From address is empty")
	}
	if len(msg.To) == 0 {
		return sdk.ErrInvalidAddress("To address is empty")
	}
//	if !msg.Amount.IsPositive() {
//		return sdk.ErrInvalidCoins("Amount is not positive")
//	}
	return nil
}

// Implements Msg. JSON encode the message.
func (msg MsgSend) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// Implements Msg. Return the signer.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// Returns the sdk.Tags for the message
func (msg MsgSend) Tags() sdk.Tags {
	//return sdk.NewTags("sender", []byte(msg.From.String())).
	//	AppendTag("receiver", []byte(msg.To.String()))
	return sdk.NewTags("sender", []byte(msg.From.String())).
		AppendTag("receiver", msg.To.String())
}

//------------------------------------------------------------------
// Handler for the message

// Handle MsgSend.
// NOTE: msg.From, msg.To, and msg.Amount were already validated
// in ValidateBasic().
func handleMsgSend(key *sdk.KVStoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		sendMsg, ok := msg.(MsgSend)
		if !ok {
			// Create custom error message and return result
			// Note: Using unreserved error codespace
			return sdk.NewError(bankCodespace, 1, "MsgSend is malformed").Result()
		}

		// Load the store.
		store := ctx.KVStore(key)

		// Debit from the sender.
		if res := handleFrom(store, sendMsg.From, sendMsg.Amount); !res.IsOK() {
			return res
		}

		// Credit the receiver.
		if res := handleTo(store, sendMsg.To, sendMsg.Amount); !res.IsOK() {
			return res
		}

		// Return a success (Code 0).
		// Add list of key-value pair descriptors ("tags").
		return sdk.Result{
			Tags: sendMsg.Tags(),
		}
	}
}

// Convenience Handlers
func handleFrom(store sdk.KVStore, from sdk.AccAddress, amt sdk.Coins) sdk.Result {
	// Get sender account from the store.
	accBytes := store.Get(from)
	if accBytes == nil {
		// Account was not added to store. Return the result of the error.
		return sdk.NewError(bankCodespace, 101, "Account not added to store").Result()
	}

	// Unmarshal the JSON account bytes.
	var acc appAccount
	err := json.Unmarshal(accBytes, &acc)
	if err != nil {
		// InternalError
		return sdk.ErrInternal("Error when deserializing account").Result()
	}

	// Deduct msg amount from sender account.
	senderCoins := acc.Coins.Sub(amt)

	// If any coin has negative amount, return insufficient coins error.
	if senderCoins.IsAnyNegative() {
		return sdk.ErrInsufficientCoins("Insufficient coins in account").Result()
	}

	// Set acc coins to new amount.
	acc.Coins = senderCoins

	// Encode sender account.
	accBytes, err = json.Marshal(acc)
	if err != nil {
		return sdk.ErrInternal("Account encoding error").Result()
	}

	// Update store with updated sender account
	store.Set(from, accBytes)
	return sdk.Result{}
}

func handleTo(store sdk.KVStore, to sdk.AccAddress, amt sdk.Coins) sdk.Result {
	// Add msg amount to receiver account
	accBytes := store.Get(to)
	var acc appAccount
	if accBytes == nil {
		// Receiver account does not already exist, create a new one.
		acc = appAccount{}
	} else {
		// Receiver account already exists. Retrieve and decode it.
		err := json.Unmarshal(accBytes, &acc)
		if err != nil {
			return sdk.ErrInternal("Account decoding error").Result()
		}
	}

	// Add amount to receiver's old coins
	receiverCoins := acc.Coins.Add(amt)

	// Update receiver account
	acc.Coins = receiverCoins

	// Encode receiver account
	accBytes, err := json.Marshal(acc)
	if err != nil {
		return sdk.ErrInternal("Account encoding error").Result()
	}

	// Update store with updated receiver account
	store.Set(to, accBytes)
	return sdk.Result{}
}

// Simple account struct
type appAccount struct {
	Coins sdk.Coins `json:"coins"`
}

//------------------------------------------------------------------
// Tx

// Simple tx to wrap the Msg.
type app1Tx struct {
	MsgSend
}

// This tx only has one Msg.
func (tx app1Tx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.MsgSend}
}

// JSON decode MsgSend.
func tx1Decoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx app1Tx
	err := json.Unmarshal(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode(err.Error())
	}
	return tx, nil
}
