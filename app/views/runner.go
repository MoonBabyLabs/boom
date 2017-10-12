package views

import "github.com/MoonBabyLabs/kek/service"

type Runner interface {
	Run(cnt service.KekDoc, urlRoute string) Runner
}
