package payload

type Property struct {
	key, value string
}

func (p *Property) Key() string {
	return p.key
}

func (p *Property) Value() string {
	return p.value
}

func NewProperty(key, value string) Property {
	return Property{
		key:   key,
		value: value,
	}
}

type Properties []Property

type Element struct {
	id         string
	properties Properties
}

func (e *Element) Properties() Properties {
	return e.properties
}

func (e *Element) AddProperty(p Property) {
	e.properties = append(e.properties, p)
}

func (e *Element) ID() string {
	return e.id
}

func NewElement(id string) *Element {
	return &Element{
		id:         id,
		properties: make(Properties, 0, 1),
	}
}

type Data []Element
