package parser

import (
	"fmt"

	"github.com/piot/scrawl-go/src/definition"
)

func MakeComponentFields(root *definition.Root, fields []*definition.Field) ([]*definition.ComponentField, error) {
	var componentFields []*definition.ComponentField
	for _, fieldComponent := range fields {
		componentType := root.FindComponent(fieldComponent.FieldType())
		if componentType == nil {
			return nil, fmt.Errorf("unknown component type:%v", fieldComponent.FieldType())
		}
		componentField := definition.NewComponentField(len(componentFields), fieldComponent.Name(), componentType)
		componentFields = append(componentFields, componentField)
	}

	return componentFields, nil
}
