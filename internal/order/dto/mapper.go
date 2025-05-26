package dto

import "github.com/MingPV/clean-go-template/internal/entities"

func ToOrderResponse(order *entities.Order) *OrderResponse {
	return &OrderResponse{
		ID:    order.ID,
		Total: order.Total,
	}
}

func ToOrderResponseList(orders []*entities.Order) []*OrderResponse {
	result := make([]*OrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, ToOrderResponse(o))
	}
	return result
}
