# purescript-bubbletea

A powerful TUI (Terminal User Interface) framework for PureScript-native, built on top of Charm Bracelet's Bubbletea.

Build rich, interactive terminal applications with the Elm architecture pattern: Model-Update-View.

## Installation

```yaml
# spago.yaml
dependencies:
  - bubbletea: "https://github.com/i-am-the-slime/purescript-bubbletea.git"
```

## Quick Start

```purescript
module Main where

import Prelude
import Bubbletea as Tea
import Effect (Effect)
import Effect.Console (log)
import Data.Maybe (Maybe(..))

type Model = { count :: Int }

data Msg = Increment | Decrement | Quit

init :: Model -> Maybe (Tea.Command Tea.Message)  
init _ = Nothing

update :: Model -> Tea.Message -> Tea.UpdateResult Model
update model msg = case parseMessage msg of
  Just Increment -> { model: model { count = model.count + 1 }, cmd: Nothing }
  Just Decrement -> { model: model { count = model.count - 1 }, cmd: Nothing }  
  Just Quit -> { model, cmd: Just Tea.quit }
  Nothing -> { model, cmd: Nothing }

view :: Model -> String
view model = 
  "Counter: " <> show model.count <> "\n\n" <>
  "Press '+' to increment, '-' to decrement, 'q' to quit"

main :: Effect Unit
main = do
  program <- Tea.newProgram { count: 0 } init update view
  result <- Tea.runProgram program
  case result of
    Just error -> log $ "Error: " <> error
    Nothing -> log "Program exited successfully"
```

## Features

- **Elm Architecture**: Predictable state management with Model-Update-View pattern
- **Rich Messaging**: Handle keyboard input, window resize events, and custom messages
- **Command System**: Perform side effects and async operations with commands
- **Cross-platform**: Works on all platforms supported by Go
- **Type-safe**: Fully typed PureScript interface with Go performance

## API

### Core Functions

- `newProgram :: Model -> Init -> Update -> View -> Effect Program` - Create a new TUI program
- `runProgram :: Program -> Effect (Maybe String)` - Run the program (blocks until exit)
- `setWindowTitle :: String -> Effect Unit` - Set terminal window title

### Commands  

- `quit :: Command Message` - Exit the program
- `clearScreen :: Command Message` - Clear the terminal screen
- `batchCommands :: Array (Command Message) -> Command Message` - Batch multiple commands
- `noCommand :: Maybe (Command Message)` - No-op command

### Message Handling

Programs receive various message types:
- **Keyboard events** - Key presses and combinations
- **Window events** - Terminal resize, focus changes  
- **Custom messages** - User-defined application messages

### Debugging

- `enableFileLogging :: Effect Unit -> Effect Unit` - Enable debug logging to file
- `debug :: String -> a -> a` - Debug utility for tracing values

## Architecture

Bubbletea follows the Elm architecture:

1. **Model** - Your application state
2. **Init** - Initialize the model and return initial commands  
3. **Update** - Handle messages and update the model
4. **View** - Render the current model as a string

The runtime handles the event loop, calling your update function when messages arrive and re-rendering when the model changes.

## Cross-Platform Support

This library works on all platforms supported by Go:
- macOS
- Linux  
- Windows
- BSD variants

## Requirements

- PureScript with purescript-native backend
- Go 1.24.3+

## Part of Charm Bracelet Ecosystem

Works great with:
- [purescript-lipgloss](https://github.com/i-am-the-slime/purescript-lipgloss) - Styling and layout
- [purescript-bubbles](https://github.com/i-am-the-slime/purescript-bubbles) - UI components