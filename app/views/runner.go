package views

import "github.com/MoonBabyLabs/kek"

type Runner interface {
	Run(cnt kek.KekDoc, urlRoute string) Runner
}
