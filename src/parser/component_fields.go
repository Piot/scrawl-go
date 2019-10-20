package parser

import (
	"fmt"

	"github.com/piot/scrawl-go/src/definition"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func MakeComponentFields(root *definition.Root, fields []*definition.Field, validComponentTypes []string) ([]*definition.ComponentField, error) {
	var componentFields []*definition.ComponentField
	for _, fieldComponent := range fields {
		componentReference := root.FindComponent(fieldComponent.FieldType())
		var componentType *definition.ComponentDataType
		if componentReference == nil {
			rawType := fieldComponent.FieldType()
			if !Contains(validComponentTypes, rawType) {
				return nil, fmt.Errorf("unknown component type:%v", fieldComponent.FieldType())
			}
			fieldComponent.ForceNewIndex(0)
			componentType = definition.NewComponentTypeUsingField(fieldComponent)
		} else {
			componentType = definition.NewComponentTypeUsingComponent(componentReference)
		}
		componentField := definition.NewComponentField(len(componentFields), fieldComponent.Name(), componentType)
		componentFields = append(componentFields, componentField)
	}

	return componentFields, nil
}
