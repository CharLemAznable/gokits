package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestKeyPair(t *testing.T) {
    a := assert.New(t)

    plainText := "{ mac=\"MAC Address\", appId=\"16位字符串\", signature=SHA1(\"appId=xxx&mac=yyy\") }"
    keyPair, _ := GenerateRSAKeyPairDefault()
    privateKeyString, _ := keyPair.RSAPrivateKeyEncoded()
    publicKeyString, _ := keyPair.RSAPublicKeyEncoded()

    privateKey, _ := RSAPrivateKeyDecoded(privateKeyString)
    publicKey, _ := RSAPublicKeyDecoded(publicKeyString)

    cipherBytes, _ := EncryptByRSAKeyString([]byte(plainText), publicKeyString)
    plainBytes, _ := DecryptByRSAKeyString(cipherBytes, privateKeyString)
    a.Equal(plainText, string(plainBytes))

    cipherBytes, _ = EncryptByRSAKey([]byte(plainText), publicKey)
    plainBytes, _ = DecryptByRSAKey(cipherBytes, privateKey)
    a.Equal(plainText, string(plainBytes))

    pair := &RSAKeyPair{}
    _, errPrv := pair.RSAPrivateKeyEncoded()
    a.NotNil(errPrv)
    _, errPub := pair.RSAPublicKeyEncoded()
    a.NotNil(errPub)
}

func TestSigner(t *testing.T) {
    a := assert.New(t)

    plainText := "{ mac=\"MAC Address\", appId=\"16位字符串\", signature=SHA1(\"appId=xxx&mac=yyy\") }"
    keyPair, _ := GenerateRSAKeyPairDefault()
    privateKeyString, _ := keyPair.RSAPrivateKeyEncoded()
    publicKeyString, _ := keyPair.RSAPublicKeyEncoded()
    privateKey := keyPair.PrivateKey
    publicKey := keyPair.PublicKey

    sign1, _ := SHA1WithRSA.SignBase64ByRSAKeyString(plainText, privateKeyString)
    a.Nil(SHA1WithRSA.VerifyBase64ByRSAKeyString(plainText, sign1, publicKeyString))

    a.NotNil(SHA1WithRSA.VerifyBase64ByRSAKeyString(plainText, sign1, publicKeyString[1:]))
    a.NotNil(SHA1WithRSA.VerifyBase64ByRSAKeyString(plainText, sign1[1:], publicKeyString))
    _, err1 := SHA1WithRSA.SignBase64ByRSAKeyString(plainText, privateKeyString[1:])
    a.NotNil(err1)

    sign256, _ := SHA256WithRSA.SignBase64ByRSAKey(plainText, privateKey)
    a.Nil(SHA256WithRSA.VerifyBase64ByRSAKey(plainText, sign256, publicKey))

    a.NotNil(SHA256WithRSA.VerifyBase64ByRSAKey(plainText, sign256[1:], publicKey))
}
