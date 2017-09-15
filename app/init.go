package app

import (
	"github.com/revel/revel"
	"strings"
	"log"
	"encoding/json"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		ValidateDomainBasePath,	   	   // We don't want anyone trying to snake into the api without having the write path.
		RemoveDomainBasePath,	   	   // We've validated the path. Now lets set the needed api routes
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		ParseJsonBodyFilter,
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}


	// register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

func ParseJsonBodyFilter(c *revel.Controller, fc []revel.Filter) {
	content, err := json.Marshal(c.Request.Body)

	if err != nil {
		log.Print(err)
	}

	if content != nil && err != nil {
		c.Params.JSON = content
	}

	fc[0](c, fc[1:])
}

func ValidateDomainBasePath(c *revel.Controller, fc []revel.Filter) {
	log.Print(c.Request.RequestURI)
	log.Print(revel.Config.StringDefault("domain.base.path", ""))
	path := c.Request.RequestURI

	if c.Request.RequestURI == "/favicon.ico" {
		fc[0](c, fc[1:])
	}

	if !strings.Contains(path, revel.Config.StringDefault("domain.base.path", "")) {
		c.Request.Request.URL.Path = "/404"
	}

	fc[0](c, fc[1:])
}

// Lets remove base path from the app structure so its sitting on its own root.
func RemoveDomainBasePath(c *revel.Controller, fc []revel.Filter) {
	c.Request.Request.URL.Path = "/" + strings.Trim(c.Request.Request.URL.Path, revel.Config.StringDefault("domain.base.path", ""))

	fc[0](c, fc[1:])
}


// HeaderFilter adds common security headers
// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	accessHostsConfig := revel.Config.StringDefault("access.hosts", "")
	accessHosts := strings.Split(accessHostsConfig, ",")
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	if inStringArray(c.Request.Host, accessHosts) || accessHostsConfig == "*" {
		protocol := "http://"

		if strings.Contains(c.Request.Proto, "HTTPS") {
			protocol = "https://"
		}

		log.Print(c.Response.Out.Header().Get("Access-Control-Allow-Origin"))

		c.Response.Out.Header().Add(
			"Access-Control-Allow-Origin",
			protocol + c.Request.Host)

		if accessHostsConfig == "*" {
			c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
		}


		c.Response.Out.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
		c.Response.Out.Header().Add("Access-Control-Allow-Headers", "content-type,write,delete,update,publish")
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}


func inStringArray(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
