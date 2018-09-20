package voucher

import (
	"os"
	"strconv"
	"testing"
)

const snakeoilKeyID = "1E92E2B4BB73E885"
const snakeoilKeyFingerprint = "90E942641C07A4C466BA97161E92E2B4BB73E885"
const testSignedValue = "test value to sign"

func newTestKeyRing(t *testing.T) *KeyRing {
	t.Helper()

	keyring := NewKeyRing()

	keyFile, err := os.Open("testdata/testkey.asc")
	if nil != err {
		t.Fatalf("failed to open key file: %s", err)
	}

	defer keyFile.Close()

	err = AddKeyToKeyRingFromReader(keyring, "snakeoil", keyFile)
	if nil != err {
		t.Fatalf("Failed to add key to keyring: %s", err)
	}

	return keyring
}

func getTestKeyID(t *testing.T) uint64 {
	t.Helper()

	keyID, err := strconv.ParseUint(snakeoilKeyID, 16, 64)
	if nil != err {
		t.Fatalf("Failed to convert snakeoilKeyID to uint: %s", err)
	}

	return keyID
}

func TestGetKeyAndSign(t *testing.T) {

	keyring := newTestKeyRing(t)

	entity, err := keyring.GetSignerByName("snakeoil")
	if nil != err {
		t.Fatalf("Failed to get signing key from KeyRing: %s", err)
	}

	if nil == entity.PrimaryKey {
		t.Fatalf("Failed to get private key from KeyRing.")
	}

	keyID := getTestKeyID(t)

	if entity.PrimaryKey.KeyId != keyID {
		t.Fatalf("Failed to get same key ID from KeyRing: \n%d vs \n%d", entity.PrimaryKey.KeyId, keyID)
	}

	signedValue, err := Sign(entity, testSignedValue)
	if nil != err {
		t.Fatalf("Failed to sign message: %s", err)
	}

	_, err = Verify(keyring, signedValue)
	if nil != err {
		t.Fatalf("Failed to verify signed message: %s", err)
	}
}

func TestOpenpgpKeyRing(t *testing.T) {

	keyring := newTestKeyRing(t)

	keyID := getTestKeyID(t)

	keys := keyring.KeysById(keyID)
	if len(keys) != 1 {
		t.Fatal("too many keys returned by KeysByID")
	}

	for _, key := range keys {
		if key.PublicKey.KeyId != keyID {
			t.Errorf("returned key that shouldn't have been, key ID is %X, should be %s", key.PublicKey.Fingerprint, snakeoilKeyFingerprint)
		}
	}

	encKeys := keyring.DecryptionKeys()
	if len(encKeys) != 0 {
		t.Fatal("too many keys returned by DecryptionKeys")
	}

}
