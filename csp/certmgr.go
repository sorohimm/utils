package csp

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/google/uuid"
)

func NewCertMGR() *CertMGR {
	return &CertMGR{}
}

type CertMGR struct{}

func (o *CertMGR) GetCertificateInfo(certHash string) (string, error) {
	cmd := exec.Command(
		"certmgr",
		"-list",
		"-thumbprint",
		certHash,
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (o *CertMGR) InstallPxfCertificate(cert []byte) (string, error) {
	certFilename := uuid.NewString() + ".pfx"
	defer func() { _ = os.Remove(certFilename) }()

	err := os.WriteFile(certFilename, cert, 0o777)
	if err != nil {
		return "", fmt.Errorf("write file error %w", err)
	}

	cmd := exec.Command(
		"certmgr",
		"-install",
		"-pfx",
		"-file",
		certFilename,
		"-stdin",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("run err: %w", err)
	}

	re := regexp.MustCompile(`: CN=([^,]+)`)

	result := re.FindStringSubmatch(string(output))

	if len(result) < 2 {
		return "", fmt.Errorf("cn not found")
	}

	container := `\\.\HDIMAGE\` + result[1]

	return container, nil
}
