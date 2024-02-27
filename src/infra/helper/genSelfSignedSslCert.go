package infraHelper

import "errors"

func GenSelfSignedSslCert(dirPath string, hostname string) error {
	keyFilePath := dirPath + "/" + hostname + ".key"
	certFilePath := dirPath + "/" + hostname + ".crt"

	_, err := RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		keyFilePath,
		"-out",
		certFilePath,
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+hostname,
	)
	if err != nil {
		return errors.New("CreateSelfSignedSslFailed: " + err.Error())
	}

	return nil
}
