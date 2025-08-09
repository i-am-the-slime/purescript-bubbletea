package bubbletea

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/purescript-native/go-runtime"
)

type model struct {
	data   Any
	init   Any
	update Any
	view   Any
}

type userMsg struct{ Msg Any }

func (m *model) UpdateData(data Any) {
	m.data = data
}

func (m model) Init() tea.Cmd {
	init := m.init.(func(Any) Any)
	result := init(m.data)
	if result == nil {
		return nil
	}
	return func() tea.Msg {
		switch result := result.(type) {
		case tea.Cmd: // Sometimes this, sometimes something else, who knows anymore?
			return result()
		case func() tea.Msg: // This is the usual case
			return result()
		case func() Any:
			msg_ := result()
			return msg_.(tea.Msg)
		}
		err := fmt.Sprintf("Unknown init result type %T\n", result)
		panic(err)
		return nil
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	update := m.update.(func(Any, Any) Any)
	result_ := update(m.data, msg)
	result := result_.(Dict)
	newData := result["model"]
	m.data = newData
	cmd_ := result["cmd"]
	if cmd_ == nil {
		return m, nil
	}
	switch cmd := cmd_.(type) {
	case tea.Cmd:
		return m, cmd
	case func() Any:
		return m, func() tea.Msg {
			r_ := cmd()
			if r_ == nil {
				return nil
			}
			return r_.(tea.Msg)
		}
	default:
		err := fmt.Sprintf("Unknown cmd type %T\n", cmd)
		panic(err)
		return nil, nil
	}
}

func (m model) View() string {
	view := m.view.(func(Any) Any)
	result := view(m.data)
	return result.(string)
}
func init() {
	exports := Foreign("Bubbletea")

	exports["setWindowTitleImpl"] = func(title_ Any) Any {
		title := title_.(string)
		tea.SetWindowTitle(title)
		return nil
	}

	exports["noCommand"] = nil

	exports["newProgramImpl"] = func(modelData_ Any, init_ Any, update_ Any, view_ Any) Any {
		initialModel := model{
			data:   modelData_,
			init:   init_,
			update: update_,
			view:   view_,
		}
		return tea.NewProgram(initialModel)
	}

	exports["batch"] = func(cmds_ Any) Any {
		cmds := cmds_.([]Any)
		teaCmds := make([]tea.Cmd, 0)
		for _, cmd := range cmds {
			if cmd == nil {
				continue
			}
			switch cmd := cmd.(type) {
			case tea.Cmd:
				teaCmds = append(teaCmds, func() tea.Msg {
					if cmd == nil {
						return nil
					}
					return cmd()
				})
			case func() tea.Msg:
				teaCmds = append(teaCmds, func() tea.Msg {
					return cmd()
				})
			case func() Any:
				teaCmds = append(teaCmds, func() tea.Msg {
					result_ := cmd()
					if result_ == nil {
						return nil
					}
					return result_.(tea.Msg)
				})
			default:
				err := fmt.Sprintf("Unknown cmd type %T\n", cmd)
				panic(err)
			}
		}
		return func() Any { return tea.Batch(teaCmds...)() }
	}

	exports["runProgramImpl"] = func(program_ Any) Any {
		program := program_.(*tea.Program)
		_, err := program.Run()
		return err
	}

	exports["quit"] = func() Any { return tea.Quit() }
	exports["clearScreen"] = func() Any { return tea.ClearScreen() }

	exports["convertMessage"] = func(_converters Any) Any {
		converters := _converters.(Dict)
		return func(msg_ Any) Any {
			if msg_ == nil {
				return nil
			}
			msg := msg_.(tea.Msg)
			switch msg := msg.(type) {
			case tea.WindowSizeMsg:
				toWindowSizeMessage := converters["toWindowSizeMessage"].(func(Any, Any) Any)
				return toWindowSizeMessage(msg.Width, msg.Height)
			case tea.KeyMsg:
				toKeyMessage := converters["toKeyMessage"].(func(Any) Any)
				return toKeyMessage(msg.String())
			case userMsg:
				log.Println("Msg:", msg)
				return msg.Msg
			case Any:
				toUnknownMessage := converters["toUnknownMessage"].(func(Any) Any)
				return toUnknownMessage(msg)
			}
			panic("Unknown message type")
			return msg
		}
	}

	exports["convertMessageV"] = exports["convertMessage"]

	exports["userMessageToTeaMessage"] = func(msg Any) Any {
		return userMsg{Msg: msg}
	}

	exports["loggingToFileImpl"] = func(effect_ Any) Any {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("Error creating log file:", err)
			os.Exit(1)
		}
		// Empty the file
		if err := os.Truncate("debug.log", 0); err != nil {
			fmt.Println("Error truncating log file:", err)
			os.Exit(1)
		}
		log.Println("\n\n\n\n\n")
		log.Println("====================")
		defer f.Close()
		effect := effect_.(func() Any)
		effect()
		return nil
	}

	exports["spy"] = func(msg_ Any) Any {
		return func(v Any) Any {
			msg := msg_.(string)
			log.Printf("spy %s %s %T \n", msg, v, v)
			return v
		}
	}
}
