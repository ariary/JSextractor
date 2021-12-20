package config

// Config holds the gitar configuration
type Config struct {
	Url       string
	GatherSrc bool
	SkipSrc   bool
	SkipEvent bool
	SkipTag   bool
}
