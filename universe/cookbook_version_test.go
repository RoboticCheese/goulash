package universe

import (
	"testing"
)

func cvdata1() (data1 CookbookVersion) {
	data1 = CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
		},
	}
	return
}

func cvdata2() (data2 CookbookVersion) {
	data2 = CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
		},
	}
	return
}

func Test_CVEquals_1_Equal(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
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

func Test_CVEquals_2_DifferentLocationType(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.LocationType = "copsode"
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

func Test_CVEquals_3_DifferentLocationPath(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.LocationPath = "https://example2.com"
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

func Test_CVEquals_4_DifferentDownloadURL(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.DownloadURL = "https://example2.com/dl2"
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

func Test_CVEquals_5_DifferentDependenciesKeys(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.Dependencies["thing2"] = ">= 0.0.0"
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

func Test_CVEquals_6_DifferentDependenciesValues(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.Dependencies["thing1"] = "~> 1.0.0"
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
