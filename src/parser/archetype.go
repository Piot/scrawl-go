package parser

import "github.com/piot/scrawl-go/src/definition"

func (p *Parser) parseGenericArchetype() (*definition.EntityArchetype, error) {
	name, meta, nameErr := p.parseArchetypeNameAndStartScope()
	if nameErr != nil {
		return nil, nameErr
	}

	entityArchetypeItems, entityArchetypeItemsErr := p.parseEntityArchetypeItemsUntilEndScope()

	if entityArchetypeItemsErr != nil {
		return nil, entityArchetypeItemsErr
	}

	mainLod := definition.NewEntityArchetypeLOD(0, entityArchetypeItems)
	archetype := definition.NewEntityArchetype(name, definition.NewEntityIndex(0xff),
		[]*definition.EntityArchetypeLOD{mainLod}, meta)

	return archetype, nil
}
