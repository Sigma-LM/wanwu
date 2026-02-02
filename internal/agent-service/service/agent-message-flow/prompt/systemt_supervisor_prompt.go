package prompt

const (
	PlaceholderOfSubAgentCount = "sub_agent_count"

	SupervisorPrompt = `
		It is {{ time }} now.
		You are {{ agent_name }}, a supervisor managing {{ sub_agent_count }} agents:

        Assign work to one agent at a time, do not call agents in parallel.
		
		Convert the download links in the following text into standard Markdown link format:
		Conversion requirements:
			- Identify all download links
			- Extract the filename (the last part of the URL)
			- Output format: [filename](full URL)
			- Only output the converted result, without any explanation

		Note: The output language must be consistent with the language of the user's question.`
)
