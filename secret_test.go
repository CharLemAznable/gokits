package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

const PasswordKey = "192006250b4c09247ec02edce69f6a2d"

const AESCipherKey = "0b4c09247ec02edc"

func TestHmacSha256Base64(t *testing.T) {
    sign := HmacSha256Base64("Abc123", PasswordKey)
    assert.Equal(t, "RSbCrv07dc+f9NffWnaz4/p07z0oXL+u6Jtjl7XK6Bg=", sign)
}

func TestAESEncrypt(t *testing.T) {
    encrypted := AESEncrypt("The quick brown fox jumps over the lazy dog", AESCipherKey)
    assert.Equal(t, "3781dU72kqM+ulqyVv7aQlEoowO5jjGkTIjNNPKILa06LZ61DrAl7bhFFR20Ioao", encrypted)
}

func TestAESDecrypt(t *testing.T) {
    decrypted := AESDecrypt("3781dU72kqM+ulqyVv7aQlEoowO5jjGkTIjNNPKILa06LZ61DrAl7bhFFR20Ioao", AESCipherKey)
    assert.Equal(t, "The quick brown fox jumps over the lazy dog", decrypted)
}
