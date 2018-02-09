package dockerconfig

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/docker/cli/cli/config/configfile"
)

func Transform(reader io.Reader) (io.Reader, error) {
	var c configfile.ConfigFile
	err := json.NewDecoder(reader).Decode(&c)
	if err != nil {
		return nil, err
	}

	all, err := c.GetAllCredentials()
	if err != nil {
		return nil, err
	}

	for server, cred := range all {
		ac := c.AuthConfigs[server]
		ac.Auth = encodeAuth(cred.Username, cred.Password)
		c.AuthConfigs[server] = ac
	}

	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func encodeAuth(username, password string) string {
	if username == "" && password == "" {
		return ""
	}

	b := []byte(username + ":" + password)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(encoded, b)
	return string(encoded)
}