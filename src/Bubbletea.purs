module Bubbletea where

import Prelude

import Data.Function.Uncurried (Fn1, Fn2, Fn3, Fn4, runFn1, runFn2, runFn3, runFn4)
import Data.Maybe (Maybe(..))
import Effect (Effect)
import Effect.Uncurried (EffectFn1, EffectFn2, EffectFn3, EffectFn4, runEffectFn1, runEffectFn2, runEffectFn3, runEffectFn4)

-- | Core Bubbletea types
foreign import data Program :: Type
foreign import data Command :: Type -> Type
foreign import data Message :: Type

-- | Model update result containing new model and command
type UpdateResult model = { model :: model, cmd :: Maybe (Command Message) }

-- | Message converters for handling different Bubbletea message types
type MessageConverters = 
  { toWindowSizeMessage :: Int -> Int -> Message
  , toKeyMessage :: String -> Message
  , toUnknownMessage :: Message -> Message
  }

-- | Create a new Bubbletea program
foreign import newProgramImpl :: forall model. EffectFn4 model (model -> Maybe (Command Message)) (model -> Message -> UpdateResult model) (model -> String) Program

newProgram :: forall model. 
  model ->
  (model -> Maybe (Command Message)) ->
  (model -> Message -> UpdateResult model) ->
  (model -> String) ->
  Effect Program
newProgram model init update view = runEffectFn4 newProgramImpl model init update view

-- | Run a Bubbletea program
foreign import runProgramImpl :: EffectFn1 Program (Maybe String)

runProgram :: Program -> Effect (Maybe String)
runProgram = runEffectFn1 runProgramImpl

-- | Set window title
foreign import setWindowTitleImpl :: EffectFn1 String Unit

setWindowTitle :: String -> Effect Unit
setWindowTitle = runEffectFn1 setWindowTitleImpl

-- | No-op command
foreign import noCommand :: Maybe (Command Message)

-- | Batch multiple commands
foreign import batch :: Fn1 (Array (Command Message)) (Command Message)

batchCommands :: Array (Command Message) -> Command Message
batchCommands = runFn1 batch

-- | Quit the program
foreign import quit :: Command Message

-- | Clear the screen
foreign import clearScreen :: Command Message

-- | Convert Bubbletea messages to PureScript messages
foreign import convertMessage :: Fn1 MessageConverters (Message -> Maybe Message)

convertMessage' :: MessageConverters -> Message -> Maybe Message
convertMessage' = runFn1 convertMessage

-- | Convert variant of convertMessage
foreign import convertMessageV :: Fn1 MessageConverters (Message -> Maybe Message)

-- | Convert user message to tea message
foreign import userMessageToTeaMessage :: Fn1 Message Message

userToTeaMessage :: Message -> Message
userToTeaMessage = runFn1 userMessageToTeaMessage

-- | Enable logging to file (for debugging)
foreign import loggingToFileImpl :: EffectFn1 (Effect Unit) Unit

enableFileLogging :: Effect Unit -> Effect Unit
enableFileLogging = runEffectFn1 loggingToFileImpl

-- | Debug spy function
foreign import spy :: Fn1 String (forall a. a -> a)

debug :: forall a. String -> a -> a
debug msg = runFn1 (spy msg)