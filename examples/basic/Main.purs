{- 
  This example demonstrates a complete Bubbletea application.
  
  To run this example:
  1. Add purescript-bubbletea to your purescript-native project
  2. Copy this file to your src/ directory  
  3. Run with: spago run -m Main
  
  Or see it running in the midigo project: https://github.com/i-am-the-slime/midigo
-}
module Main where

import Prelude

import Bubbletea as Tea
import Data.Maybe (Maybe(..))
import Effect (Effect)
import Effect.Console (log)

-- | Application model
type Model = { message :: String, step :: Int }

-- | Initialize the model
init :: Model -> Maybe (Tea.Command Tea.Message)
init _ = Nothing

-- | Update the model based on messages
update :: Model -> Tea.Message -> Tea.UpdateResult Model
update model _ = 
  let newStep = model.step + 1
      newMessage = if newStep > 10 
                  then "Done! Press any key to exit..."
                  else "Step " <> show newStep <> " - Press any key to continue"
  in if newStep > 11
     then { model, cmd: Just Tea.quit }
     else { model: { message: newMessage, step: newStep }, cmd: Nothing }

-- | Render the model as a string
view :: Model -> String  
view model = 
  "=== Bubbletea Example ===\n\n" <>
  model.message <> "\n\n" <>
  "(This is a simple counter that will auto-exit after 10 steps)"

main :: Effect Unit
main = do
  log "Starting Bubbletea example..."
  
  let initialModel = { message: "Welcome to Bubbletea! Press any key to start...", step: 0 }
  
  program <- Tea.newProgram initialModel init update view
  result <- Tea.runProgram program
  
  case result of
    Just error -> log $ "Program error: " <> error  
    Nothing -> log "Program completed successfully!"