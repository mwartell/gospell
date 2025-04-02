# GoSpell

GoSpell is a terminal-based spelling practice application that uses text-to-speech to help users improve their spelling skills

![GoSpell Demo](https://img.shields.io/badge/demo-coming%20soon-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/jharlan-hash/gospell)](https://goreportcard.com/report/github.com/jharlan-hash/gospell)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

GoSpell is an interactive TUI spelling practice tool that:

- Speaks words aloud using Google Cloud Text-to-Speech
- Provides word definitions to help with context
- Tracks your correct spelling streak
- Creates a satisfying practice environment with a clean TUI

I think it would be great for students, language learners, or anyone that wants to improve their spelling skills.

## Features

- **Text-to-Speech Integration**: Hear words spoken (mostly) clearly with Google Cloud TTS
- **Word Definitions**: See word definitions to hopefully understand context and meaning
- **Progress Tracking**: Keep track of your spelling streak
- **Pretty good TUI**: Clean terminal user interface using [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Installation

### Prerequisites

- Go 1.17+
- Google Cloud Platform account with Text-to-Speech API enabled
- Google Cloud credentials JSON file (Make sure it is a key from a service account on Google Cloud)

### Install from source

```bash
# Clone the repository
git clone https://github.com/jharlan-hash/gospell.git
cd gospell

# Install dependencies
go mod download

# Build the application
go build -o gospell

# Run the application
./gospell --credentials=/path/to/your-credentials.json
```

## Usage

```bash
# Basic usage
./gospell --credentials=/path/to/your-credentials.json
```

### Key Commands

- **Enter**: Submit your spelling
- **Ctrl+R**: Repeat the current word
- **Ctrl+C/Ctrl+D/Esc**: Exit the application

## Configuration

GoSpell accepts the following command-line flags:

| Flag | Short | Description |
|------|-------|-------------|
| `--credentials` | `-c` | Path to Google Cloud credentials JSON file (required) |
| `--help` | `-h` | Display help |

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - goated TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions for TUI apps
- [Google Cloud Text-to-Speech](https://cloud.google.com/text-to-speech) - TTS API
- [Beep](https://github.com/gopxl/beep) - Audio playback 

## Contributing

Feel free to contribute if you want - this is just a personal project.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Free Dictionary API](https://dictionaryapi.dev/) for providing word definitions
- [Google Cloud Text-to-Speech](https://cloud.google.com/text-to-speech) for TTS capabilities

---

Made with ❤️ by Jack Sovern
