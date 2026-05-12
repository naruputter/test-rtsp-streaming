package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Camera represents a single CCTV camera configuration.
type Camera struct {
	ID      string `yaml:"id" json:"id"`
	Name    string `yaml:"name" json:"name"`
	RTSPURL string `yaml:"rtsp_url" json:"rtsp_url"`
	Enabled bool   `yaml:"enabled" json:"enabled"`
}

// Config holds the full application configuration.
type Config struct {
	ServerAddr   string   `yaml:"server_addr"`
	HLSOutputDir string   `yaml:"hls_output_dir"`
	Cameras      []Camera `yaml:"cameras"`
}

// Load reads the YAML config file from path and returns a Config.
// Sensible defaults are applied for missing fields.
func Load(path string) (*Config, error) {
	cfg := &Config{
		ServerAddr:   ":8080",
		HLSOutputDir: "./hls_output",
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
