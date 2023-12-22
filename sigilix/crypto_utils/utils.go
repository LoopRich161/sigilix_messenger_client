package crypto_utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math/big"
)

var EllipticCurve = elliptic.P256()

var EllipticCurveBitSize = EllipticCurve.Params().BitSize // 256

func SerializeSignature(r, s *big.Int) []byte {
	keySizeInBytes := (EllipticCurveBitSize + 7) / 8
	ret := make([]byte, 2*keySizeInBytes)
	r.FillBytes(ret[:keySizeInBytes])
	s.FillBytes(ret[keySizeInBytes:])
	return ret
}

func DeserializeSignature(signature []byte) (*big.Int, *big.Int, error) {
	keySizeInBytes := (EllipticCurveBitSize + 7) / 8
	if len(signature) != 2*keySizeInBytes {
		return nil, nil, errors.New("signature is not the correct size")
	}
	r := new(big.Int).SetBytes(signature[:keySizeInBytes])
	s := new(big.Int).SetBytes(signature[keySizeInBytes:])
	return r, s, nil
}

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(EllipticCurve, crand.Reader)
}

func Base64ToBytes(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func BytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func PublicECDSAKeyFromBytes(data []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(EllipticCurve, data)
	if x == nil {
		return nil, errors.New("invalid ECDSA public key")
	}
	return &ecdsa.PublicKey{Curve: EllipticCurve, X: x, Y: y}, nil
}

func PublicECDSAKeyToBytes(key *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(key.Curve, key.X, key.Y)
}

func PublicRSAKeyToBytes(key *rsa.PublicKey) ([]byte, error) {
	// DER format with X.509 subjectPublicKeyInfo with PKCS#1
	return x509.MarshalPKIXPublicKey(key)
}

func MustPublicRSAKeyToBytes(key *rsa.PublicKey) []byte {
	bytes, err := PublicRSAKeyToBytes(key)
	if err != nil {
		panic(err)
	}
	return bytes
}

func PublicRSAKeyFromBytes(data []byte) (*rsa.PublicKey, error) {
	p, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	return p.(*rsa.PublicKey), nil
}

func MustPublicRSAKeyFromBytes(data []byte) *rsa.PublicKey {
	key, err := PublicRSAKeyFromBytes(data)
	if err != nil {
		panic(err)
	}
	return key
}

func HashData(data []byte) []byte {
	hashed := sha256.Sum256(data)
	return hashed[:]
}

func ValidateECDSASignature(pubKey *ecdsa.PublicKey, data []byte, signature []byte) (bool, error) {
	hashed := HashData(data)

	r, s, err := DeserializeSignature(signature)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(pubKey, hashed, r, s), nil
}

func ValidateECDSASignatureFromBase64(pubKeyBytes []byte, data []byte, signatureBase64 string) (bool, error) {
	pubKey, err := PublicECDSAKeyFromBytes(pubKeyBytes)
	if err != nil {
		return false, err
	}

	signature, err := Base64ToBytes(signatureBase64)
	if err != nil {
		return false, err
	}
	return ValidateECDSASignature(pubKey, data, signature)
}

func GenerateUserIdByPublicKey(publicKey *ecdsa.PublicKey) uint64 {
	hashBytes := sha256.Sum256(elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y))
	buf := make([]byte, 8)
	copy(buf[4:], hashBytes[:4])
	return binary.BigEndian.Uint64(buf)
}

func SignMessage(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := HashData(data)

	r, s, err := ecdsa.Sign(crand.Reader, privateKey, hashed)
	if err != nil {
		return nil, err
	}

	signature := SerializeSignature(r, s)
	return signature, nil
}

func SignMessageBase64(privateKey *ecdsa.PrivateKey, data []byte) (string, error) {
	signature, err := SignMessage(privateKey, data)
	if err != nil {
		return "", err
	}
	return BytesToBase64(signature), nil
}

