package views

import "github.com/MoonBabyLabs/kek"

type Runner interface {
	Run(cnt kek.Doc, urlRoute string) Runner
}
