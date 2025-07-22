package main

import (
	"fmt"
	"math/rand"
	"os"
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

type model struct {
	prompt            string
	userText          []string
	startTime         time.Time
	lastPress         time.Time
	totalKeystrokes   int
	correctKeystrokes int
}

func initialModel() model {
	prompt := ""
	for range 20 {
		prompt += words[rand.Intn(len(words))] + " "
	}

	return model{
		prompt:   prompt,
		userText: []string{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			newPrompt := ""
			for range 20 {
				newPrompt += words[rand.Intn(len(words))] + " "
			}
			m.prompt = newPrompt
			m.userText = []string{}
			m.correctKeystrokes = 0
			m.totalKeystrokes = 0
		case "backspace":
			if len(m.prompt) <= len(m.userText) {
				return m, nil
			}
			if len(m.userText) > 0 {
				m.userText = m.userText[:len(m.userText)-1]
			}
		default:
			if len(m.prompt) <= len(m.userText) {
				return m, nil
			}
			m.lastPress = time.Now()
			if m.totalKeystrokes == 0 {
				m.startTime = time.Now()
			}
			m.totalKeystrokes++
			text := msg.String()
			if len(m.userText) < len(m.prompt) && m.prompt[len(m.userText)] == text[0] {
				m.userText = append(m.userText, green.Render(text))
				m.correctKeystrokes++
			} else {
				m.userText = append(m.userText, red.Render(text))
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Type this:\n\n"
	s += m.prompt
	s += "\n\n"
	for _, l := range m.userText {
		s += l
	}
	s += "\n\n"

	elapsed := 0.0
	if !m.startTime.IsZero() {
		elapsed = m.lastPress.Sub(m.startTime).Seconds()
	}

	wpm := 0.0
	if elapsed > 0 {
		minutes := elapsed / 60.0
		wpm = (float64(m.correctKeystrokes) / 5.0) / minutes
	}

	accPct := 0.0
	if m.totalKeystrokes > 0 {
		accPct = float64(m.correctKeystrokes) / float64(m.totalKeystrokes) * 100.0
	}

	cpm := 0.0
	if elapsed > 0 {
		cpm = (float64(m.correctKeystrokes) / elapsed) * 60.0
	}

	stats := fmt.Sprintf(
		"Time: %.0fs | WPM: %.1f | Accuracy: %.1f%%\nCPM: %.1f",
		elapsed,
		wpm,
		accPct,
		cpm,
	)
	s += stats
	s += "\n"
	if len(m.prompt) <= len(m.userText) {
		s += "finished"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
