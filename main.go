package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
)

var words = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
	"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
	"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
	"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
	"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	"is", "am", "are", "was", "were", "been", "being", "did", "had", "has",
	"having", "may", "might", "must", "shall", "should", "can", "could", "will", "would",
	"do", "does", "did", "doing", "say", "says", "said", "go", "goes", "went",
	"gone", "see", "sees", "saw", "seen", "know", "knows", "knew", "known", "think",
	"thinks", "thought", "get", "gets", "got", "gotten", "make", "makes", "made", "want",
	"wants", "wanted", "give", "gives", "gave", "given", "use", "uses", "used", "find",
	"finds", "found", "tell", "tells", "told", "ask", "asks", "asked", "work", "works",
	"worked", "seem", "seems", "seemed", "feel", "feels", "felt", "try", "tries", "tried",
	"leave", "leaves", "left", "call", "calls", "called", "good", "new", "first", "last",
	"long", "great", "little", "own", "other", "old", "right", "big", "high", "different",
	"small", "large", "next", "early", "young", "important", "few", "public", "bad", "same",
	"able",
}

var wordCount int = 50

type model struct {
	prompt          string
	userText        string
	startTime       time.Time
	lastPress       time.Time
	totalKeystrokes int
	width           int
}

func initialModel() model {
	prompt := ""
	for range wordCount {
		prompt += words[rand.Intn(len(words))] + " "
	}

	return model{
		prompt:   prompt,
		userText: "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			newPrompt := ""
			for range wordCount {
				newPrompt += words[rand.Intn(len(words))] + " "
			}
			m.prompt = newPrompt
			m.userText = ""
			m.totalKeystrokes = 0
			m.startTime = time.Time{}
			m.lastPress = time.Time{}
		case "backspace":
			if len(m.userText) > 0 && len(m.userText) < len(m.prompt) {
				m.userText = m.userText[:len(m.userText)-1]
			}
		default:
			if len(msg.String()) == 1 {
				if len(m.userText) < len(m.prompt) {
					if m.totalKeystrokes == 0 {
						m.startTime = time.Now()
					}
					m.lastPress = time.Now()
					m.totalKeystrokes++
					m.userText += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var correctKeystrokes int
	var styledPrompt strings.Builder
	for i, char := range m.prompt {
		if i < len(m.userText) {
			if m.userText[i] == byte(char) {
				styledPrompt.WriteString(green.Render(string(char)))
				correctKeystrokes++
			} else {
				styledPrompt.WriteString(red.Render(string(char)))
			}
		} else {
			styledPrompt.WriteString(string(char))
		}
	}

	elapsed := 0.0
	if !m.startTime.IsZero() {
		elapsed = m.lastPress.Sub(m.startTime).Seconds()
	}

	wpm := 0.0
	if elapsed > 0 {
		minutes := elapsed / 60.0
		wpm = (float64(correctKeystrokes) / 5.0) / minutes
	}

	accPct := 0.0
	if m.totalKeystrokes > 0 {
		accPct = float64(correctKeystrokes) / float64(m.totalKeystrokes) * 100.0
	}

	cpm := 0.0
	if elapsed > 0 {
		cpm = float64(correctKeystrokes) / elapsed * 60.0
	}

	mistakes := m.totalKeystrokes - correctKeystrokes

	stats := fmt.Sprintf(
		"Time: %.0fs | WPM: %.1f | Accuracy: %.1f%% | CPM: %.1f | Mistakes: %d",
		elapsed,
		wpm,
		accPct,
		cpm,
		mistakes,
	)

	promptStyle := lipgloss.NewStyle().
		Width(m.width-4).
		Border(lipgloss.RoundedBorder(), true).
		Align(lipgloss.Center).
		Padding(1).
		Margin(1)

	promptBox := promptStyle.Render(styledPrompt.String())

	basicStyle := lipgloss.NewStyle().
		Width(m.width - 2).
		Align(lipgloss.Center).
		Margin(1)

	statsStyled := basicStyle.Render(stats)
	infoStyled := basicStyle.Render("ctrl+r to restart, ctrl+c to quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		promptBox,
		statsStyled,
		infoStyled,
		"\n",
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occurred: %v", err)
		os.Exit(1)
	}
}

