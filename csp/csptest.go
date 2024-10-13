package csp

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

const hashAlg = "GOST12_256"

type Config struct {
	Container string
}

type CSPTest struct {
	Container string
}

func NewCSPTest(cfg *Config) *CSPTest {
	return &CSPTest{
		Container: cfg.Container,
	}
}

func (o *CSPTest) Sign(file []byte) ([]byte, error) {
	stdin := bytes.NewBuffer(file)
	stderr := &bytes.Buffer{}

	name := uuid.New()
	infileName := name.String()
	outfileName := infileName + ".sgn"
	defer func() { _ = os.Remove(infileName) }()
	defer func() { _ = os.Remove(outfileName) }()

	err := os.WriteFile(infileName, file, 0o777)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		"csptest",
		"-keyset",
		"-sign", hashAlg,
		"-cont", o.Container,
		"-in", infileName,
		"-out", outfileName,
		"-keytype", "signature",
	)

	cmd.Stdin = stdin
	cmd.Stderr = stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	print(out)

	signedFileBytes, err := os.ReadFile(outfileName)
	if err != nil {
		return nil, err
	}

	return signedFileBytes, nil
}
