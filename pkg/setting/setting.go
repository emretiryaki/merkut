package setting

import (
	"gopkg.in/ini.v1"
	"time"
)

type Cfg struct {

	Raw *ini.File
	// Paths
	ProvisioningPath string

	// SMTP email settings
	Smtp SmtpSettings

	// Rendering
	ImagesDir                        string
	PhantomDir                       string
	RendererUrl                      string
	DisableBruteForceLoginProtection bool

	TempDataLifetime time.Duration
}

func NewCfg() *Cfg {
	return &Cfg{
		Raw: ini.Empty(),
	}
}