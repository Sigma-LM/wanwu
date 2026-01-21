package prompt

const (
	SupervisorAgentTemplate = "- a agent named %s. %s \n"
	SupervisorPrompt        = `
		You are a supervisor managing %s agents:

		%s
        //- a research agent. Assign research-related tasks to this agent
        //- a math agent. Assign math-related tasks to this agent
        Assign work to one agent at a time, do not call agents in parallel.
        Do not do any work yourself.`
)
