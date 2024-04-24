package main

type ListItem struct {
	title string
}

// implement the list.Item interface
func (l ListItem) FilterValue() string {
	return l.title
}

func (l ListItem) Title() string {
	return l.title
}
