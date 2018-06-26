package client

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/loomnetwork/go-loom/common/evmcompat"
)

type LoomTx struct {
	Slot         uint64
	Denomination uint32 //TODO should be bigint
	Owner        common.Address
	PrevBlock    *big.Int
	Signature    []byte
}

func (l *LoomTx) Sig() []byte {
	return l.Signature
}

func (l *LoomTx) NewOwner() common.Address {
	return l.Owner
}

func (l *LoomTx) Sign(key *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(l.Hash(), key)
	if err != nil {
		return nil, err
	}
	if len(sig) != 65 {
		return nil, errors.New(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}

	r := make([]byte, 1, 65)
	r = append(r, sig[:64]...)
	return append(r, sig[64]+27), nil
}

func (l *LoomTx) RlpEncode() ([]byte, error) {
	logdebug("RlpEncode()")

	fmt.Printf("rlpencode-%d-%d", len(l.Owner.Bytes()), len(l.Owner))

	fmt.Printf("rlpencode-%v-%v", l.Owner.Bytes(), l.Owner)

	return rlp.EncodeToBytes([]interface{}{
		uint64(l.Slot),
		l.PrevBlock.Uint64(),
		l.Denomination,
		l.Owner,
	})
}

func (l *LoomTx) Hash() []byte {
	//TODO is Previous block included block?
	//    if l.IncludeBlock.Mod(1000) == 0 {
	//            ret = w3.sha3(rlp.encode(self, UnsignedTransaction))
	//   }

	//      else
	data, err := soliditySha3(l.Slot)
	if err != nil {
		panic(err) //TODO
	}
	if len(data) != 32 {
		panic(fmt.Sprintf("wrong hash size! expected 32, got %v", len(data)))
	}
	return data
}

func (l *LoomTx) MerkleHash() []byte {
	data, err := l.rlpEncodeWithSha3()
	if err != nil {
		panic(err) //TODO
	}
	panic("TODO")

	return data
}

func soliditySha3(data uint64) ([]byte, error) {
	pairs := []*evmcompat.Pair{&evmcompat.Pair{"uint64", strconv.FormatUint(data, 10)}}
	hash, err := evmcompat.SoliditySHA3(pairs)
	if err != nil {
		return []byte{}, err
	}
	return hash, err
}

func (l *LoomTx) rlpEncodeWithSha3() ([]byte, error) {
	hash, err := l.RlpEncode()
	if err != nil {
		return []byte{}, err
	}
	d := sha3.NewKeccak256()
	d.Write(hash)
	return d.Sum(nil), nil
}
