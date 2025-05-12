package repository

import (
	"context"
	"math"
	"route256/cart/internal/domain/model"
	"slices"
	"sort"
	"sync"

	"go.opentelemetry.io/otel"
)

type CartRepository struct {
	cartItems map[int64][]*model.CartItemModel
	lock      sync.RWMutex
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		cartItems: make(map[int64][]*model.CartItemModel),
	}
}

// CreateItem implements CartRepository.
func (c *CartRepository) CreateItem(ctx context.Context, item *model.CartItemModel) (bool, error) {
	_, span := otel.Tracer("repository").Start(ctx, "cart_repository.CreateItem")
	defer span.End()

	c.lock.Lock()
	defer c.lock.Unlock()

	cartItems, hasEntry := c.cartItems[item.UserId]
	if !hasEntry {
		return c.createNewUserEntry(item)
	}

	return c.addToExistingCart(cartItems, item)
}

// DeleteBySku implements CartRepository.
func (c *CartRepository) DeleteBySku(ctx context.Context, userId int64, skuId int64) {
	_, span := otel.Tracer("repository").Start(ctx, "cart_repository.DeleteBySku")
	defer span.End()

	c.lock.Lock()
	defer c.lock.Unlock()

	cartItems, hasEntry := c.cartItems[userId]
	if hasEntry {
		var item = slices.IndexFunc(cartItems, func(item *model.CartItemModel) bool {
			return item.SkuId == skuId
		})

		if item == -1 {
			return
		}

		cartItems = slices.Delete(cartItems, item, item+1)

		if len(cartItems) == 0 {
			delete(c.cartItems, userId)
		} else {
			c.cartItems[userId] = cartItems
		}
	}
}

// DeleteAll implements CartRepository.
func (c *CartRepository) DeleteAll(ctx context.Context, userId int64) {
	_, span := otel.Tracer("repository").Start(ctx, "cart_repository.DeleteAll")
	defer span.End()

	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.cartItems, userId)
}

// GetAllOrderBySku implements CartRepository.
func (c *CartRepository) GetAllOrderBySku(ctx context.Context, userId int64) []model.CartItemModel {
	_, span := otel.Tracer("repository").Start(ctx, "cart_repository.GetAllOrderBySku")
	defer span.End()

	var items []model.CartItemModel = c.GetByUserId(ctx, userId)
	if items == nil {
		return nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SkuId < items[j].SkuId
	})

	return items
}

func (c *CartRepository) GetByUserId(ctx context.Context, userId int64) []model.CartItemModel {
	_, span := otel.Tracer("repository").Start(ctx, "cart_repository.GetByUserId")
	defer span.End()

	c.lock.RLock()
	defer c.lock.RUnlock()

	cartItems, hasEntry := c.cartItems[userId]
	if !hasEntry {
		return nil
	}

	items := make([]model.CartItemModel, len(cartItems))
	for i, item := range cartItems {
		items[i] = *item
	}

	return items
}

func (c *CartRepository) getTotalItemCount() int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var total int64 = 0

	for _, cartItems := range c.cartItems {
		for _, item := range cartItems {
			total += int64(item.Count)
		}
	}

	return total
}

func (c *CartRepository) createNewUserEntry(item *model.CartItemModel) (bool, error) {
	newCartItems := make([]*model.CartItemModel, 1)
	newCartItems[0] = item
	c.cartItems[item.UserId] = newCartItems

	return true, nil
}

func (c *CartRepository) addToExistingCart(cartItems []*model.CartItemModel, item *model.CartItemModel) (bool, error) {
	var itemIndex = slices.IndexFunc(cartItems, func(currItem *model.CartItemModel) bool {
		return currItem.SkuId == item.SkuId
	})

	if itemIndex >= 0 {
		var newCount = uint64(cartItems[itemIndex].Count) + uint64(item.Count)
		if newCount > math.MaxUint32 {
			return false, model.ErrTotalCountExceeded
		}

		cartItems[itemIndex].Count += item.Count
	} else {
		cartItems = append(cartItems, item)
		c.cartItems[item.UserId] = cartItems
	}

	return itemIndex == -1, nil
}
