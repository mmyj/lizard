package apiaccessor

import (
	"errors"
)

// Accessor is the interface use to check the availability of the request's arguments
type Accessor interface {
	CheckSignature() error
	CheckTimestamp() error
	CheckNonce() error
}

var (
	errArgLack            = errors.New("arg lack")
	errSignatureUnmatched = errors.New("signature is unmatched")
	errTimestampTimeout   = errors.New("timestamp time out")
	errNonceUsed          = errors.New("nonce is used")
)

const (
	nonceTag     = "nonce"
	timestampTag = "timestamp"
	secretKeyTag = "secret_key"
	signatureTag = "signature"
)

type arg struct {
	k string
	v string
}

type args struct {
	kv map[string]string
	l  []*arg
}

func newArgs() args {
	return args{
		kv: make(map[string]string),
		l:  make([]*arg, 0),
	}
}

func (a *args) append(k, v string) {
	a.kv[k] = v
	a.l = append(a.l, &arg{k: k, v: v})
}

// EvalSignature evaluating the signature of the request's arguments
type EvalSignature func(origin string) (signature string)

// TimestampChecker checking the availability of the request's timestamp argument
type TimestampChecker func(timestamp int64) error

// NonceChecker checking the availability of the request's nonce argument
type NonceChecker func(nonce string) error
