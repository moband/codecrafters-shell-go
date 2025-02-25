package main

type Shell struct {
	builtins map[string]func(*CommandHandler)
}

func NewShell() *Shell {
	s := &Shell{
		builtins: make(map[string]func(*CommandHandler)),
	}
	s.registerBuiltins()
	return s
}

func (s *Shell) registerBuiltins() {
	s.builtins["echo"] = (*CommandHandler).handleEcho
	s.builtins["exit"] = (*CommandHandler).handleExit
	s.builtins["pwd"] = (*CommandHandler).handlePwd
	s.builtins["cd"] = (*CommandHandler).handleCd
	s.builtins["type"] = (*CommandHandler).handleType
}
