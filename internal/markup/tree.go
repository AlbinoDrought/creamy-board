package markup

type Node struct {
	Quoted  bool
	Content string
}

func Parse(text string) []Node {
	nodes := []Node{}

	textLen := len(text)
	var (
		node           Node
		lastNodeOffset int
		newLine        bool
	)
	newLine = true

	commitNode := func(to int) {
		if lastNodeOffset == to {
			return
		}
		node.Content = text[lastNodeOffset:to]
		lastNodeOffset = to
		nodes = append(nodes, node)
		node = Node{}
	}

	for i := range text {
		if newLine {
			if text[i] == '>' {
				commitNode(i)
				node.Quoted = true
			}
		}

		newLine = text[i] == '\n'

		if newLine && node.Quoted {
			// end quote
			commitNode(i)
		}
	}
	commitNode(textLen)

	return nodes
}
