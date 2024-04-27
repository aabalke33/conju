# TUI Map

The Bubbletea TUI is broken into nested models:

Main Model: Stores State
    Setting Model: Stores Settings for Conju Game
        Language Model: View to select Language
        Tense Model: View to select Tense
        Duration Model: View to select Duration
        Confirmation Model: View to confirm settings before game
        (item): Model for lists used in Language and Tense
    Game Model: Stores game state (verbs, score, time left)
        Round Model: View current round and verb
    Performance Model: Displays final assessment info
