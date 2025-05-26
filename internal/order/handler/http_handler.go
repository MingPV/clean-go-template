package handler

import (
	"strconv"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/order/usecase"
	response "github.com/MingPV/clean-go-template/pkg/responses"
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
	// var order entities.Order
	order := &entities.Order{}
	if err := c.BodyParser(order); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	if err := h.orderUseCase.CreateOrder(order); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(order)
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
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(orders)
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
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "id is required")
	}

	// Convert id to int
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	order, err := h.orderUseCase.FindOrderByID(orderID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(order)
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
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "id is required")
	}

	// Convert id to int
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	order := &entities.Order{}
	if err := c.BodyParser(&order); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	// Patch Order
	if err := h.orderUseCase.PatchOrder(orderID, order); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Fetch the updated order
	updatedOrder, err := h.orderUseCase.FindOrderByID(orderID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(updatedOrder)
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
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "id is required")
	}

	// Convert id to int
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.orderUseCase.DeleteOrder(orderID); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Message(c, fiber.StatusOK, "order deleted")
}
