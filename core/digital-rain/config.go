// Package main implements a terminal-based Matrix digital rain effect visualization
package main

import (
	"time"
)

// Language represents the supported languages
type Language int

const (
	// English represents English language
	English Language = iota
	// Chinese represents Chinese language
	Chinese
)

// ToString converts Language to string
func (l Language) ToString() string {
	switch l {
	case Chinese:
		return "cn"
	default:
		return "en"
	}
}

// Default configuration values
const (
	DefaultRows            = 30
	DefaultCols            = 80
	DefaultRefreshRate     = 50 * time.Millisecond
	MinRefreshRate         = 10 * time.Millisecond
	DefaultLanguage        = English
	DefaultDropColor       = "#00FF00" // Matrix green
	DefaultTrailColor      = "#008800" // Darker green for trail
	DefaultBackgroundColor = "#000000" // Black background
	DefaultMinSpeed        = 1
	DefaultMaxSpeed        = 5
	DefaultDropLength      = 10
	DefaultCharSet         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	DefaultLogFile         = ""
	DefaultProfilePort     = 6061
	DefaultProfileInterval = 5 * time.Second
)

// Config holds the configuration for the digital rain
type Config struct {
	DropColor       string
	TrailColor      string
	BackgroundColor string
	CharSet         string
	MinSpeed        int
	MaxSpeed        int
	DropLength      int
	Language        Language
}

// SetLanguage sets the language from a string
func (c *Config) SetLanguage(lang string) {
	switch lang {
	case "cn", "chinese", "zh":
		c.Language = Chinese
	default:
		c.Language = English
	}
}

// Check validates and fixes the configuration
func (c *Config) Check() {
	if c.DropColor == "" {
		c.DropColor = DefaultDropColor
	}
	if c.TrailColor == "" {
		c.TrailColor = DefaultTrailColor
	}
	if c.BackgroundColor == "" {
		c.BackgroundColor = DefaultBackgroundColor
	}
	if c.CharSet == "" {
		c.CharSet = DefaultCharSet
	}
	if c.MinSpeed < 1 {
		c.MinSpeed = DefaultMinSpeed
	}
	if c.MaxSpeed < c.MinSpeed {
		c.MaxSpeed = DefaultMaxSpeed
	}
	if c.DropLength < 1 {
		c.DropLength = DefaultDropLength
	}
}
