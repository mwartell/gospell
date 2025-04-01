package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jharlan-hash/gospell/internal/api"
	"github.com/jharlan-hash/gospell/internal/definition"
	"github.com/jharlan-hash/gospell/internal/tts"
	"github.com/jharlan-hash/gospell/internal/wpm"
	"github.com/muesli/reflow/wordwrap"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pborman/getopt"
	"google.golang.org/api/option"
)

func main() {
	credentialFlag := getopt.StringLong("credentials", 'c', "", "Path to Google Cloud credentials JSON file")
	helpFlag := getopt.BoolLong("help", 'h', "display help")

	getopt.Parse()

	if *helpFlag {
		getopt.Usage()
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	model := initialModel(*credentialFlag)
	model.ctx = ctx

	p := tea.NewProgram(&model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type wordMessage struct {
	word       string
	definition string
}

type correctMessage struct{}
type incorrectMessage struct{}

type model struct {
	textInput       textinput.Model
	client          *texttospeech.Client
	ctx             context.Context
	definition      string
	credentialPath  string
	word            string
	initialTime     time.Time
	finalTime       time.Time
	streak          int
	correction      string
	width           int
	height          int
	definitionState *definition.State
}

// initialModel initializes the model with a text input field and a random word.
func initialModel(credentialPath string) model {
	ti := textinput.New()
	ti.Placeholder = "spell spoken word..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	state := &definition.State{}
	return model{
		textInput:       ti,
		credentialPath:  credentialPath,
		correction:      "\n",
		word:            api.GetAcceptableWord(),
		definitionState: state,
		definition:      state.GetDefinition(api.GetAcceptableWord()),
	}
}

func (m *model) Init() tea.Cmd {
	if m.credentialPath != "" {
		client, err := texttospeech.NewClient(m.ctx, option.WithCredentialsFile(m.credentialPath))
		if err != nil {
			tea.ExitAltScreen()
			log.Fatal("Bad credentials file")
		}

		m.client = client
	} else {
		tea.ExitAltScreen()
		log.Fatal("No credentials file provided")
	}

	go tts.SayWord(m.ctx, *m.client, m.word)
	return textinput.Blink
}

// Command to generate a new word.
func getNewWord(m *model) tea.Cmd {
	return func() tea.Msg {
		word := api.GetAcceptableWord()
		def := m.definitionState.GetDefinition(word)

		// Play the word audio in a goroutine to avoid blocking.
		go tts.SayWord(m.ctx, *m.client, word)

		return wordMessage{
			word:       word,
			definition: def,
		}
	}
}

// Update function handling messages
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case wordMessage:
		// Update model with new word.
		m.word = msg.word
		m.definition = msg.definition
		return m, nil

	case tea.KeyMsg:
		if m.textInput.Value() == "" && msg.Type != tea.KeyCtrlR {
			m.initialTime = time.Now() // start timer on first key press.
		}

		m.finalTime = time.Now() // update timer on every key press.

		switch msg.Type {
		case tea.KeyEnter: // submit word while ignoring empty input.
			return m.submitWord()
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD: // exit.
			os.Remove("temp.wav")
			return m, tea.Quit

		case tea.KeyCtrlR: // repeat word.
            go tts.SayWord(m.ctx, *m.client, m.word)
			return m, nil

		case tea.KeyDown:
			// If the user presses down, we want to get the next definition.
			m.definition = m.definitionState.NextDefinition()

		case tea.KeyUp:
			// If the user presses up, we want to get the previous definition.
			m.definition = m.definitionState.PrevDefinition()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case correctMessage:
		m.streak++
		m.definition = wordwrap.String(m.definition, 100)

		temparr := strings.Split(m.definition, "\n")
		for index := range temparr {
			temparr[index] = "\033[32m" + strings.TrimSpace(temparr[index]) + "\033[0m"
		}

		m.definition = strings.Join(temparr, "\n")
		m.correction = "\n"
		return m, getNewWord(m)

	case incorrectMessage:
		m.streak = 0
		m.definition = wordwrap.String(m.definition, 100)

		temparr := strings.Split(m.definition, "\n")
		for index := range temparr {
			temparr[index] = "\033[31m" + strings.TrimSpace(temparr[index]) + "\033[0m"
		}

		m.definition = strings.Join(temparr, "\n")
		m.correction = fmt.Sprintf("\nCorrect spelling: %s", m.word)
		return m, getNewWord(m)
	}

	// Handle text input updates.
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) submitWord() (tea.Model, tea.Cmd) {
	if m.textInput.Value() == "" {
		return m, nil
	}

	userInput := m.textInput.Value()
	m.textInput.Reset()

	if userInput == m.word { // Correct answer.
		return m, func() tea.Msg { return correctMessage{} }
	} else { // Incorrect answer.
		return m, func() tea.Msg { return incorrectMessage{} }
	}

}

func (m model) View() string {
	// Create a container style that will be centered as a whole.
	containerStyle := lipgloss.NewStyle().Padding(1, 2)

	// Center the welcome text within the container.
	welcomeText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100). // Set a fixed width for the centered elements.
		Render("Welcome to gospell!")

	// Center the input field within the container.
	inputView := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(m.textInput.View())

		// Center the definition but keep it within the container's width.
	definitionText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(m.definition)

	correctionText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(m.correction)

	wpmText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render(fmt.Sprintf("WPM: %d", wpm.CalculateWpm(m.textInput.Value(), m.initialTime, m.finalTime)))

	// Center the streak counter.
	streakText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(100).
		Render("Streak: " + fmt.Sprintf("%d", m.streak))

	// Combine all elements with the container style.
	content := containerStyle.Render(
		welcomeText + "\n\n" +
			wpmText + "\n" +
			inputView + "\n" +
			definitionText + "\n" +
			// wordwrap.String(m.definition, 100) + "\n" +
			correctionText + "\n" +
			streakText,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
