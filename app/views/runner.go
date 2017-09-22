package views

type Runner interface {
	Run(cnt map[string]interface{}, fields []map[string]map[string]interface{}, urlRoute string, cType string) Runner
}
