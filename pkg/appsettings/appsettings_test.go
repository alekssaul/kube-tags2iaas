package appsettings

import (
	"fmt"
	"os"
	"testing"
)

func TestGetAppSettings(t *testing.T) {
	os.Setenv("INFRASTRUCTURETAGS", `{"foo": "hello","bar": "world"}`)
	result, err := GetAppSettings()
	if err != nil {
		t.Fatalf("GetAppSettings: Could not initialize the application properly")
	}

	if fmt.Sprintf("%#v", result) == fmt.Sprintf("%#v", AppSettings{}) {
		t.Fatalf("GetAppSettings: Returned an empty AppSetting")
	}

	if result.InfrastructureTags[0].Key != "foo" {
		t.Fatalf("GetAppSettings: Did not initialize the Infrastructure Tags as expected")
	}

	os.Setenv("INFRASTRUCTURETAGS", `{"foo": "hello","bar": "world","zoo": "test"}`)
	result, err = GetAppSettings()
	if err == nil {
		t.Fatalf("GetAppSettings: Did not check for MaximumNumberofTags")
	}

	os.Exit(0)
}
