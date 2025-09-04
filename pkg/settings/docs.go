// Package settings provides configuration management for the Omarchy Theme Generator
// using Viper for flexible, layered configuration support.
//
// Configuration Sources and Precedence:
//  1. Built-in defaults
//  2. System config: /etc/omarchy/omarchy-theme-gen.json
//  3. User config: $XDG_CONFIG_HOME/omarchy/omarchy-theme-gen.json
//  4. Workspace config: ./omarchy-theme-gen.json
//  5. Environment variables: OMARCHY_THEME_GEN_*
//
// Usage:
//
//	settings, err := settings.Load()
//	ctx := settings.WithSettings(context.Background(), settings)
//	s := settings.FromContext(ctx)
package settings
