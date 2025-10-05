package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	MaxConcurrent     int
	FfmpegTimeout     time.Duration
}

func init() {
	_ = godotenv.Load()
}

var config Config

func LoadConfig() Config {
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}

	maxConcurrentStr := os.Getenv("MAX_CONCURRENT")
	if maxConcurrentStr == "" {
		config.MaxConcurrent = 10
	} else {
		config.MaxConcurrent, _ = strconv.Atoi(maxConcurrentStr)
	}

	timeoutStr := os.Getenv("FFMPEG_TIMEOUT")
	if timeoutStr == "" {
		config.FfmpegTimeout = 30 * time.Second
	} else {
		d, err := time.ParseDuration(timeoutStr)
		if err != nil {
			fmt.Printf("⚠️ Invalid FFMPEG_TIMEOUT: %s. Using 30s\n", timeoutStr)
			config.FfmpegTimeout = 30 * time.Second
		} else {
			config.FfmpegTimeout = d
		}
	}
	fmt.Printf("✅ Config loaded:\n")
	fmt.Printf("   PORT: %s\n", config.Port)
	fmt.Printf("   MAX_CONCURRENT: %d\n", config.MaxConcurrent)
	fmt.Printf("   FFMPEG_TIMEOUT: %v\n", config.FfmpegTimeout)
	return config

}

func GetConfig() Config {
	if config.Port == "" {
		return LoadConfig()
	} else {
		return config
	}
}