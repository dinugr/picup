package config

import (
	"path/filepath"
	"picup/core/domain/schema"
)

var Server schema.Server = schema.Server{}
var Engine schema.EngineMap = schema.EngineMap{}
var Variant schema.VariantMap = schema.VariantMap{}

func GetDirMap() map[string]string {

	dirs := map[string]string{}
	dirs["temp_dir"] = filepath.Join(Server.TempDir)
	dirs["data_dir"] = filepath.Join(Server.DataDir)
	for _, v := range Variant {
		dirs[v.Key] = filepath.Join(Server.DataDir, v.Name)
	}

	return dirs
}
