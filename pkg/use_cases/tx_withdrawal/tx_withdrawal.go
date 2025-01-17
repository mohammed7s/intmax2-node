package tx_withdrawal

import (
	"context"
	"errors"
	"intmax2-node/configs"
	intMaxAcc "intmax2-node/internal/accounts"
	"intmax2-node/internal/logger"
	"intmax2-node/internal/mnemonic_wallet"
	"intmax2-node/internal/open_telemetry"
	service "intmax2-node/internal/tx_withdrawal_service"
	txWithdrawal "intmax2-node/internal/use_cases/tx_withdrawal"

	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrEmptyUserPrivateKey   = errors.New("user private key is empty")
	ErrEmptyRecipientAddress = errors.New("recipient address is empty")
	ErrEmptyAmount           = errors.New("amount is empty")
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	sb  ServiceBlockchain
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	sb ServiceBlockchain,
) txWithdrawal.UseCaseTxWithdrawal {
	return &uc{
		cfg: cfg,
		log: log,
		sb:  sb,
	}
}

func (u *uc) Do(ctx context.Context, args []string, recipientAddressHex, amount, userEthPrivateKey string, resumeIncompleteWithdrawals bool) (err error) {
	const (
		hName     = "UseCase TxWithdrawal"
		senderKey = "sender"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if recipientAddressHex == "" {
		return ErrEmptyRecipientAddress
	}

	service.ResumeWithdrawalRequest(spanCtx, u.cfg, u.log, u.sb, recipientAddressHex, resumeIncompleteWithdrawals)

	if resumeIncompleteWithdrawals {
		u.log.Infof("Complete the withdrawal request")
		return nil
	}

	if userEthPrivateKey == "" {
		return ErrEmptyUserPrivateKey
	}

	wallet, err := mnemonic_wallet.New().WalletFromPrivateKeyHex(userEthPrivateKey)
	if err != nil {
		u.log.Errorf("fail to parse user private key: %v", err)
	}

	// The userPrivateKey is acceptable in either format:
	// it may include the '0x' prefix at the beginning,
	// or it can be provided without this prefix.
	userAccount, err := intMaxAcc.NewPrivateKeyFromString(wallet.IntMaxPrivateKey)
	if err != nil {
		return err
	}

	userAddress := userAccount.ToAddress()
	span.SetAttributes(
		attribute.String(senderKey, userAddress.String()),
	)

	if amount == "" {
		return ErrEmptyAmount
	}

	return service.WithdrawalTransaction(spanCtx, u.cfg, u.log, u.sb, args, amount, recipientAddressHex, userEthPrivateKey)
}
