# Conju: A Terminal-Based Language Learning App

Conju is a TUI Language Learning app, currently Supporting language conjugation
assessments/games.

Practice conjugating words in your desired language.
Built with golang, ffmpeg, charm/bubbletea and sqlite
![conju225](https://github.com/aabalke33/conju/assets/22086435/c16f0d6f-e83d-4581-a03a-37f77a5fdba3)

## Youtube Breakdown
[Youtube Video](https://www.youtube.com/watch?v=fB5a8g0nMJY)

## Architecture / Design

Conju is built using Golang and exports to a binary. The TUI itself uses the
[Charm Bubbletea TUI framework, and additional libraries lipgloss, and Bubbles](https://github.com/charmbracelet/bubbletea).
This framework allows easy creation of nested TUI models, one for each view,
by following the [Elm architecture](https://guide.elm-lang.org/architecture/) for interactuve programs.

### Models
Main Model: the main model handles the state of the entire app
    Setting Model: Handles the state of which setting options are selected
        Language Selection Model
        Tense Selection Model
        Duration Input Model
        Confirmation Model
    Game Model: Handles the game state, ex. count of correct options
        Round Model: Repeated randomized model for each verb/conjugation match
    Performance Model: Final model to display final data to the user

Scores are exported to ./data/conju.csv, a csv file to keep track of progress for
personal benchmarking.

### Audio
Uses ffmpeg/ffplay to play small mp3 files included in the repo for audio feedback during the game.
This will require ffmpeg is install on the device.

## Data Storage

Language data is stored in individual sqlite databases in ./data, with tables per language tense.
Each table then has the following column template:

| infinitve | meaning | conjugation | conjugation... | ... |
|-----------|---------|-------------|----------------|-----|
| verb1     |         |             |                |     |
| ser       |  to be  |    soy      |     eres       | es  |

In many languages, the conjugations will be broken down by pronouns. For example,
in Spanish, the first/second/third person * single/plural conjugations would follow the meaning column.

## Contributing
If you would like to contribute language databases, please do! Follow the template
outlined in the Data Storage section, and check the Spanish.db file for an example.

## Roadmap
- Add Pronoun Setting/Menu Selection
- Adaptive Color profiles depending on support
- Add All Spanish Tenses
- Keep track of the most common mistake words
- Add preposition options to conjugations
- Add verb meaning = verb game mode

## Misc Notes
- Requires Unofficial Bubbletea release for special char support in Windows, see go.mod
- Special Character support is provided by keyboard configuration and is not built into the app.
- The archive folder holds previously used SQL Queries to provide an audit trail of db creation.
