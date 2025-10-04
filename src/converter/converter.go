package converter // ← пакет называется "utils"

import (

	"context"
	"fmt"
	"io"

	"mime/multipart"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)



func Convert(file *multipart.FileHeader, to string)  ([]byte, error){
    var input []byte
    src, err := file.Open()
    if err != nil {
        return nil, err
    }
    defer src.Close()

    input, err = io.ReadAll(src)
    if err != nil {
        return nil, err
    }

    id, err := uuid.NewV7()
    if err != nil {
        return nil, err
    }

    inputFileName := id.String() + "." + strings.Split(file.Filename, ".")[1]
    outpurFileName := id.String() + "." + to
    os.WriteFile(inputFileName, input, 0644)
    defer os.Remove(inputFileName)
    defer os.Remove(outpurFileName)


    cmd := exec.Command(
		"ffmpeg",
		"-i", inputFileName,
        "-c:a", "libmp3lame",
		"-b:a", "128k",
		"-ac", "2",
		"-ar", "44100",
		"-loglevel", "quiet",
		"-threads", "0",
		outpurFileName,
	)


	if err := cmd.Start(); err != nil {
		return nil, err
	}

    if err := cmd.Wait(); err != nil {
        fmt.Println(fmt.Errorf("ffmpeg process failed: %w", err))
		return nil, err
	}

    files, _ := os.ReadFile(outpurFileName)

    return files, err
}

func ConvertV2(file *multipart.FileHeader, _ string) ([]byte, error){
    var input []byte
    src, err := file.Open()
    if err != nil {
        return nil, err
    }
    defer src.Close()

    input, err = io.ReadAll(src)
    if err != nil {
        return nil, err
    }

    id, err := uuid.NewV7()
        if err != nil {
        return nil, err
    }
    inputFileName := id.String() + "." + strings.Split(file.Filename, ".")[1]
    defer os.Remove(inputFileName)
    os.WriteFile(inputFileName, input, 0644)
    

	// 5. Запускаем FFmpeg — теперь с **файлом на диске**
	cmd := exec.CommandContext(
		context.Background(),
		"ffmpeg",
		"-i", inputFileName, // ← FFmpeg теперь читает файл — НЕ pipe!
		"-c:a", "libopus",      // кодек Opus — ✅ ПРАВИЛЬНО
        "-b:a", "96k",          // битрейт: 96–128k — достаточно для музыки
        "-ar", "48000",         // ✅ ОБЯЗАТЕЛЬНО: 48 kHz — ЕДИНСТВЕННО ПОДДЕРЖИВАЕМАЯ ЧАСТОТА ДЛЯ OPUS
        "-application", "audio", // оптимизация для аудио (не голоса)
        "-frame_duration", "20", // оптимальная длительность фрейма
        "-compression_level", "10", // максимальное сжатие
        "-vbr", "on",           // переменный битрейт — лучше качество
        "-f", "webm",           // выходной контейнер
        "pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// 6. Запускаем процесс
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// 7. Читаем результат
	output, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read ffmpeg output: %w", err)
	}

	// 8. Ждём завершения — чтобы убедиться, что FFmpeg не упал
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ffmpeg process failed: %w", err)
	}
	return output, nil
}

