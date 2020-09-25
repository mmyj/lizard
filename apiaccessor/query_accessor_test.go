package apiaccessor

import (
	"errors"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestNewQueryAccessor(t *testing.T) {
	query := url.Values{}
	_, err := NewQueryAccessor(query, "123")
	assert.Equal(t, errors.Is(err, ErrArgLack), true)

	query = url.Values{
		nonceTag: []string{"12345"},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, errors.Is(err, ErrArgLack), true)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"12345"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"12345"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"a":          []string{"12345"},
		"b":          []string{"12345"},
		"c":          []string{"12345"},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
}

func TestCheckSignature(t *testing.T) {
	query := url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"12345"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err := NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckSignature()
	assert.Equal(t, errors.Is(err, ErrSignatureUnmatched), true)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{"12345"},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckSignature()
	assert.Equal(t, err, nil)
}

func TestCheckTimestamp(t *testing.T) {
	query := url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{"12345"},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err := NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckTimestamp()
	assert.Equal(t, errors.Is(err, ErrTimestampTimeout), true)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckTimestamp()
	assert.Equal(t, err, nil)
}

func TestCheckNonce(t *testing.T) {
	nonceMap := make(map[string]bool)
	mockNonceChecker := func(nonce string) error {
		if _, ok := nonceMap[nonce]; ok {
			return ErrNonceUsed
		}
		nonceMap[nonce] = true
		return nil
	}

	query := url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{"12345"},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err := NewQueryAccessor(query, "123", WithNonceChecker(mockNonceChecker))
	assert.Equal(t, err, nil)
	err = qa.CheckNonce()
	assert.Equal(t, err, nil)
	err = qa.CheckNonce()
	assert.Equal(t, errors.Is(err, ErrNonceUsed), true)
}

func TestCheckSignatureWithAfs(t *testing.T) {
	nonceMap := make(map[string]bool)
	mockNonceChecker := func(nonce string) error {
		if _, ok := nonceMap[nonce]; ok {
			return ErrNonceUsed
		}
		nonceMap[nonce] = true
		return nil
	}

	query := url.Values{
		nonceTag:         []string{"C0988D12-1BA2-41FA-865F-8DE270F9D85F"},
		signatureTag:     []string{"3efb19a26f64ccd0bac7e557c9a541df"},
		timestampTag:     []string{"1601028202"},
		"phone":          []string{"10700000490"},
		"afs_need":       []string{"1"},
		"afs_scene":      []string{"FFFF000000000178E014"},
		"afs_session_id": []string{"0152JIZgtMjy7iQLwB8JakWQphFFDc93Y1QFPgYL9RaZev50WXioPHBC61baDu7FHd4fNmwl2po70LSfvN6Wgm-tILksAsi_jOiYiArUQCcwdB868MvOl8tcUwL1pP4CnexPyktVeWJNkNKtAiXTEpgu6OwDPAEa64M5I7WpeSLKzWEDeT8BtEwD-ZOQ9xcY8286kToIZFsl19mhs_mLSHIA"},
		"afs_sig":        []string{"05a1C7nT4bR5hcbZlAujcdyaZLw91T0q9MLNR2wf8RL_T9aEQKj8zdHHb7-_Jl7mBAaionUk_C4JJYxFHMbjjc_znF3NrxhXR7GCTgV1H9HINiGp4-RSlo2wwzH2kqN58cyJAc-rseGkvNsN2xE2OyrE9pJhgUkPAe9marpHnnfb7TmctUUd7VcCO04fU54VmrblJzl3WydXL0RAPobz3qEKaS2Yc46UFkIg8oF-id3xmJ5qhY_MrcYaNOV0HVaB0v01Ar3YahFiLVfGqfBS59zjXSryGtw1MsTAMVpuGqiSoG6vx2BZmd8Ma7Ao-BOpoelZXfOp07Rp4cEumFSZ8HqiBowDKwosPH9ANv4aTkEtLx0X76doZjF3RB0bvQp12PYe7zToQu_oogYc6di7nrTz0ZFS1mz4z4HP5D9PkHdQRJNlnscrecaXC0Vy4bsder05gCwmvP5PB-3ilI1PdmQt-FUl2xBANnr9Zo97C3HfA"},
		"afs_token":      []string{"FFFF000000000178E014%3A1601028186726%3A0.5543301835657096"},
	}
	qa, err := NewQueryAccessor(query, "123", WithNonceChecker(mockNonceChecker))
	assert.Equal(t, err, nil)
	err = qa.CheckSignature()
	assert.Equal(t, err, nil)
}
