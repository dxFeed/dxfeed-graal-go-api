package main

type Tool interface {
	Run(args []string)
}

func Run(args []string) {
	tool := createTool(args)
	Tool.Run(tool, args)
}

func createTool(args []string) Tool {
	switch args[1] {
	case "connect":
		return Connect{}
	default:
		return nil
	}
}
