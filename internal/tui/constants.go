package tui

import "time"

const (
	// --- Layout Ratios ---
	// leftPanelWidthRatio defines the percentage of the total width for the left column.
	// rightPanelWidthRatio is "1 - leftPanelWidthRatio".
	leftPanelWidthRatio = 0.35
	// helpViewWidthRatio defines the width of the help view relative to the total width.
	helpViewWidthRatio = 0.5
	// helpViewHeightRatio defines the height of the help view relative to the total height.
	helpViewHeightRatio = 0.75
	// expandedPanelHeightRatio defines the height of a focused panel relative to the total height.
	expandedPanelHeightRatio = 0.4

	// --- Layout Dimensions ---
	// collapsedPanelHeight is the fixed height for a panel when it's not in focus.
	collapsedPanelHeight = 3
	// borderWidth is the horizontal space taken by left and right borders.
	borderWidth = 2
	// titleBarHeight is the vertical space taken by top and bottom borders with titles.
	titleBarHeight = 2
	// statusPanelHeight is the fixed height for the status panel.
	statusPanelHeight = 3

	// --- Help View Styling ---
	// helpTitleMargin is the left margin for the title in the help view.
	helpTitleMargin = 9
	// helpKeyWidth is the fixed width for the keybinding column in the help view.
	helpKeyWidth = 12
	// helpDescMargin is the right margin for the keybinding column in the help view.
	helpDescMargin = 1

	// --- Characters & Symbols ---
	scrollThumbChar       = "▐"
	graphNodeChar         = "○"
	dirExpandedIcon       = "▼ "
	repoRootNodeName      = "."
	gitRenameDelimiter    = " -> "
	initialContentLoading = "Loading..."

	// --- File Watcher ---
	// fileWatcherPollInterval is the debounce interval for repository file system events.
	fileWatcherPollInterval = 500 * time.Millisecond

	// --- Git Status Parsing ---
	// porcelainStatusPrefixLength is the length of the status prefix in `git status --porcelain`.
	porcelainStatusPrefixLength = 3
)

// --- Border Characters ---
const (
	borderTop         = "─"
	borderBottom      = "─"
	borderLeft        = "│"
	borderRight       = "│"
	borderTopLeft     = "╭"
	borderTopRight    = "╮"
	borderBottomLeft  = "╰"
	borderBottomRight = "╯"
)

// --- Tree Characters ---
const (
	treeConnector     = ""
	treeConnectorLast = ""
	treePrefix        = "    "
	treePrefixLast    = "   "
)
