package ibis


import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	// "github.com/rs/zerolog/log"

	// "github.com/rs/zerolog/log"
	"math/rand"
	"time"

	"gitlab.basis-plus.ru/ibis/proto"
	pb "gitlab.tesonero-computers.ru/ibis/authproto/pkg/endpoint"

	db "gitlab.tesonero-computers.ru/ibis/aaa/internal/data"
	_ "gitlab.tesonero-computers.ru/ibis/aaa/internal/model"
)

// Check authenticates peers
func auth(p *pb.EndpointParams) (*pb.EndpointParams, error) {

	rand.Seed(time.Now().UnixNano())

	fmt.Println(p.Account)
	// fmt.Println(p.AuthType)

	accIbis, err := db.GetIbisAccount(p.AuthType, p.Account)
	errorText := ""
	if err != nil {
		errorText = "db rows error"
		fmt.Println(err, errorText)
		return p, proto.ErrorStatuses[proto.StatusAuthError]
	}

	fmt.Println(accIbis, errorText)

	switch p.AuthType {
	case proto.AuthTypeText:

		expectedLogin  := []byte(accIbis.ClientLogin)
		expectedPasswd := []byte(accIbis.ClientPassword)
		serverLogin    := []byte(accIbis.ServerLogin)
		serverPasswd   := []byte(accIbis.ServerPassword)


		if !authText(expectedLogin, expectedPasswd, p.Rnd, p.Hash) {
			//fmt.Println("expectedLogin-", expectedLogin, "Rnd-", p.Rnd)
			//fmt.Println("expectedPasswd-", expectedPasswd, "Hash-", p.Hash)
			return p, proto.ErrorStatuses[proto.StatusAuthError]
		}

		p.Rnd  = serverLogin
		p.Hash = serverPasswd
		p.Timestamp = time.Now().UTC().Unix()

		return p, nil

	case proto.AuthTypeMD5:
		key := []byte(accIbis.PresharedKey)
		if !authMD5(p.Rnd, p.Hash, key, p.Timestamp) {
			// fmt.Println("MD5 AUT:ERROR - ", p.Account)
			return p, proto.ErrorStatuses[proto.StatusAuthError]
		}

		// fmt.Println("MD5 AUTH:OK - ", p.Account)

		p.Timestamp = time.Now().UTC().Unix()
		buf := make([]byte, 32)
		rand.Read(buf)
		p.Rnd = buf
		hash := hashMD5(p.Rnd, key, p.Timestamp)
		p.Hash = hash[:]

		return p, nil

	case proto.AuthTypeSHA256:
		key := []byte(accIbis.PresharedKey)
		if !authSHA256(p.Rnd, p.Hash, key, p.Timestamp) {
			return p, proto.ErrorStatuses[proto.StatusAuthError]
		}
		p.Timestamp = time.Now().UTC().Unix()
		buf := make([]byte, 32)
		rand.Read(buf)
		p.Rnd = buf
		hash := hashSHA256(p.Rnd, key, p.Timestamp)
		p.Hash = hash[:]

		return p, nil

	default:

		return p, proto.ErrorStatuses[proto.StatusAuthTypeMismatch]
	}
}

// Text authenticates the remote client by comparing the received login
// and password with saved ones
func authText(login, passwd, inLogin, inPasswd []byte) bool {
	var a, b, c, d [32]byte
	copy(a[:], login)
	copy(b[:], inLogin)
	copy(c[:], passwd)
	copy(d[:], inPasswd)

	if !(bytes.Compare(a[:], b[:]) == 0 &&
		bytes.Compare(c[:], d[:]) == 0) {
		return false
	}
	return true
}

// authMD5 authenticates the remote client by computing the hash
// and comparing with the received one
func authMD5(rnd, hash, key []byte, ts int64) bool {
	var a, b, c [32]byte
	copy(a[:], rnd)
	copy(b[:], key)
	copy(c[:], hash)
	computed := hashMD5(rnd, key, ts)
	if bytes.Compare(c[:], computed[:]) != 0 {
		return false
	}
	return true
}

// hashMD5 computes hash for rnd + key + timestamp
func hashMD5(rnd, key []byte, ts int64) (hash [32]byte) {
	var a, b [32]byte
	copy(a[:], rnd)
	copy(b[:], key)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(ts))

	sum := md5.Sum(append(append(a[:], b[:]...), buf...))
	copy(hash[:], sum[:])
	return
}

// authSHA256 authenticates the remote client by computing the hash
// and comparing with the received one
func authSHA256(rnd, hash, key []byte, ts int64) bool {
	var a, b, c [32]byte
	copy(a[:], rnd)
	copy(b[:], key)
	copy(c[:], hash)
	computed := hashSHA256(rnd, key, ts)
	if bytes.Compare(c[:], computed[:]) != 0 {
		return false
	}
	return true
}

// hashSHA256 computes hash for rnd + key + timestamp
func hashSHA256(rnd, key []byte, ts int64) (hash [32]byte) {
	var a, b [32]byte
	copy(a[:], rnd)
	copy(b[:], key)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(ts))

	sum := sha256.Sum256(append(append(a[:], b[:]...), buf...))
	copy(hash[:], sum[:])
	return
}



