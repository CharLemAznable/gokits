package gokits

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "errors"
)

type RSAKeyPair struct {
    PrivateKey *rsa.PrivateKey
    PublicKey  *rsa.PublicKey
}

func GenerateRSAKeyPairDefault() (*RSAKeyPair, error) {
    return GenerateRSAKeyPair(1024)
}

func GenerateRSAKeyPair(keySize int) (*RSAKeyPair, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
    if err != nil {
        return nil, err
    }
    publicKey := &privateKey.PublicKey
    return &RSAKeyPair{
        PrivateKey: privateKey,
        PublicKey:  publicKey}, nil
}

func (p *RSAKeyPair) RSAPrivateKeyEncoded() (string, error) {
    if nil == p.PrivateKey {
        return "", errors.New("PrivateKeyEmpty")
    }
    bytes, err := x509.MarshalPKCS8PrivateKey(p.PrivateKey)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

func (p *RSAKeyPair) RSAPublicKeyEncoded() (string, error) {
    if nil == p.PublicKey {
        return "", errors.New("PublicKeyEmpty")
    }
    bytes, err := x509.MarshalPKIXPublicKey(p.PublicKey)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

func RSAPrivateKeyDecoded(privateKeyString string) (*rsa.PrivateKey, error) {
    bytes, err := base64.StdEncoding.DecodeString(privateKeyString)
    if nil != err {
        return nil, err
    }
    privateKey, err := x509.ParsePKCS8PrivateKey(bytes)
    if nil != err {
        return nil, err
    }
    return privateKey.(*rsa.PrivateKey), nil
}

func RSAPublicKeyDecoded(publicKeyString string) (*rsa.PublicKey, error) {
    bytes, err := base64.StdEncoding.DecodeString(publicKeyString)
    if nil != err {
        return nil, err
    }
    publicKey, err := x509.ParsePKIXPublicKey(bytes)
    if nil != err {
        return nil, err
    }
    return publicKey.(*rsa.PublicKey), nil
}

func EncryptByRSAKeyString(plainBytes []byte, publicKeyString string) ([]byte, error) {
    publicKey, err := RSAPublicKeyDecoded(publicKeyString)
    if nil != err {
        return nil, err
    }
    return EncryptByRSAKey(plainBytes, publicKey)
}

func EncryptByRSAKey(plainBytes []byte, publicKey *rsa.PublicKey) ([]byte, error) {
    return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainBytes)
}

func DecryptByRSAKeyString(cipherBytes []byte, privateKeyString string) ([]byte, error) {
    privateKey, err := RSAPrivateKeyDecoded(privateKeyString)
    if nil != err {
        return nil, err
    }
    return DecryptByRSAKey(cipherBytes, privateKey)
}

func DecryptByRSAKey(cipherBytes []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
}

type Signer struct {
    hash crypto.Hash
}

func (s *Signer) SignBase64ByRSAKeyString(plainText, privateKeyString string) (string, error) {
    sign, err := s.SignByRSAKeyString(plainText, privateKeyString)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(sign), nil
}

func (s *Signer) SignBase64ByRSAKey(plainText string, privateKey *rsa.PrivateKey) (string, error) {
    sign, err := s.SignByRSAKey(plainText, privateKey)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(sign), nil
}

func (s *Signer) SignByRSAKeyString(plainText, privateKeyString string) ([]byte, error) {
    privateKey, err := RSAPrivateKeyDecoded(privateKeyString)
    if nil != err {
        return nil, err
    }
    return s.SignByRSAKey(plainText, privateKey)
}

func (s *Signer) SignByRSAKey(plainText string, privateKey *rsa.PrivateKey) ([]byte, error) {
    hash := s.hash.New()
    hash.Write([]byte(plainText))
    return rsa.SignPKCS1v15(rand.Reader, privateKey, s.hash, hash.Sum(nil))
}

func (s *Signer) VerifyBase64ByRSAKeyString(plainText, signText, publicKeyString string) error {
    sign, err := base64.StdEncoding.DecodeString(signText)
    if nil != err {
        return err
    }
    return s.VerifyByRSAKeyString(plainText, sign, publicKeyString)
}

func (s *Signer) VerifyBase64ByRSAKey(plainText, signText string, publicKey *rsa.PublicKey) error {
    sign, err := base64.StdEncoding.DecodeString(signText)
    if nil != err {
        return err
    }
    return s.VerifyByRSAKey(plainText, sign, publicKey)
}

func (s *Signer) VerifyByRSAKeyString(plainText string, sign []byte, publicKeyString string) error {
    publicKey, err := RSAPublicKeyDecoded(publicKeyString)
    if nil != err {
        return err
    }
    return s.VerifyByRSAKey(plainText, sign, publicKey)
}

func (s *Signer) VerifyByRSAKey(plainText string, sign []byte, publicKey *rsa.PublicKey) error {
    hash := s.hash.New()
    hash.Write([]byte(plainText))
    return rsa.VerifyPKCS1v15(publicKey, s.hash, hash.Sum(nil), sign)
}

var (
    SHA1WithRSA   = Signer{hash: crypto.SHA1}
    SHA256WithRSA = Signer{hash: crypto.SHA256}
)
