package provider

import "reflect"

type Contract interface {
	Construct() Contract
}

type ProviderList map[string]Contract

type Base struct {
	Providers ProviderList
}



func (b Base) Get (provider Contract) interface{} {
	name := reflect.TypeOf(provider).Name()

	if b.Providers == nil {
		b.Set(provider, b.Providers)
	}

	return b.Providers[name]
}

func (b Base) Set (provider Contract, list ProviderList) ProviderList {
	name := reflect.TypeOf(provider).Name()
	list[name] = provider.Construct()

	return list
}