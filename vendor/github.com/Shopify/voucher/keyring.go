package voucher

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

var errEmptyKeyring = errors.New("keyring is empty")
var errNoKeysInEjson = errors.New("no keys in ejson file")
var errKeysNotMap = errors.New("keys in ejson is not a map")

// KeyRing wraps an OpenPGP EntityList (which implements openpgp.KeyRing),
// adding support for determining which key is associated with which check.
// KeyRing implements openpgp.KeyRing, thus can be used in place of it where
// appropriate.
type KeyRing struct {
	keyIds   map[string]uint64
	entities openpgp.EntityList
}

// KeysById returns the set of keys that have the given key id.
func (keyring *KeyRing) KeysById(id uint64) []openpgp.Key {
	return keyring.entities.KeysById(id)
}

// KeysByIdUsage returns the set of keys with the given id
// that also meet the key usage given by requiredUsage.
// The requiredUsage is expressed as the bitwise-OR of
// packet.KeyFlag* values.
func (keyring *KeyRing) KeysByIdUsage(id uint64, requiredUsage byte) []openpgp.Key {
	return keyring.entities.KeysByIdUsage(id, requiredUsage)
}

// DecryptionKeys returns all private keys that are valid for
// decryption.
func (keyring *KeyRing) DecryptionKeys() []openpgp.Key {
	return keyring.entities.DecryptionKeys()
}

// GetSignerByName gets the first available signing key associated with the passed name.
func (keyring *KeyRing) GetSignerByName(name string) (*openpgp.Entity, error) {
	keyID := keyring.keyIds[name]

	// ensure we only get PGP keys that are specifically configured for signing.
	for _, key := range keyring.KeysByIdUsage(keyID, packet.KeyFlagSign) {
		if nil != key.PublicKey {
			if keyID == key.PublicKey.KeyId && nil != key.Entity {
				return key.Entity, nil
			}
		}
	}

	return nil, fmt.Errorf("no signing entity exists for check name \"%s\"", name)
}

// AddEntities adds new keys from the passed EntityList to the keyring for
// use.
func (keyring *KeyRing) AddEntities(name string, input openpgp.EntityList) {
	for _, entity := range input {
		if nil != entity.PrimaryKey {
			keyring.entities = append(keyring.entities, entity)
			keyring.keyIds[name] = entity.PrimaryKey.KeyId
		}
	}
}

// NewKeyRing creates a new keyring from the passed EntityList. The keys in
// the input EntityList are then associated with the
func NewKeyRing() *KeyRing {
	keyring := new(KeyRing)

	keyring.keyIds = make(map[string]uint64)
	keyring.entities = make(openpgp.EntityList, 0)

	return keyring
}

// AddKeyToKeyRingFromReader imports the PGP keys stored in a Reader into the
// passed KeyRing.
func AddKeyToKeyRingFromReader(keyring *KeyRing, name string, reader io.Reader) error {
	var err error
	var newKeyring openpgp.EntityList

	newKeyring, err = openpgp.ReadArmoredKeyRing(reader)
	if nil != err {
		return err
	}

	keyring.AddEntities(name, newKeyring)

	return nil
}
