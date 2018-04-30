// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {

	Period time.Duration 			`config:"period"`
	RegistryFile     string 		`config:"registry_file"`
	Url string  					`config:"url"`
	Database string					`config:"database"`
	Collection string 				`config:"collection"`

}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	RegistryFile: "registry",
	Url: "localhost",
	Database: "elmah",
	Collection: "Elmah",
}
