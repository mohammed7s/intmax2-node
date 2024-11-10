package blockchain

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	errorsB "intmax2-node/internal/blockchain/errors"
	"intmax2-node/internal/open_telemetry"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prodadidb/go-validation"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ErrScrollChainIDInvalid = fmt.Errorf(errorsB.ErrScrollChainIDInvalidStr, ScrollMainNetChainID, ScrollSepoliaChainID)

var ErrEthereumChainIDInvalid = fmt.Errorf(
	errorsB.ErrEthereumChainIDInvalidStr, EthereumMainNetChainID, EthereumSepoliaChainID,
)

type ChainIDType string

const (
	EthereumMainNetChainID ChainIDType = "5115"
	EthereumSepoliaChainID ChainIDType = "5115"

	ScrollMainNetChainID ChainIDType = "5115"
	ScrollSepoliaChainID ChainIDType = "5115"
)

type ChainLinkEvmJSONRPC string

const (
	EthereumMainNetChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.testnet.citrea.xyz"
	EthereumSepoliaChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.testnet.citrea.xyz"

	ScrollMainNetChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.testnet.citrea.xyz"
	ScrollSepoliaChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.testnet.citrea.xyz"
)

type ChainLinkExplorer string

const (
	EthereumMainNetChainLinkExplorer ChainLinkExplorer = "https://explorer.testnet.citrea.xyz"
	EthereumSepoliaChainLinkExplorer ChainLinkExplorer = "https://explorer.testnet.citrea.xyz"

	ScrollMainNetChainLinkExplorer ChainLinkExplorer = "https://explorer.testnet.citrea.xyz"
	ScrollSepoliaChainLinkExplorer ChainLinkExplorer = "https://explorer.testnet.citrea.xyz"
)

func (sb *serviceBlockchain) scrollNetworkChainIDValidator() error {
	return validation.Validate(sb.cfg.Blockchain.ScrollNetworkChainID,
		validation.Required,
		validation.In(
			string(ScrollMainNetChainID), string(ScrollSepoliaChainID),
		),
	)
}

func (sb *serviceBlockchain) SetupScrollNetworkChainID(ctx context.Context) error {
	const (
		hName = "ServiceBlockchain func:SetupScrollNetworkChainID"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	err := sb.scrollNetworkChainIDValidator()
	if err != nil {
		const (
			enterMSG = "Enter the Scroll network chain-ID:"
			crlf     = '\n'
		)
		fmt.Printf(enterMSG)
		var chainID string
		chainID, err = bufio.NewReader(os.Stdin).ReadString(crlf)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return errors.Join(errorsB.ErrStdinProcessingFail, err)
		}
		sb.cfg.Blockchain.ScrollNetworkChainID = strings.TrimSpace(chainID)
	}

	err = sb.scrollNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return errors.Join(ErrScrollChainIDInvalid, err)
	}

	return nil
}

func (sb *serviceBlockchain) ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error) {
	const (
		hName                   = "ServiceBlockchain func:ScrollNetworkChainLinkEvmJSONRPC"
		scrollNetworkChainIDKey = "scroll_network_chain_id"
		emptyKey                = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(scrollNetworkChainIDKey, sb.cfg.Blockchain.ScrollNetworkChainID),
		))
	defer span.End()

	err := sb.scrollNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrScrollChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.ScrollNetworkChainID, string(ScrollMainNetChainID)) {
		return string(ScrollMainNetChainLinkEvmJSONRPC), nil
	}

	return string(ScrollSepoliaChainLinkEvmJSONRPC), nil
}

func (sb *serviceBlockchain) ScrollNetworkChainLinkExplorer(ctx context.Context) (string, error) {
	const (
		hName                   = "ServiceBlockchain func:ScrollNetworkChainLinkExplorer"
		scrollNetworkChainIDKey = "scroll_network_chain_id"
		emptyKey                = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(scrollNetworkChainIDKey, sb.cfg.Blockchain.ScrollNetworkChainID),
		))
	defer span.End()

	err := sb.scrollNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrScrollChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.ScrollNetworkChainID, string(ScrollMainNetChainID)) {
		return string(ScrollMainNetChainLinkExplorer), nil
	}

	return string(ScrollSepoliaChainLinkExplorer), nil
}

