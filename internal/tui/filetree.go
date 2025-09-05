package tui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Node represents a file or directory within the file tree structure.
type Node struct {
	name      string
	status    string // Git status prefix (e.g., "M ", "MM", "??"), only for file nodes.
	path      string // Full path relative to the repo root
	isRenamed bool   // Flag to indicate a renamed/copied file
	children  []*Node
}

// BuildTree parses the output of `git status --porcelain` to construct a file tree.
func BuildTree(gitStatus string) *Node {
	root := &Node{name: "."}
	lines := strings.Split(strings.TrimSpace(gitStatus), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return root // No changes.
	}

	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status := line[:2]
		path := strings.TrimSpace(line[3:])
		isRenamed := false

		if status[0] == 'R' || status[0] == 'C' {
			parts := strings.Split(path, " -> ")
			if len(parts) == 2 {
				path = parts[1]
				isRenamed = true
			}
		}

		parts := strings.Split(path, string(filepath.Separator))
		currentNode := root
		for i, part := range parts {
			childNode := currentNode.findChild(part)
			if childNode == nil {
				childNode = &Node{name: part}
				currentNode.children = append(currentNode.children, childNode)
			}
			currentNode = childNode

			if i == len(parts)-1 {
				currentNode.status = status
				currentNode.path = path
				currentNode.isRenamed = isRenamed
			}
		}
	}

	root.sort()
	root.compact()
	return root
}

// findChild searches for an immediate child node by name.
func (n *Node) findChild(name string) *Node {
	for _, child := range n.children {
		if child.name == name {
			return child
		}
	}
	return nil
}

// sort recursively sorts the children of a node.
func (n *Node) sort() {
	if n.children == nil {
		return
	}
	sort.SliceStable(n.children, func(i, j int) bool {
		isDirI := len(n.children[i].children) > 0
		isDirJ := len(n.children[j].children) > 0
		if isDirI != isDirJ {
			return isDirI // Directories first.
		}
		return n.children[i].name < n.children[j].name
	})

	for _, child := range n.children {
		child.sort()
	}
}

// compact recursively merges directories that contain only a single sub-directory.
func (n *Node) compact() {
	if n.children == nil {
		return
	}
	for _, child := range n.children {
		child.compact()
	}
	if n.name == "." {
		return
	}
	for len(n.children) == 1 && len(n.children[0].children) > 0 {
		child := n.children[0]
		n.name = filepath.Join(n.name, child.name)
		n.children = child.children
	}
}

// Render traverses the tree and returns a slice of strings for display.
func (n *Node) Render(theme Theme) []string {
	return n.renderRecursive("", theme)
}

// renderRecursive creates raw, tab-delimited strings for the view to parse.
func (n *Node) renderRecursive(prefix string, theme Theme) []string {
	var lines []string
	for i, child := range n.children {
		connector := theme.Tree.Connector
		newPrefix := theme.Tree.Prefix
		if i == len(n.children)-1 {
			connector = theme.Tree.ConnectorLast
			newPrefix = theme.Tree.PrefixLast
		}

		if len(child.children) > 0 { // It's a directory
			// Format: "prefix\tconnector\tname"
			lines = append(lines, fmt.Sprintf("%s%s▼\t\t%s", prefix, connector, child.name))
			lines = append(lines, child.renderRecursive(prefix+newPrefix, theme)...)
		} else { // It's a file
			displayName := child.name
			if child.isRenamed {
				displayName = child.path
			}
			// Format: "prefix\tconnector\tstatus\tname"
			lines = append(lines, fmt.Sprintf("%s%s\t%s\t%s", prefix, connector, child.status, displayName))
		}
	}
	return lines
}
