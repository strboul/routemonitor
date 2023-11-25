package main

import (
	"os"
	"testing"
)

func TestReadConfigMissingFile(t *testing.T) {
	_, err := ReadConfig("/dev/null/foo")
	expectedErrMsg := "cannot read file"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestReadConfigInvalidFile(t *testing.T) {
	_, err := ReadConfig("/dev/null")
	expectedErrMsg := "could not parse config file=\"/dev/null\" EOF"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestReadConfigSuccessful(t *testing.T) {
	content := []byte(`
route:
  - name: Route 1
    ip: 192.168.1.1
    expect:
      - when:
          device: Device 1
`)
	tempFile, err := os.CreateTemp("", "routemonitor-*.yml")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tempFile.Name())

	os.WriteFile(tempFile.Name(), content, 0644)
	config, configErr := ReadConfig(tempFile.Name())
	if configErr != nil {
		t.Errorf(err.Error())
	}

	route := config.Route[0]

	if route.Name != "Route 1" {
		t.Errorf("current %v, expected %v", route.Name, "Route 1")
	}

	if route.IP != "192.168.1.1" {
		t.Errorf("current %v, expected %v", route.IP, "192.168.1.1")
	}

	expectWhen := route.Expect[0].When
	if expectWhen.Device != "Device 1" {
		t.Errorf("current %v, expected %v", expectWhen.Device, "Device 1")
	}
}
