package controllers

import "github.com/revel/revel"

type Options struct {
	*revel.Controller
	Base
}

// Options provides a route request for OPTIONS based routes.
// It is generally useful for preflight requests from browsers.
// Returns a success message in the revel render result.
func (c Options) Options() revel.Result {
	success := make(map[string]bool)
	success["success"] = true

	return c.RenderJSON(success)
}
