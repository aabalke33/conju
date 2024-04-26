# Conju: A Terminal-Based Language Conjugation App

# Work in Progress, do not touch

Practice conjugating words in your desired language.
Built with golang, ffmpeg, charm/bubbltea and sqlite

### REQUIRES EXPERIMENTAL VERSION
Charm's Bubbletea package only supports latin special chars in version
"github.com/charmbracelet/bubbletea v0.25.1-0.20240422164726-702b43d6b062"
and later. Assume v0.26.0 will be minimum offical verion required.

Accent Support is provided through your own keyboard configuration.

Current Support:
1. Supports Spanish Present & Preterite Tense

Roadmap:
2. Support Spanish
    a. with and without vosotros form
    b. Include pronouns
    c. include tenses
    d. Change the time commit
    e. Make pretty with Bubbletea
    f. Store User Info for data pages

Development:
1. Create Local DB for language.
2. Local DB stores user data, ie score, speed, etc and language words
