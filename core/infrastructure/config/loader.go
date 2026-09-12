package config

import (
	"errors"
	"fmt"
	"picup/core/domain/constants"
	"picup/core/domain/schema"
	"picup/core/domain/variantargs"
	"picup/core/infrastructure/security"
	"strings"

	"gopkg.in/ini.v1"
)

func Load(path string) error {
	cfg, err := ini.ShadowLoad(path)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	s := cfg.Section("server")
	Server.Host = s.Key("host").MustString("localhost")
	Server.Port = s.Key("port").MustUint64(9906)

	Server.DataDir = s.Key("data_dir").MustString("./data")
	Server.TempDir = s.Key("temp_dir").MustString("./temp")
	Server.MaxSize = s.Key("max_size").MustInt64(10485760)
	Server.FileTypes = s.Key("file_types").Strings(",")
	Server.ExifToolPath = s.Key("exiftool_path").MustString("exiftool")
	Server.VariantsSequence = s.Key("variants_sequence").Strings(",")

	// URL path prefix under which image variants are served.
	// Keep it without trailing slash, e.g. "/assets/images".
	Server.AssetsBasePath = s.Key("assets_base_path").MustString("/assets/images")

	// data source name
	Server.DSN = s.Key("dsn").MustString("sqlite://./picup.db")

	// webui
	Server.WebUIDir = s.Key("webui_dir").MustString("./webui")
	Server.WebUIEnabled = s.Key("webui_enabled").MustBool(false)

	// auth
	defaultJWTSecret := security.NewSecret(32)
	Server.Auth.JWTSecret = s.Key("jwt_secret").MustString(defaultJWTSecret)
	Server.Auth.JWTIssuer = s.Key("jwt_issuer").MustString("picup")
	Server.Auth.JWTAudience = s.Key("jwt_audience").MustString("picup")
	Server.Auth.JWTExpirationSeconds = s.Key("jwt_expiration_seconds").MustInt64(3600)
	Server.Auth.JWTSecretAutogen = (defaultJWTSecret == Server.Auth.JWTSecret)

	Server.Auth.AdminUsername = s.Key("admin_username").MustString("admin")
	Server.Auth.AdminPassword = s.Key("admin_password").MustString("admin")

	// variant and engine populations
	for _, section := range cfg.Sections() {
		if kname, ok := strings.CutPrefix(section.Name(), "variant."); ok {
			variant := &schema.Variant{
				Key:     kname,
				Name:    section.Key("name").MustString(kname),
				Type:    section.Key("type").String(),
				Engine:  section.Key("engine").MustString(constants.RESERVED_ENGINE_DEFAULT_KEY),
				RawArgs: section.Key("args").String(),
			}

			Variant[kname] = variant
			continue
		}

		if kname, ok := strings.CutPrefix(section.Name(), "engine."); ok {
			engine := &schema.Engine{
				Key:        kname,
				Name:       section.Key("name").MustString(kname),
				Provider:   section.Key("provider").String(),
				Parameters: map[string]string{},
			}
			for _, key := range cfg.Section(section.Name()).Keys() {
				engine.Parameters[key.Name()] = key.String()
			}

			Engine[engine.Name] = engine
			continue
		}
	}

	if errs := validate(); len(errs) > 0 {
		return fmt.Errorf("Configuration error: %s", strings.Join(errs, "; "))
	}

	registerReservedEngine(Engine)
	registerReservedVariant(Variant)

	if err := CompileVariants(Variant); err != nil {
		return fmt.Errorf("variant argument compilation failed: %w", err)
	}

	return nil

}

// CompileVariants parses raw variant arguments once during configuration loading.
func CompileVariants(variants schema.VariantMap) error {
	var errs []error

	for name, current := range variants {
		if current == nil {
			errs = append(errs, fmt.Errorf("variant %q: configuration is nil", name))
			continue
		}

		spec, err := variantargs.Parse(current.RawArgs)
		if err != nil {
			errs = append(errs, fmt.Errorf("variant %q: %w", name, err))
			continue
		}
		current.Pipeline = spec
	}

	return errors.Join(errs...)
}
