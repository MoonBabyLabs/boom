package provider


type Contract interface {
	Construct() Contract
	GetName() string
}

type ProviderList map[string]Contract

type Base struct {
	Providers ProviderList
}



func (b Base) Get (provider Contract) interface{} {
	if b.Providers == nil {
		b.Set(provider, b.Providers)
	}

	return b.Providers[provider.GetName()]
}

func (b Base) Set (provider Contract, list ProviderList) ProviderList {
	list[provider.GetName()] = provider.Construct()

	return list
}