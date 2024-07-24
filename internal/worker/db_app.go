package worker

import (
	"context"
	"encoding/json"
	mDBApp "intmax2-node/pkg/sql_db/db_app/models"

	"github.com/holiman/uint256"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=worker_test -source=db_app.go

type SQLDriverApp interface {
	GenericCommandsApp
	Signatures
	Transactions
	TxMerkleProofs
}

type GenericCommandsApp interface {
	Exec(ctx context.Context, input interface{}, executor func(d interface{}, input interface{}) error) (err error)
}

type Signatures interface {
	CreateSignature(signature string) (*mDBApp.Signature, error)
	SignatureByID(txID string) (*mDBApp.Signature, error)
}

type Transactions interface {
	CreateTransaction(
		senderPublicKey, txHash, signatureID string,
	) (*mDBApp.Transactions, error)
	TransactionByID(txID string) (*mDBApp.Transactions, error)
}

type TxMerkleProofs interface {
	CreateTxMerkleProofs(
		senderPublicKey, txHash, txID string,
		txTreeIndex *uint256.Int,
		txMerkleProof json.RawMessage,
	) (*mDBApp.TxMerkleProofs, error)
	TxMerkleProofsByID(id string) (*mDBApp.TxMerkleProofs, error)
}