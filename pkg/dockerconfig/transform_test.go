package dockerconfig

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestTransform(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err string
	}{
		{
			`{}`,
			`{"auths":null}`,
			"",
		},
		{
			`{"auths":{}}`,
			`{"auths":{}}`,
			"",
		},
		{
			`{"auths":{"https://index.docker.io/v1/":{"auth":"aWxvcmVtOmlwc3Vt"}}}`,
			`{"auths":{"https://index.docker.io/v1/":{"auth":"aWxvcmVtOmlwc3Vt"}}}`,
			"",
		},
		// TODO: how to test / mock?
		//{
		//	`{"auths":{"https://index.docker.io/v1/":{}},"credsStore":"..."}`,
		//	`{"auths":{"https://index.docker.io/v1/":{"auth":"aWxvcmVtOmlwc3Vt"}}}`,
		//	"",
		//},
	}

	for i, test := range tests {
		input := bytes.NewBuffer([]byte(test.in))
		output, err := Transform(input)

		if err != nil && test.err == "" {
			t.Fatalf("#%d\tDid not expect an error: \n\treceived: \"%s\"", i, err)
		}

		if err != nil && err.Error() != test.err {
			t.Fatalf("#%d\tWrong error: \n\treceived: \"%s\", \n\texpected: \"%s\"", i, err, test.err)
		}

		if err == nil && test.err != "" {
			t.Fatalf("#%d\tExpected an error: \n\texpected: \"%s\"", i, test.err)
		}

		b, err := ioutil.ReadAll(output)
		if err != nil {
			t.Fatalf("#%d\tCould not read from (output io.Reader)", i)
		}

		s := string(b)
		if s != test.out {
			t.Fatalf("#%d\tWrong output: \n\treceived: \"%s\", \n\texpected: \"%s\"", i, s, test.out)
		}
	}
}
