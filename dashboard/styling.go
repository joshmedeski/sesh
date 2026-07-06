package dashboard

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// GroupNameRender returns a lipgloss style for the given group name
func GroupNameRender(name string, width int) lipgloss.Style {
	return NewStyle(width, width, 1, 1, 15, false, []int{0, 0, 0, 0})
}

// Session line styling
func SessionLineRender(num, name, branch, gitStatus, meta string, totalWidth int) string {

	// Calculate proportional column widths (Let Lipgloss handle truncation via .Width().MaxWidth())
	overhead := 8
	colSpace := max(totalWidth-overhead, 1)

	nameW := max(int(float64(colSpace)*0.35), 5)
	branchW := max(int(float64(colSpace)*0.30), 0)
	gitW := max(int(float64(colSpace)*0.30), 0)
	metaW := max(int(float64(colSpace)*0.05), 5)

	// Use Lipgloss layout blocks for columns to bypass manual loop padding math
	numPart := NewStyle(5, 5, 1, 1, 7, true, []int{0, 0, 0, 0}).Render(num)
	namePart := NewStyle(nameW, nameW, 1, 1, 15, false, []int{0, 0, 0, 2}).Render(name)
	branchPart := NewStyle(branchW, branchW, 1, 1, 5, false, []int{0, 0, 0, 0}).Render(branch)
	gitPart := NewStyle(gitW, gitW, 1, 1, 10, false, []int{0, 0, 0, 0}).Render(gitStatus)
	metaPart := NewStyle(metaW, metaW, 1, 1, 15, false, []int{0, 0, 0, 0}).Render(meta)

	// TODO: make the colors configurable/use the current terminal theme

	leftSide := lipgloss.JoinHorizontal(lipgloss.Top, numPart, namePart)
	rightSide := lipgloss.JoinHorizontal(lipgloss.Top, branchPart, gitPart, metaPart)

	// Calculate space between left and right sides
	leftWidth := lipgloss.Width(leftSide)
	rightWidth := lipgloss.Width(rightSide)

	gapWidth := max(totalWidth-leftWidth-rightWidth, 0)

	spaceGap := strings.Repeat(" ", gapWidth)

	return lipgloss.JoinHorizontal(lipgloss.Center, leftSide, spaceGap, rightSide)
}
