package config

type Problem struct {
	ProblemId   string `json:"problem_id"`
	ProblemName string `json:"problem_name"`
}

type Compiler struct {
	CompilerId   string `json:"compiler_id"`
	CompilerName string `json:"compiler_name"`
	CompilerLang string `json:"compiler_lang"`
}
