// Package assets embeds the six graphical, font and music resources recovered
// from the supplied Windows production.
package assets

import "embed"

// Files holds the original artwork, font data and YM soundtrack.
//
//go:embed original/*
var Files embed.FS
