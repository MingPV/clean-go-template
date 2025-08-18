package rest

import (
	"strconv"

	"github.com/MingPV/clean-go-template/pkg/apperror"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/dto"
	"github.com/MingPV/clean-go-template/internal/order/usecase"
	responses "github.com/MingPV/clean-go-template/pkg/responses"
	"github.com/gofiber/fiber/v2"
)

type HttpOrderHandler struct {
	orderUseCase usecase.OrderUseCase
}

func NewHttpOrderHandler(useCase usecase.OrderUseCase) *HttpOrderHandler {
	return &HttpOrderHandler{orderUseCase: useCase}
}

// CreateOrder godoc
// @Summary Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body entities.Order true "Order payload"
// @Success 201 {object} entities.Order
// @Router /orders [post]
func (h *HttpOrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.ErrorWithMessage(c, err, "invalid request")
	}

	order := &entities.Order{Total: req.Total}
	if err := h.orderUseCase.CreateOrder(order); err != nil {
		return responses.Error(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToOrderResponse(order))
}

// FindAllOrders godoc
// @Summary Get all orders
// @Tags orders
// @Produce json
// @Success 200 {array} entities.Order
// @Router /orders [get]
func (h *HttpOrderHandler) FindAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderUseCase.FindAllOrders()
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToOrderResponseList(orders))
}

// FindOrderByID godoc
// @Summary Get order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} entities.Order
// @Router /orders/{id} [get]
func (h *HttpOrderHandler) FindOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return responses.ErrorWithMessage(c, err, "invalid id")
	}

	order, err := h.orderUseCase.FindOrderByID(orderID)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToOrderResponse(order))
}

// PatchOrder godoc
// @Summary Update an order partially
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param order body entities.Order true "Order update payload"
// @Success 200 {object} entities.Order
// @Router /orders/{id} [patch]
func (h *HttpOrderHandler) PatchOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return responses.ErrorWithMessage(c, err, "invalid id")
	}

	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.ErrorWithMessage(c, err, "invalid request")
	}

	order := &entities.Order{Total: req.Total}

	msg, err := validatePatchOrder(order)
	if err != nil {
		return responses.ErrorWithMessage(c, err, msg)
	}

	updatedOrder, err := h.orderUseCase.PatchOrder(orderID, order)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToOrderResponse(updatedOrder))
}

// DeleteOrder godoc
// @Summary Delete an order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.MessageResponse
// @Router /orders/{id} [delete]
func (h *HttpOrderHandler) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return responses.ErrorWithMessage(c, err, "invalid id")
	}

	if err := h.orderUseCase.DeleteOrder(orderID); err != nil {
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "order deleted")
}

func validatePatchOrder(order *entities.Order) (string, error) {

	if order.Total <= 0 {
		return "total must be positive", apperror.ErrInvalidData
	}

	return "", nil
}
