package pongo2

// Doc = { ( Filter | Tag | HTML ) }
func (p *Parser) parseDocElement() (INode, *Error) {
	t := p.Current()

	switch t.Typ {
	case TokenHTML:
		if p.template.Options.TrimWhitespace {
			t.Val, p.htmlInQuote, p.htmlLastChar = stripWhitespace(t.Val, p.htmlInQuote, p.htmlLastChar)
		}
		n := &nodeHTML{token: t}
		left := p.PeekTypeN(-1, TokenSymbol)
		right := p.PeekTypeN(1, TokenSymbol)
		n.trimLeft = left != nil && left.TrimWhitespaces
		n.trimRight = right != nil && right.TrimWhitespaces
		p.Consume() // consume HTML element
		return n, nil
	case TokenSymbol:
		switch t.Val {
		case "{{":
			// parse variable
			if p.template.Options.TrimWhitespace {
				// We don't reset lastChar here anymore, as we want to preserve context across variables
				// But we might want to ensure we don't accidentally strip space if the variable is the first thing?
				// Actually, stripWhitespace handles the logic.
				// If we had `foo {{`, stripWhitespace for `foo ` would have kept the space (lastChar='o').
				// If we had `> {{`, stripWhitespace for `> ` would have stripped it (lastChar='>').
				// So we don't need to do anything here.
			}
			variable, err := p.parseVariableElement()
			if err != nil {
				return nil, err
			}
			return variable, nil
		case "{%":
			// parse tag
			tag, err := p.parseTagElement()
			if err != nil {
				return nil, err
			}
			return tag, nil
		}
	}
	return nil, p.Error("Unexpected token (only HTML/tags/filters in templates allowed)", t)
}

func (tpl *Template) parse() *Error {
	tpl.parser = newParser(tpl.name, tpl.tokens, tpl)
	doc, err := tpl.parser.parseDocument()
	if err != nil {
		return err
	}
	tpl.root = doc
	return nil
}

func (p *Parser) parseDocument() (*nodeDocument, *Error) {
	doc := &nodeDocument{}

	for p.Remaining() > 0 {
		node, err := p.parseDocElement()
		if err != nil {
			return nil, err
		}
		doc.Nodes = append(doc.Nodes, node)
	}

	return doc, nil
}
