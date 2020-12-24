package types

type Properties interface {
	AddProperty(p *MDProperty)
}

type MDSchemasType struct {
	O                    *OpenApiType // Original object
	Name                 string
	Description          string
	Type                 string
	AllOff               bool
	AdditionalProperties bool
	Properties           []*MDProperty
}

func (o *MDSchemasType) AddProperty(p *MDProperty) {
	if o.Properties == nil {
		o.Properties = []*MDProperty{p}
		return
	}
	o.Properties = append(o.Properties, p)
}

type MDProperty struct {
	P           *OpenApiType // Original property
	Name        string
	Type        string
	Mandatory   string
	Description string
	Properties  []*MDProperty
}

func (o *MDProperty) AddProperty(p *MDProperty) {
	if o.Properties == nil {
		o.Properties = []*MDProperty{p}
		return
	}
	o.Properties = append(o.Properties, p)
}
