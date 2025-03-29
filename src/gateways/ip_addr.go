package gateways

import (
	"go-fiber-template/src/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetIp(ctx *fiber.Ctx) error {
	data, err := h.IPService.GetIp()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get ip address"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}
