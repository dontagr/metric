package models_test

import (
	"fmt"

	"github.com/dontagr/metric/models"
)

func Example() {
	delta := int64(123)

	new := models.Metrics{
		ID:    "1",
		MType: models.Counter,
		Delta: &delta,
		Hash:  "test",
	}

	fmt.Println(new.Hash)

	// Output:
	// test
}
