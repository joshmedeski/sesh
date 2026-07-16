// styling.go
package dashboard

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func GroupNameRender(name string, width int) lipgloss.Style {
	return NewStyle(width, width, 1, 1, 15, false, []int{0, 0, 0, 0})
}

func SessionLineRender(num, name, branch, gitStatus, meta string, totalWidth int) string {
	overhead := 8
	colSpace := max(totalWidth-overhead, 1)

	nameW := max(int(float64(colSpace)*0.35), 5)
	branchW := max(int(float64(colSpace)*0.30), 0)
	gitW := max(int(float64(colSpace)*0.30), 0)
	metaW := max(int(float64(colSpace)*0.05), 5)

	// Fix: Apply Height/MaxHeight rules directly onto strings via lipgloss styles to restrict output to 1 line
	numPart := lipgloss.NewStyle().Width(5).MaxWidth(5).Height(1).MaxHeight(1).Foreground(lipgloss.ANSIColor(7)).Render(num)
	namePart := lipgloss.NewStyle().Width(nameW).MaxWidth(nameW).Height(1).MaxHeight(1).Foreground(lipgloss.ANSIColor(15)).Render(name)
	branchPart := lipgloss.NewStyle().Width(branchW).MaxWidth(branchW).Height(1).MaxHeight(1).Foreground(lipgloss.ANSIColor(5)).Render(branch)
	gitPart := lipgloss.NewStyle().Width(gitW).MaxWidth(gitW).Height(1).MaxHeight(1).Foreground(lipgloss.ANSIColor(10)).Render(gitStatus)
	metaPart := lipgloss.NewStyle().Width(metaW).MaxWidth(metaW).Height(1).MaxHeight(1).Foreground(lipgloss.ANSIColor(15)).Render(meta)

	leftSide := lipgloss.JoinHorizontal(lipgloss.Top, numPart, namePart)
	rightSide := lipgloss.JoinHorizontal(lipgloss.Top, branchPart, gitPart, metaPart)

	leftWidth := lipgloss.Width(leftSide)
	rightWidth := lipgloss.Width(rightSide)
	gapWidth := max(totalWidth-leftWidth-rightWidth, 0)
	spaceGap := strings.Repeat(" ", gapWidth)

	return lipgloss.JoinHorizontal(lipgloss.Center, leftSide, spaceGap, rightSide)
}
