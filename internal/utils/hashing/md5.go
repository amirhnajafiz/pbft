package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// MD5HashRequestMsg hashing gets a request and returns the digest of that message.
func MD5HashRequestMsg(request *pbft.RequestMsg) string {
	text := fmt.Sprintf(
		"%d-%s-%s-%d",
		request.GetTransaction().GetTimestamp(),
		request.GetTransaction().GetSender(),
		request.GetTransaction().GetReciever(),
		request.GetTransaction().GetAmount(),
	)

	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// MD5HashViewMsg hashing gets a view change message and returns the digest of that message.
func MD5HashViewMsg(msg *pbft.ViewChangeMsg) string {
	text := fmt.Sprintf("%d-%d", msg.GetView(), msg.GetSequenceNumber())

	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
