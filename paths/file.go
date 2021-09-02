package paths

import (
	"bytes"
	"encoding/json"
	"fmt"

	vgcrypto "code.vegaprotocol.io/shared/libs/crypto"
	vgfs "code.vegaprotocol.io/shared/libs/fs"

	"github.com/zannen/toml"
)

func ReadStructuredFile(path string, v interface{}) error {
	buf, err := vgfs.ReadFile(path)
	if err != nil {
		return fmt.Errorf("couldn't read file at %s: %w", path, err)
	}

	if _, err := toml.Decode(string(buf), v); err != nil {
		return fmt.Errorf("couldn't decode buffer: %w", err)
	}

	return nil
}

func WriteStructuredFile(path string, v interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(v); err != nil {
		return fmt.Errorf("couldn't encode buffer: %w", err)
	}

	if err := vgfs.WriteFile(path, buf.Bytes()); err != nil {
		return fmt.Errorf("couldn't write file at %s: %w", path, err)
	}

	return nil
}

func ReadEncryptedFile(path string, passphrase string, v interface{}) error {
	encryptedBuf, err := vgfs.ReadFile(path)
	if err != nil {
		return fmt.Errorf("couldn't read secure file at %s: %w", path, err)
	}

	buf, err := vgcrypto.Decrypt(encryptedBuf, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't decrypt buffer: %w", err)
	}

	err = json.Unmarshal(buf, v)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal object: %w", err)
	}

	return nil
}

func WriteEncryptedFile(path string, passphrase string, v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("couldn't marshal object: %w", err)
	}

	encryptedBuf, err := vgcrypto.Encrypt(buf, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't encrypt buffer: %w", err)
	}

	if err := vgfs.WriteFile(path, encryptedBuf); err != nil {
		return fmt.Errorf("couldn't write secure file at %s: %w", path, err)
	}

	return nil
}
