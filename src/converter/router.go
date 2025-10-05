package converter

import (
	"audioconverter/utils/config"

	"github.com/gofiber/fiber/v3"
)
var env = config.GetConfig()
var semaphore = make(chan struct{}, env.MaxConcurrent) // 10 concurrent connections

// @Summary Конвертация аудиофайла 11
// @Description Принимает аудиофайл и целевой формат, возвращает конвертированный файл
// @Tags Audio Processing
// @Accept multipart/form-data
// @Produce application/json
// @Param audio formData file true "Аудиофайл для конвертации"
// @Param to formData string false "Целевой формат (например, mp3, wav)" Format(string) default(mp3)
// @Success 200 {file} []byte "Конвертированный аудиофайл"
// @Failure 400 {object} map[string]string "Ошибка валидации или обработки"
// @Router /api/v1/convert [post]
// @Deprecated true
func ConverHeandler(c fiber.Ctx) error {
		file, err := c.FormFile("audio")
		
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "audio as file as requred",
			})
		}
		to := c.FormValue("to")
		if to == "" {
			to = "mp3"
		}
		
		c.Response().Header.Set("Content-Type", "audio/mpeg")
		out, err := Convert(file, to)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Send(out) 
}

// @Summary Конвертация аудиофайла v2 использовать для webm
// @Description Принимает аудиофайл и целевой формат, возвращает конвертированный файл
// @Tags Audio Processing
// @Accept multipart/form-data
// @Produce application/json
// @Param audio formData file true "Аудиофайл для конвертации"
// @Param to formData string false "Целевой формат (например, mp3, wav)" Format(string) default(mp3)
// @Success 200 {file} []byte "Конвертированный аудиофайл"
// @Failure 400 {object} map[string]string "Ошибка валидации или обработки"
// @Router /api/v2/convert [post]
func ConverV2Heandler(c fiber.Ctx) error {
	file, err := c.FormFile("audio")
		
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "audio as file as requred",
		})
	}

	to := c.FormValue("to")
	if to == "" {
		to = "webm"
	}
	// ✅ БЛОКИРУЕМ — ЖДЕМ, ПОКА ОСВОБОДИТСЯ МЕСТО В ПУЛЕ
	semaphore <- struct{}{} // ← БЛОКИРУЕТ, ПОКА НЕ БУДЕТ СВОБОДНОГО МЕСТА
	defer func() { <-semaphore }() // ← ОСВОБОЖДАЕМ МЕСТО, КОГДА ЗАДАЧА ЗАВЕРШИТСЯ
	
	c.Response().Header.Set("Content-Type", "audio/"+to)
	out, err := ConvertV2(file, to)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Send(out) 
}