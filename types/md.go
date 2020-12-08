package types

type MDSchemasType struct {
	O                    *ObjectType // Original object
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
	P           *ObjectType // Original property
	Name        string
	Type        string
	Mandatory   string
	Description string
	SubElement  []*MDSchemasType
}

func (o *MDProperty) AddSubElement(md *MDSchemasType) {
	if o.SubElement == nil {
		o.SubElement = []*MDSchemasType{md}
		return
	}
	o.SubElement = append(o.SubElement, md)
}
