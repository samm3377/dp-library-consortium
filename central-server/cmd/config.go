package main

import (
	"io/fs"
	"log"
	"os"
	"path"
	"strconv"
)

type Config struct {
	DbPath       string
	ServerPort   string
	SecretKey    string
	TemplateDir  string
	ReadTimeout  int64
	WriteTimeout int64
}

func LoadConfig() *Config {
	config := &Config{}

	// Configuration from .env
	config.DbPath = os.Getenv("DB_PATH")
	config.ServerPort = os.Getenv("CENTRAL_PORT")
	if config.ServerPort == "" {
		config.ServerPort = "8080" // Default if not set
	}

	config.SecretKey = os.Getenv("SESSION_SECRET")
	if config.SecretKey == "" {
		log.Fatal("SESSION_SECRET environment variable must be set")
	}

	config.TemplateDir = os.Getenv("TEMPLATE_DIR")
	if config.TemplateDir == "" {
		config.TemplateDir = "templates" // Default if not set
	}

	readtimeout, err := strconv.Atoi(os.Getenv("READTIMEOUT"))
	if err != nil || readtimeout == 0 {
		readtimeout = 10
	}

	config.ReadTimeout = int64(readtimeout)

	writetimeout, err := strconv.Atoi(os.Getenv("WRITETIMEOUT"))
	if err != nil || writetimeout == 0 {
		writetimeout = 10
	}

	config.WriteTimeout = int64(writetimeout)

	return config
}

func LookupTemplates(templateDirectory string) map[string][]string {
	mapping := make(map[string][]string)

	layoutPath := path.Join(templateDirectory, "layouts")
	layoutDir := os.DirFS(layoutPath)

	layoutFiles, err := fs.Glob(layoutDir, "*.html")

	if err != nil {
		log.Fatal("error while looking for layout files", err)
	}

	for i := range layoutFiles {
		layoutFiles[i] = path.Join(layoutPath, layoutFiles[i])
	}

	templateDir := os.DirFS(templateDirectory)
	templateFiles, err := fs.Glob(templateDir, "*.html")

	if err != nil {
		log.Fatal("error while looking for template files", err)
	}

	for _, tf := range templateFiles {
		mapping[path.Base(tf)] = append(layoutFiles, path.Join(templateDirectory, tf))
	}

	return mapping
}
