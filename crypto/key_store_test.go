// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package crypto

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/randentropy"
)

func TestKeyStorePlain(t *testing.T) {
	ks := NewKeyStorePlain(common.DefaultDataDir())
	pass := "" // not used but required by API
	k1, err := ks.GenerateNewKey(randentropy.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}

	k2 := new(Key)
	k2, err = ks.GetKey(k1.Address, pass)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k2.Address, pass)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeyStorePassphrase(t *testing.T) {
	ks := NewKeyStorePassphrase(common.DefaultDataDir(), KDFStandard)
	pass := "foo"
	k1, err := ks.GenerateNewKey(randentropy.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	k2 := new(Key)
	k2, err = ks.GetKey(k1.Address, pass)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k2.Address, pass) // also to clean up created files
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeyStorePassphraseKDFLight(t *testing.T) {
	ks := NewKeyStorePassphrase(common.DefaultDataDir(), KDFLight)
	pass := "foo"
	k1, err := ks.GenerateNewKey(randentropy.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	k2 := new(Key)
	k2, err = ks.GetKey(k1.Address, pass)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k2.Address, pass) // also to clean up created files
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeyStorePassphraseDecryptionFail(t *testing.T) {
	ks := NewKeyStorePassphrase(common.DefaultDataDir(), KDFStandard)
	pass := "foo"
	k1, err := ks.GenerateNewKey(randentropy.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ks.GetKey(k1.Address, "bar") // wrong passphrase
	if err == nil {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k1.Address, "bar") // wrong passphrase
	if err == nil {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k1.Address, pass) // to clean up
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeyStorePassphraseKDFLightDecryptionFail(t *testing.T) {
	ks := NewKeyStorePassphrase(common.DefaultDataDir(), KDFLight)
	pass := "foo"
	k1, err := ks.GenerateNewKey(randentropy.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ks.GetKey(k1.Address, "bar") // wrong passphrase
	if err == nil {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k1.Address, "bar") // wrong passphrase
	if err == nil {
		t.Fatal(err)
	}

	err = ks.DeleteKey(k1.Address, pass) // to clean up
	if err != nil {
		t.Fatal(err)
	}
}

// Test and utils for the key store tests in the Ethereum JSON tests;
// tests/KeyStoreTests/basic_tests.json
type KeyStoreTest struct {
	Json     encryptedKeyJSON
	Password string
	Priv     string
}

func Test_PBKDF2_1(t *testing.T) {
	tests := loadKeyStoreTest("tests/test_vector.json", t)
	testDecrypt(tests["wikipage_test_vector_pbkdf2"], t)
}

func Test_PBKDF2_2(t *testing.T) {
	tests := loadKeyStoreTest("../tests/files/KeyStoreTests/basic_tests.json", t)
	testDecrypt(tests["test1"], t)
}

func Test_PBKDF2_3(t *testing.T) {
	tests := loadKeyStoreTest("../tests/files/KeyStoreTests/basic_tests.json", t)
	testDecrypt(tests["python_generated_test_with_odd_iv"], t)
}

func Test_PBKDF2_4(t *testing.T) {
	tests := loadKeyStoreTest("../tests/files/KeyStoreTests/basic_tests.json", t)
	testDecrypt(tests["evilnonce"], t)
}

func Test_Scrypt_1(t *testing.T) {
	tests := loadKeyStoreTest("tests/v3_test_vector.json", t)
	testDecrypt(tests["wikipage_test_vector_scrypt"], t)
}

func Test_Scrypt_2(t *testing.T) {
	tests := loadKeyStoreTest("../tests/files/KeyStoreTests/basic_tests.json", t)
	testDecrypt(tests["test2"], t)
}

func testDecrypt(test KeyStoreTest, t *testing.T) {
	privBytes, _, err := decryptKey(&test.Json, test.Password)
	if err != nil {
		t.Fatal(err)
	}
	privHex := hex.EncodeToString(privBytes)
	if test.Priv != privHex {
		t.Fatal(fmt.Errorf("Decrypted bytes not equal to test, expected %v have %v", test.Priv, privHex))
	}
}

func loadKeyStoreTest(file string, t *testing.T) map[string]KeyStoreTest {
	tests := make(map[string]KeyStoreTest)
	err := common.LoadJSON(file, &tests)
	if err != nil {
		t.Fatal(err)
	}
	return tests
}

func TestKeyForDirectICAP(t *testing.T) {
	key := NewKeyForDirectICAP(randentropy.Reader)
	if !strings.HasPrefix(key.Address.Hex(), "0x00") {
		t.Errorf("Expected first address byte to be zero, have: %s", key.Address.Hex())
	}
}
