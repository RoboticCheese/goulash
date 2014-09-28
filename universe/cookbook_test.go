package universe

import (
	"github.com/RoboticCheese/goulash/universe"
	"testing"
)

func data1() (data1 universe.Cookbook) {
	data1 = universe.Cookbook{
		"0.1.0": universe.CookbookVersion{
			LocationType: "opscode",
			LocationPath: "https://example1.com",
			DownloadURL:  "https://example1.com/dl1",
			Dependencies: map[string]string{
				"thing1": ">= 0.0.0",
			},
		},
	}
	return
}

func data2() (data2 universe.Universe) {
	data2 = universe.Cookbook{
		"0.1.0": universe.CookbookVersion{
			LocationType: "opscode",
			LocationPath: "https://example1.com",
			DownloadURL:  "https://example1.com/dl1",
			Dependencies: map[string]string{
				"thing1": ">= 0.0.0",
			},
		},
	}
	return
}

func Test_Equals_1_Equal(t *testing.T) {
	data1 := data1()
	data2 := data2()
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Equals_2_MoreVersions(t *testing.T) {
	data1 := data1()
	data2 := data2()
	data2["0.2.0"] = universe.CookbookVersion{}
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_3_FewerVersions(t *testing.T) {
	data1 := data1()
	data2 := data2()
	data2 = universe.Cookbook{}
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_4_DifferentVersions(t *testing.T) {
	data1 := data1()
	data2 := data2()
	data2["0.1.0"] = universe.CookbookVersion{}
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}
