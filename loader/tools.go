package loader

type Tools struct {
	tools []Tool
}

func NewTools() *Tools {
	return &Tools{tools: make([]Tool, 0)}
}

func (t *Tools) Add(tool Tool) {
	t.tools = append(t.tools, tool)
}

func (t *Tools) All() []Tool {
	return t.tools
}

func (t *Tools) FindByName(name string) Tool {
	for _, tool := range t.tools {
		if tool.GetName() == name {
			return tool
		}
	}
	return nil
}
