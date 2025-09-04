package tui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Node represents a file or directory within the file tree structure.
type Node struct {
	name     string
	status   string // Git status prefix (e.g., "M", "??"), only for file nodes.
	children []*Node
}

// BuildTree parses the output of `git status --porcelain` to construct a file
// tree. It processes each line, builds a hierarchical structure of nodes,
// sorts them, and compacts single-child directories for a cleaner display.
func BuildTree(gitStatus string) *Node {
	root := &Node{name: "."}

	lines := strings.Split(strings.TrimSpace(gitStatus), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return root // No changes, return the root.
	}

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		spaceIndex := strings.Index(line, " ")
		if spaceIndex == -1 {
			continue
		}

		status := strings.TrimSpace(line[:spaceIndex])
		path := line[spaceIndex+1:]

		parts := strings.Split(path, string(filepath.Separator))
		currentNode := root
		for i, part := range parts {
			// Traverse the tree, creating nodes as necessary.
			childNode := currentNode.findChild(part)
			if childNode == nil {
				childNode = &Node{name: part}
				currentNode.children = append(currentNode.children, childNode)
			}
			currentNode = childNode

			// The last part of the path is the file, so set its status.
			if i == len(parts)-1 {
				currentNode.status = status
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

// sort recursively sorts the children of a node. Directories are listed first,
// then files, with both groups sorted alphabetically.
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
		return n.children[i].name < n.children[j].name // Then sort alphabetically.
	})

	for _, child := range n.children {
		child.sort()
	}
}

// compact recursively merges directories that contain only a single sub-directory.
// For example, a path like "src/main/go" becomes a single node.
func (n *Node) compact() {
	if n.children == nil {
		return
	}

	// First, compact all children in a post-order traversal.
	for _, child := range n.children {
		child.compact()
	}

	// Merge this node with its child if it's a single-directory container,
	// but do not compact the root node itself.
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
func (n *Node) Render() []string {
	return n.renderRecursive("")
}

// renderRecursive performs a depth-first traversal to generate the visual
// representation of the tree, using box-drawing characters to show hierarchy.
func (n *Node) renderRecursive(prefix string) []string {
	var lines []string
	for i, child := range n.children {
		// Use different connectors for the last child in a list.
		connector := "├─"
		newPrefix := "│  "
		if i == len(n.children)-1 {
			connector = "└─"
			newPrefix = "   "
		}

		if len(child.children) > 0 {
			// It's a directory.
			lines = append(lines, fmt.Sprintf("%s%s▼ %s", prefix, connector, child.name))
			lines = append(lines, child.renderRecursive(prefix+newPrefix)...)
		} else {
			// It's a file.
			lines = append(lines, fmt.Sprintf("%s%s %s %s", prefix, connector, child.status, child.name))
		}
	}
	return lines
}
