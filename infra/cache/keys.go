package cache

import "fmt"

const (
	productListPrefix = "product:list:"
)

func ProductById(id uint) string {
	return fmt.Sprintf("product:id:%d", id)
}

func ProductList(page, limit int) string {
	return fmt.Sprintf("%spage:%d:limit:%d", productListPrefix, page, limit)
}

func ProductListPrefix() string {
	return productListPrefix
}
