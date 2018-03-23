package command

import (
	"io/ioutil"
	"testing"
)

func TestTryDockerCompose(t *testing.T) {
	tests := []struct {
		fileName        string
		projectName     string
		expectedContent string
		expectedError   string
	}{
		{
			fileName:      "sloppy.yml",
			projectName:   "apache",
			expectedError: "no conversion required for sloppy files",
		},
		{
			fileName:        "docker.yml",
			projectName:     "foobar",
			expectedContent: "project: foobar\nservices:\n  apps:\n    foo:\n      image: golang\n      port: 8080\nversion: v1\n",
		},
	}

	for idx, test := range tests {
		reader, err := tryDockerCompose("testdata/"+test.fileName, test.projectName)

		if err != nil {
			if test.expectedError == "" {
				t.Fatalf(
					"#%d Did not expect error.\nReceived: %s",
					idx,
					err,
				)
			} else if err.Error() != test.expectedError {
				t.Fatalf(
					"#%d Expected different error.\nReceived: %s\nExpected: %s",
					idx,
					err,
					test.expectedError,
				)
			}
		}

		if reader == nil {
			if test.expectedContent != "" {
				t.Fatalf(
					"#%d Did not receive a reader.\nExpected: \n%s",
					idx,
					test.expectedContent,
				)
			}
		} else {
			bs, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatalf(
					"#%d Did not expect error while reading from reader.\nReceived: %s",
					idx,
					err,
				)
			}

			if s := string(bs); s != test.expectedContent {
				t.Fatalf(
					"#%d Content mismatch.\nReceived: %q\nExpected: %q",
					idx,
					s,
					test.expectedContent,
				)
			}
		}
	}
}