func PrivateKeyToBytes(privateKey *ecdsa.PrivateKey) []byte {
	keySizeInBytes := (EllipticCurveBitSize + 7) / 8
	ret := make([]byte, keySizeInBytes)
	privateKey.D.FillBytes(ret)
	return append(ret, PublicECDSAKeyToBytes(&privateKey.PublicKey)...)
}

func PrivateKeyToBytesBase64(privateKey *ecdsa.PrivateKey) string {
	return BytesToBase64(PrivateKeyToBytes(privateKey))
}

func PrivateKeyFromBytes(data []byte) (*ecdsa.PrivateKey, error) {
	keySizeInBytes := (EllipticCurveBitSize + 7) / 8
	if len(data) != (3*keySizeInBytes)+1 { // 1 byte for the curve type
		return nil, errors.New("invalid ECDSA private key")
	}
	d := new(big.Int).SetBytes(data[:keySizeInBytes])
	x, y := elliptic.Unmarshal(EllipticCurve, data[keySizeInBytes:])
	if x == nil {
		return nil, errors.New("invalid ECDSA public key")
	}
	return &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: EllipticCurve, X: x, Y: y}, D: d}, nil
}

func PrivateKeyFromBytesBase64(data string) (*ecdsa.PrivateKey, error) {
	bytes, err := Base64ToBytes(data)
	if err != nil {
		return nil, err
	}
	return PrivateKeyFromBytes(bytes)
}

func EncryptMessageChunked(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// OAEP has a limit on the size of the data that can be encrypted. This limit is reflected
	// in 'maxDataChunkSize'. If the data is larger than this limit, it will be split into chunks.
	// Each chunk will be encrypted separately and the resulting encrypted chunks will be concatenated.
	// if the data is smaller than the limit, it will be encrypted as a single chunk.

	maxDataChunkSize := (pubKey.Size() - 2*sha256.Size - 2) - 1
	outputChunkSize := pubKey.Size()

	chunkCount := (len(data) + maxDataChunkSize - 1) / maxDataChunkSize // round up
	ret := make([]byte, 0, chunkCount*outputChunkSize)
	for i := 0; i < chunkCount; i++ {
		chunk := data[i*maxDataChunkSize:]
		if len(chunk) > maxDataChunkSize {
			chunk = chunk[:maxDataChunkSize]
		}
		encryptedChunk, err := rsa.EncryptOAEP(sha256.New(), crand.Reader, pubKey, chunk, nil)
		if err != nil {
			return nil, err
		}
		ret = append(ret, encryptedChunk...)
	}
	return ret, nil
}

func DecryptMessageChunked(privKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	dataChunkSize := (privKey.Size() - 2*sha256.Size - 2) - 1
	outputChunkSize := privKey.Size()

	if len(data)%outputChunkSize != 0 {
		return nil, errors.New("invalid data size")
	}
	chunkCount := len(data) / outputChunkSize

	ret := make([]byte, 0, chunkCount*dataChunkSize)
	for i := 0; i < chunkCount; i++ {
		chunk := data[i*outputChunkSize:]
		if len(chunk) > outputChunkSize {
			chunk = chunk[:outputChunkSize]
		}
		decryptedChunk, err := rsa.DecryptOAEP(sha256.New(), crand.Reader, privKey, chunk, nil)
		if err != nil {
			return nil, err
		}
		ret = append(ret, decryptedChunk...)
	}

	return ret, nil
}

func EncryptMessage(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// use OAEP
	//return rsa.EncryptOAEP(sha256.New(), crand.Reader, pubKey, data, nil)
	return EncryptMessageChunked(pubKey, data)
}

func DecryptMessage(privKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	// use OAEP
	//return rsa.DecryptOAEP(sha256.New(), crand.Reader, privKey, data, nil)
	return DecryptMessageChunked(privKey, data)
}

func NewRSAKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(crand.Reader, 2048)
}

func RsaPrivateToBytes(privKey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(privKey)
}

func RsaPrivateFromBytes(data []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(data)
}
