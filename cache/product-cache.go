package cache

import "redis-golang-cnwnc/entities"

type ProductCache interface {
	Set(key string, value *[]entities.Product)
	Get(key string) *[]entities.Product
}
