package grpc

import (
	"context"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/usecase"
	"github.com/MingPV/clean-go-template/pkg/apperror"
	orderpb "github.com/MingPV/clean-go-template/proto/order"
	"google.golang.org/grpc/status"
)

type GrpcOrderHandler struct {
	orderUseCase usecase.OrderUseCase
	orderpb.UnimplementedOrderServiceServer
}

func NewGrpcOrderHandler(uc usecase.OrderUseCase) *GrpcOrderHandler {
	return &GrpcOrderHandler{orderUseCase: uc}
}

func (h *GrpcOrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	order := &entities.Order{Total: float64(req.Total)}
	if err := h.orderUseCase.CreateOrder(order); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &orderpb.CreateOrderResponse{Order: toProtoOrder(order)}, nil
}

func (h *GrpcOrderHandler) FindOrderByID(ctx context.Context, req *orderpb.FindOrderByIDRequest) (*orderpb.FindOrderByIDResponse, error) {
	order, err := h.orderUseCase.FindOrderByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &orderpb.FindOrderByIDResponse{Order: toProtoOrder(order)}, nil
}

func (h *GrpcOrderHandler) FindAllOrders(ctx context.Context, req *orderpb.FindAllOrdersRequest) (*orderpb.FindAllOrdersResponse, error) {
	orders, err := h.orderUseCase.FindAllOrders()
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoOrders []*orderpb.Order
	for _, o := range orders {
		protoOrders = append(protoOrders, toProtoOrder(o))
	}

	return &orderpb.FindAllOrdersResponse{Orders: protoOrders}, nil
}

func (h *GrpcOrderHandler) PatchOrder(ctx context.Context, req *orderpb.PatchOrderRequest) (*orderpb.PatchOrderResponse, error) {
	order := &entities.Order{Total: float64(req.Total)}
	updatedOrder, err := h.orderUseCase.PatchOrder(int(req.Id), order)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &orderpb.PatchOrderResponse{Order: toProtoOrder(updatedOrder)}, nil
}

func (h *GrpcOrderHandler) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	if err := h.orderUseCase.DeleteOrder(int(req.Id)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &orderpb.DeleteOrderResponse{Message: "order deleted"}, nil
}

// helper function convert entities.Order to orderpb.Order
func toProtoOrder(o *entities.Order) *orderpb.Order {
	return &orderpb.Order{
		Id:    int32(o.ID),
		Total: float64(o.Total),
	}
}
