package healthcheck

import "github.com/gofiber/fiber/v3"

type HealthResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

// HealthCheck godoc
// @Summary Check service health
// @Description Get service health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheckHeandler(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
		"message": "Server is running",
	})
}