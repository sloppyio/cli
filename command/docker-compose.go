package command

import (
	"bytes"
	"io"

	"github.com/ghodss/yaml"

	"github.com/sloppyio/sloppose/pkg/converter"
)

func tryDockerCompose(fileName, projectName string) (reader io.Reader, err error) {
	r := &converter.ComposeReader{}
	buf, err := r.Read(fileName)
	if err != nil {
		return nil, err
	}

	// detect sloppy file
	var test = struct {
		Version string `json:"version,omitempty"`
		Project string `json:"project,omitempty"`
	}{}
	err = yaml.Unmarshal(buf, &test)
	if err != nil {
		return nil, err
	}
	if test.Project != "" && test.Version == "v1" {
		// sloppy file format detected no error, but no new reader as well
		return nil, nil
	}

	cf, err := converter.NewComposeFile(buf, projectName)
	if err != nil {
		return nil, err
	}

	sf, err := converter.NewSloppyFile(cf)
	if err != nil {
		return nil, err
	}

	linker := &converter.Linker{}
	err = linker.Resolve(cf, sf)
	if err != nil {
		return nil, err
	}

	bs, err := yaml.Marshal(sf)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bs), nil
}
