package txs

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	"github.com/Quantiex-Hub/cmd/ethrelayer/types"
)

// LoadPrivateKey loads the validator's private key from environment variables
func LoadPrivateKey() (key *ecdsa.PrivateKey, err error) {
	// Load config file containing environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Private key for validator's Ethereum address must be set as an environment variable
	rawPrivateKey := os.Getenv("ETHEREUM_PRIVATE_KEY")
	if strings.TrimSpace(rawPrivateKey) == "" {
		log.Fatal("Error loading ETHEREUM_PRIVATE_KEY from .env file")
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey, nil
}

// LoadSender uses the validator's private key to load the validator's address
func LoadSender() (address common.Address, err error) {
	key, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress, nil
}

// GenerateERC20ClaimMessage Generates a hashed message containing a ProphecyClaim event's data
func GenerateERC20ClaimMessage(event types.ProphecyClaimERC20Event) []byte {
	prophecyID := Int256(event.ProphecyID)
	chainName := String(event.ChainName)
	claimType := Uint8(event.ClaimType)
	sender := String(event.BinanceSender.Hex())
	recipient := String(event.EthereumReceiver.Hex())
	validator := String(event.ValidatorAddress.Hex())
	token := String(event.TokenAddress.Hex())
	symbol := String(event.Symbol)
	amount := Int256(event.Amount)
	txHash := String(event.TxHash)

	// Generate claim message using ProphecyClaim data
	return SoliditySHA3(prophecyID, chainName, claimType, sender, recipient, validator, token, symbol, amount, txHash)
}

// GenerateERC721ClaimMessage Generates a hashed message containing a ProphecyClaim event's data
func GenerateERC721ClaimMessage(event types.ProphecyClaimERC721Event) []byte {
	prophecyID := Int256(event.ProphecyID)
	chainName := String(event.ChainName)
	claimType := Uint8(event.ClaimType)
	sender := String(event.BinanceSender.Hex())
	recipient := String(event.EthereumReceiver.Hex())
	validator := String(event.ValidatorAddress.Hex())
	token := String(event.TokenAddress.Hex())
	symbol := String(event.Symbol)
	tokenId := Int256(event.TokenId)
	tokenURI := String(event.TokenURI)
	txHash := String(event.TxHash)

	// Generate claim message using ProphecyClaim data
	return SoliditySHA3(prophecyID, chainName, claimType, sender, recipient, validator, token, symbol, tokenId, tokenURI, txHash)
}

// PrefixMsg prefixes a message for verification, mimics behavior of web3.eth.sign
func PrefixMsg(msg []byte) []byte {
	return SoliditySHA3(String("\x19Ethereum Signed Message:\n32"), msg)
}

// SignClaim Signs the prepared message with validator's private key
func SignClaim(msg []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	// Sign the message
	sig, err := crypto.Sign(msg, key)
	if err != nil {
		panic(err)
	}
	return sig, nil
}



