package types

type MDProperties interface {
	AddMDProperty(p *MDProperty)
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

func (o *MDSchemasType) AddMDProperty(p *MDProperty) {
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

func (o *MDProperty) AddMDProperty(p *MDProperty) {
	if o.Properties == nil {
		o.Properties = []*MDProperty{p}
		return
	}
	o.Properties = append(o.Properties, p)
}
