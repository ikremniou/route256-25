package contest

import "strings"

type ProductValidationInput struct {
	NumberOfProducts int
	Products         map[string]string
	ProductPrices    map[string]bool
	ValidationString string
}

func ProductValidation(input []ProductValidationInput) []string {
	var result = make([]string, len(input))
	for index, productInput := range input {
		var validationParts = strings.Split(productInput.ValidationString, ",")
		var isPartValid = true

		for i := 0; i < len(validationParts) && isPartValid; i++ {
			var productAndPrice = strings.Split(validationParts[i], ":")
			if len(productAndPrice) != 2 {
				isPartValid = false
				break
			}

			var productPrice, isFound = productInput.Products[productAndPrice[0]]

			if !isFound {
				isPartValid = false
			} else {
				// compare the price
				if productPrice != productAndPrice[1] {
					isPartValid = false
				}

				var productPriceStatus, _ = productInput.ProductPrices[productPrice]
				if productPriceStatus {
					isPartValid = false
				} else {
					productInput.ProductPrices[productPrice] = true
				}
			}
		}

		if isPartValid {
			for _, price := range productInput.ProductPrices {
				if !price {
					isPartValid = false
				}
			}
		}

		if isPartValid {
			result[index] = "YES"
		} else {
			result[index] = "NO"
		}
	}

	return result
}
