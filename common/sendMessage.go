package common

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strings"
)

type SignableMessage struct {
	Version []byte
	Header  []byte
	Body    []byte
}

// SignMessage	消息签名
func SignMessage(msg string, privateKeyHex string) (string, error) {
	// 设置签名请求数据
	signableMessage := EncodeDefunct([]byte(msg))
	messageHash := hashEIP191Message(*signableMessage)

	//fmt.Println(privateKeyHex[2:], len(privateKeyHex[2:]))
	var privateKey *ecdsa.PrivateKey
	if strings.Contains(privateKeyHex, "0x") {
		privateKey, _ = crypto.HexToECDSA(privateKeyHex[2:])
	} else {
		privateKey, _ = crypto.HexToECDSA(privateKeyHex)
	}

	signature, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		fmt.Println("error3:", err)
		return "", err
	}

	return hexutil.Encode(signature), nil
}

// EncodeDefunct 使用旧的方法对消息进行编码
func EncodeDefunct(primitive []byte) *SignableMessage {
	msgLength := []byte(fmt.Sprintf("%d", len(primitive)))

	// 编码版本 E，定义在 EIP-191 中
	return &SignableMessage{
		Version: []byte{'E'},
		Header:  append([]byte("thereum Signed Message:\n"), msgLength...),
		Body:    primitive,
	}
}

func hashEIP191Message(signableMessage SignableMessage) []byte {
	version := signableMessage.Version
	if len(version) != 1 {
		log.Fatalf("The supplied message version is %x. The EIP-191 signable message standard only supports one-byte versions.", version)
	}

	joined := append([]byte("\x19"), append(version, append(signableMessage.Header, signableMessage.Body...)...)...)
	return crypto.Keccak256(joined)
}
