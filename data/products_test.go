package data

import "testing"

func TestChecksValidation(t *testing.T) {
	product := &Product{Name: "chai", Price: 1.01, SKU: "a-b-c"}
	err := product.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
