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
	status    string // Git status prefix (e.g., "M ", "??"), only for file nodes.
	path      string // Full path relative to the repository root.
	isRenamed bool
	children  []*Node
}

// BuildTree parses the output of `git status --porcelain` to construct a file tree.
func BuildTree(gitStatus string) *Node {
	root := &Node{name: repoRootNodeName, path: "."}

	lines := strings.Split(gitStatus, "\n")
	if len(lines) == 1 && lines[0] == "" {
		return root
	}

	for _, line := range lines {
		if len(line) < porcelainStatusPrefixLength {
			continue
		}
		status := line[:2]
		fullPath := strings.TrimSpace(line[porcelainStatusPrefixLength:])
		isRenamed := false

		if status[0] == 'R' || status[0] == 'C' {
			parts := strings.Split(fullPath, gitRenameDelimiter)
			if len(parts) == 2 {
				fullPath = parts[1]
				isRenamed = true
			}
		}

		parts := strings.Split(fullPath, string(filepath.Separator))
		currentNode := root
		for i, part := range parts {
			childNode := currentNode.findChild(part)
			if childNode == nil {
				// Construct path for the new node based on its parent
				nodePath := filepath.Join(currentNode.path, part)
				if currentNode.path == "." {
					nodePath = part
				}
				childNode = &Node{name: part, path: nodePath}
				currentNode.children = append(currentNode.children, childNode)
			}
			currentNode = childNode

			if i == len(parts)-1 { // Leaf node (file)
				currentNode.status = status
				currentNode.path = fullPath // Overwrite with the full path from git
				currentNode.isRenamed = isRenamed
			}
		}
	}

	root.sort()
	root.compact()
	return root
}

// Render traverses the tree and returns a slice of formatted strings for display.
func (n *Node) Render(theme Theme) []string {
	return n.renderRecursive("", theme)
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

// sort recursively sorts the children of a node, placing directories before files.
func (n *Node) sort() {
	if n.children == nil {
		return
	}
	sort.SliceStable(n.children, func(i, j int) bool {
		isDirI := len(n.children[i].children) > 0
		isDirJ := len(n.children[j].children) > 0
		if isDirI != isDirJ {
			return isDirI
		}
		return n.children[i].name < n.children[j].name
	})

	for _, child := range n.children {
		child.sort()
	}
}

// compact recursively merges directories that contain only a single sub-directory
// to create a more concise file tree.
func (n *Node) compact() {
	if n.children == nil {
		return
	}

	// Recursively compact children first.
	for _, child := range n.children {
		child.compact()
	}

	// Do not compact the root node itself.
	if n.name == repoRootNodeName {
		return
	}

	// If a directory has only one child and that child is also a directory, merge them.
	for len(n.children) == 1 && len(n.children[0].children) > 0 {
		child := n.children[0]
		n.name = filepath.Join(n.name, child.name)
		n.path = child.path
		n.children = child.children
	}
}

// renderRecursive performs a depth-first traversal of the tree to generate
// raw, tab-delimited strings for the view to parse and style.
func (n *Node) renderRecursive(prefix string, theme Theme) []string {
	var lines []string
	for _, child := range n.children {
		newPrefix := prefix + theme.Tree.Prefix

		if len(child.children) > 0 { // It's a directory
			displayName := dirExpandedIcon + child.name
			lines = append(lines, fmt.Sprintf("%s\t\t%s\t%s", prefix, displayName, child.path))
			lines = append(lines, child.renderRecursive(newPrefix, theme)...)
		} else { // It's a file.
			displayName := child.name
			if child.isRenamed {
				displayName = child.path
			}
			lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s", prefix, child.status, displayName, child.path))
		}
	}
	return lines
}
