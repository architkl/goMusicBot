package framework

type (
	Middleware func(Context) error
	Command    func(Context)

	CommandStruct struct {
		middleware Middleware
		command    Command
		help       string
	}

	CmdMap map[string]CommandStruct

	CommandHandler struct {
		cmds CmdMap
	}
)

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
}

func (handler CommandHandler) GetCmds() CmdMap {
	return handler.cmds
}

func (handler CommandHandler) Get(name string) (*Middleware, *Command, bool) {
	cmd, found := handler.cmds[name]
	return &cmd.middleware, &cmd.command, found
}

func (handler CommandHandler) Register(name string, middleware Middleware, command Command, helpmsg string) {
	// Massage the arguments into a "Full command"
	cmdstruct := CommandStruct{middleware: middleware, command: command, help: helpmsg}
	handler.cmds[name] = cmdstruct
}

func (command CommandStruct) GetHelp() string {
	return command.help
}