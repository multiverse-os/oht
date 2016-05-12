package ricochet

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

type Ricochet struct {
	onionPrivateKey *rsa.PrivateKey
	onionHost       string
}

func (r *Ricochet) LoadPrivateKey(path string) error {
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return err
	}
	r.onionPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	return nil
}

func (r *Ricochet) OnAuthenticationChallenge(channelID int32, remoteOnionHost string, serverCookie [16]byte) {
	r.authHandler[onionHost].AddServerCookie(serverCookie[:])
	// DER Encode the Public Key
	publickeyBytes, _ := asn1.Marshal(rsa.PublicKey{
		N: r.privateKey.PublicKey.N,
		E: sr.privateKey.PublicKey.E,
	})

	signature, _ := rsa.SignPKCS1v15(nil, r.privateKey, crypto.SHA256, r.authHandler[onionHost].GenChallenge(r.onionHost, remoteOnionHost))
	signatureBytes := make([]byte, 128)
	copy(signatureBytes[:], signature[:])
	r.ricochet.SendProof(1, publickeyBytes, signatureBytes)
}