func (sb *serviceBlockchain) ethereumNetworkChainIDValidator() error {
	return validation.Validate(sb.cfg.Blockchain.EthereumNetworkChainID,
		validation.Required,
		validation.In(
			string(EthereumMainNetChainID), string(EthereumSepoliaChainID),
		),
	)
}

func (sb *serviceBlockchain) SetupEthereumNetworkChainID(ctx context.Context) error {
	const (
		hName    = "ServiceBlockchain func:SetupEthereumNetworkChainID"
		emptyKey = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	sb.cfg.Blockchain.EthereumNetworkRpcUrl = strings.TrimSpace(sb.cfg.Blockchain.EthereumNetworkRpcUrl)
	if sb.cfg.Blockchain.EthereumNetworkRpcUrl != emptyKey {
		client, err := ethclient.DialContext(spanCtx, sb.cfg.Blockchain.EthereumNetworkRpcUrl)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return errors.Join(errorsB.ErrEthClientDialFail)
		}

		var chainID *big.Int
		chainID, err = client.ChainID(spanCtx)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return errors.Join(errorsB.ErrChainIDWithEthClientFail)
		}

		sb.cfg.Blockchain.EthereumNetworkChainID = chainID.String()
	}

	err := sb.ethereumNetworkChainIDValidator()
	if err != nil {
		const (
			enterMSG = "Enter the Ethereum network chain-ID:"
			crlf     = '\n'
		)
		fmt.Printf(enterMSG)
		var chainID string
		chainID, err = bufio.NewReader(os.Stdin).ReadString(crlf)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return errors.Join(errorsB.ErrStdinProcessingFail, err)
		}
		sb.cfg.Blockchain.EthereumNetworkChainID = strings.TrimSpace(chainID)
	}

	err = sb.ethereumNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return errors.Join(ErrEthereumChainIDInvalid, err)
	}

	return nil
}

func (sb *serviceBlockchain) EthereumNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error) {
	const (
		hName                     = "ServiceBlockchain func:EthereumNetworkChainLinkEvmJSONRPC"
		ethereumNetworkChainIDKey = "ethereum_network_chain_id"
		emptyKey                  = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(ethereumNetworkChainIDKey, sb.cfg.Blockchain.EthereumNetworkChainID),
		))
	defer span.End()

	sb.cfg.Blockchain.EthereumNetworkRpcUrl = strings.TrimSpace(sb.cfg.Blockchain.EthereumNetworkRpcUrl)
	if sb.cfg.Blockchain.EthereumNetworkRpcUrl != emptyKey {
		client, err := ethclient.DialContext(spanCtx, sb.cfg.Blockchain.EthereumNetworkRpcUrl)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return emptyKey, errors.Join(errorsB.ErrEthClientDialFail)
		}

		var chainID *big.Int
		chainID, err = client.ChainID(spanCtx)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return emptyKey, errors.Join(errorsB.ErrChainIDWithEthClientFail)
		}

		sb.cfg.Blockchain.EthereumNetworkChainID = chainID.String()

		return sb.cfg.Blockchain.EthereumNetworkRpcUrl, nil
	}

	err := sb.ethereumNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrEthereumChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.EthereumNetworkChainID, string(EthereumMainNetChainID)) {
		return string(EthereumMainNetChainLinkEvmJSONRPC), nil
	}

	return string(EthereumSepoliaChainLinkEvmJSONRPC), nil
}

func (sb *serviceBlockchain) EthereumNetworkChainLinkExplorer(ctx context.Context) (string, error) {
	const (
		hName                     = "ServiceBlockchain func:EthereumNetworkChainLinkExplorer"
		ethereumNetworkChainIDKey = "ethereum_network_chain_id"
		emptyKey                  = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(ethereumNetworkChainIDKey, sb.cfg.Blockchain.EthereumNetworkChainID),
		))
	defer span.End()

	err := sb.ethereumNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrEthereumChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.ScrollNetworkChainID, string(EthereumMainNetChainID)) {
		return string(EthereumMainNetChainLinkExplorer), nil
	}

	return string(EthereumSepoliaChainLinkExplorer), nil
}
