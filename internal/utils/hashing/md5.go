package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// MD5 hashing gets a request and returns the digest of that message.
func MD5(request *pbft.RequestMsg) string {
	text := fmt.Sprintf(
		"%d-%d-%s-%s-%d",
		request.GetSequenceNumber(),
		request.GetTransaction().GetTimestamp(),
		request.GetTransaction().GetSender(),
		request.GetTransaction().GetReciever(),
		request.GetTransaction().GetAmount(),
	)

	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
