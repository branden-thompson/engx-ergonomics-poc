package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	progress progress.Model
	value    float64
	done     bool
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.progress.Init(),
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tickMsg:
		if m.value >= 1.0 {
			m.done = true
			return m, tea.Quit
		}

		// Increment progress
		m.value += 0.02 // 2% per tick
		cmd := m.progress.SetPercent(m.value)
		return m, tea.Batch(cmd, tickCmd())

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return "\n✅ Complete! Press 'q' to quit.\n"
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06B6D4")).
		Bold(true).
		Render("🛩️ Testing Progress Animation")

	progressBar := m.progress.View()
	percent := fmt.Sprintf("%.1f%% complete", m.value*100)

	return fmt.Sprintf("\n%s\n\n%s\n%s\n\nPress 'q' to quit\n", title, progressBar, percent)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
		value:    0.0,
		done:     false,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}