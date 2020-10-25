package command_util

// Command from user input
type Command struct {
	Method  string
	Params  []string
}

const (
	CommandJoin 	= "join"
	CommandLeave 	= "leave"
	CommandSend 	= "send"
	CommandSwitch 	= "switch"
	CommandDisplay 	= "display"

	CommandPut 		= "put"
	CommandGet 		= "get"
	CommandDelete 	= "delete"
	CommandList 	= "ls"
	CommandStore 	= "store"

	CommandQuit 	= "quit"
)

func IsFileCommand(command Command) bool {
	fileCommands := [...]string{
		CommandGet,
		CommandPut,
		CommandDelete,
		CommandList,
		CommandStore,
	}
	for _, s := range fileCommands {
		if command.Method == s {
			return true
		}
	}
	return false
}

func IsMemberCommand(command Command) bool {
	memberCommands := [...]string{
		CommandJoin,
		CommandLeave,
		CommandDisplay,
		CommandSwitch,
	}
	for _, s := range memberCommands {
		if command.Method == s {
			return true
		}
	}
	return false
}

