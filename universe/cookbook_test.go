package universe

import (
	"testing"
)

func cdata1() (data1 Cookbook) {
	data1 = Cookbook{
		"0.1.0": CookbookVersion{
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

func cdata2() (data2 Cookbook) {
	data2 = Cookbook{
		"0.1.0": CookbookVersion{
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

func Test_CEquals_1_Equal(t *testing.T) {
	data1 := cdata1()
	data2 := cdata2()
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

func Test_CEquals_2_MoreVersions(t *testing.T) {
	data1 := cdata1()
	data2 := cdata2()
	data2["0.2.0"] = CookbookVersion{}
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

func Test_CEquals_3_FewerVersions(t *testing.T) {
	data1 := cdata1()
	data2 := cdata2()
	data2 = Cookbook{}
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

func Test_CEquals_4_DifferentVersions(t *testing.T) {
	data1 := cdata1()
	data2 := cdata2()
	data2["0.1.0"] = CookbookVersion{}
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
