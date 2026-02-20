package completion

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Command struct {
	Name        string
	Aliases     []string
	Short       string
	Long        string
	Subcommands []Command
	Flags       []Flag
	Args        []Arg
	Hidden      bool
}

type Flag struct {
	Name       string
	Shorthand  string
	Usage      string
	Default    string
	Required   bool
	EnumValues []string
	Hidden     bool
}

type Arg struct {
	Name       string
	Usage      string
	Required   bool
	EnumValues []string
}

type Generator struct {
	cmd      Command
	progName string
}

func NewGenerator(cmd Command, progName string) *Generator {
	return &Generator{
		cmd:      cmd,
		progName: progName,
	}
}

func (g *Generator) GenerateBash(w io.Writer) error {
	script := `#!/bin/bash
# Bash completion for ` + g.progName + `

_` + g.progName + `_completion() {
    local cur prev words cword
    _init_completion || return

    local commands="` + strings.Join(g.getCommandNames(), " ") + `"

    if [[ ${cword} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
        return
    fi

    local command=${words[1]}
    case ${command} in
`

	for _, subcmd := range g.cmd.Subcommands {
		flags := g.getFlagNames(subcmd)
		if len(flags) > 0 {
			script += `        ` + subcmd.Name + `)
            COMPREPLY=($(compgen -W "` + strings.Join(flags, " ") + `" -- "${cur}"))
            return
            ;;
`
		}
	}

	script += `    esac
}

complete -F _` + g.progName + `_completion ` + g.progName + `
`

	_, err := w.Write([]byte(script))
	return err
}

func (g *Generator) GenerateZsh(w io.Writer) error {
	script := `#compdef ` + g.progName + `

` + g.progName + `() {
    local -a commands
    commands=(
`

	for _, subcmd := range g.cmd.Subcommands {
		if !subcmd.Hidden {
			escapedDesc := escapeZshString(subcmd.Short)
			script += fmt.Sprintf("        '%s:\"%s'\"\n", subcmd.Name, escapedDesc)
		}
	}

	script += `    )

    _arguments -C \
        "--help[Show help]" \
        "-h[Show help]" \
        "1: :->commands" \
        "*::arg:->args"

    case $state in
        commands)
            _describe "command" commands
            ;;
        args)
            case ${words[1]} in
`

	for _, subcmd := range g.cmd.Subcommands {
		flags := g.generateZshFlags(subcmd)
		if len(flags) > 0 {
			script += `                ` + subcmd.Name + `)
                    _arguments -C \
` + flags + `
                    ;;
`
		}
	}

	script += `            esac
            ;;
    esac
}

` + g.progName + `
`

	_, err := w.Write([]byte(script))
	return err
}

func (g *Generator) generateZshFlags(cmd Command) string {
	var flags []string
	for _, flag := range cmd.Flags {
		if flag.Hidden {
			continue
		}

		flagDef := `'--` + flag.Name + `[` + escapeZshString(flag.Usage) + `]'`

		if len(flag.EnumValues) > 0 {
			flagDef = `'(--` + flag.Name + `)--` + flag.Name + `[` + escapeZshString(flag.Usage) + `]:` + flag.Name + `:->` + flag.Name + `'`
		}

		if flag.Shorthand != "" {
			flagDef += ` '-` + flag.Shorthand + `[` + escapeZshString(flag.Usage) + `]'`
		}

		flags = append(flags, "                        "+flagDef)
	}
	return strings.Join(flags, " \\\n")
}

func (g *Generator) GenerateFish(w io.Writer) error {
	script := `# Fish completion for ` + g.progName + `

complete -c ` + g.progName + ` -f

# Main commands
`

	for _, subcmd := range g.cmd.Subcommands {
		if !subcmd.Hidden {
			script += `complete -c ` + g.progName + ` -n '__fish_use_subcommand' -a '` + subcmd.Name + `' -d '` + escapeFishString(subcmd.Short) + `'
`
		}
	}

	script += `
# Subcommand flags
`

	for _, subcmd := range g.cmd.Subcommands {
		for _, flag := range subcmd.Flags {
			if flag.Hidden {
				continue
			}

			script += `complete -c ` + g.progName + ` -n '__fish_seen_subcommand_from ` + subcmd.Name + `' -l ` + flag.Name + ` -d '` + escapeFishString(flag.Usage) + `'`

			if len(flag.EnumValues) > 0 {
				script += ` -xa '` + strings.Join(flag.EnumValues, " ") + `'`
			}

			script += `
`

			if flag.Shorthand != "" {
				script += `complete -c ` + g.progName + ` -n '__fish_seen_subcommand_from ` + subcmd.Name + `' -s ` + flag.Shorthand + ` -d '` + escapeFishString(flag.Usage) + `'
`
			}
		}
	}

	_, err := w.Write([]byte(script))
	return err
}

func (g *Generator) GeneratePowerShell(w io.Writer) error {
	script := `# PowerShell completion for ` + g.progName + `

using namespace System.Management.Automation
using namespace System.Management.Automation.Language

Register-ArgumentCompleter -Native -CommandName ` + g.progName + ` -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commands = @(
`

	for _, subcmd := range g.cmd.Subcommands {
		if !subcmd.Hidden {
			script += `        [CompletionResult]::new('` + subcmd.Name + `', '` + subcmd.Name + `', [CompletionResultType]::ParameterValue, '` + escapePowerShellString(subcmd.Short) + `')
`
		}
	}

	script += `    )

    $commandElements = $commandAst.CommandElements
    if ($commandElements.Count -eq 1) {
        return $commands
    }

    $subcommand = $commandElements[1].Value
    switch ($subcommand) {
`

	for _, subcmd := range g.cmd.Subcommands {
		flags := g.generatePowerShellFlags(subcmd)
		if len(flags) > 0 {
			script += `        '` + subcmd.Name + `' {
            @(
` + flags + `
            )
        }
`
		}
	}

	script += `        default { @() }
    }
}
`

	_, err := w.Write([]byte(script))
	return err
}

func (g *Generator) generatePowerShellFlags(cmd Command) string {
	var flags []string
	for _, flag := range cmd.Flags {
		if flag.Hidden {
			continue
		}
		flags = append(flags, `                [CompletionResult]::new('--`+flag.Name+`', '--`+flag.Name+`', [CompletionResultType]::ParameterName, '`+escapePowerShellString(flag.Usage)+`')`)
	}
	return strings.Join(flags, ",\n")
}

func (g *Generator) getCommandNames() []string {
	var names []string
	for _, subcmd := range g.cmd.Subcommands {
		if !subcmd.Hidden {
			names = append(names, subcmd.Name)
			names = append(names, subcmd.Aliases...)
		}
	}
	return names
}

func (g *Generator) getFlagNames(cmd Command) []string {
	var names []string
	for _, flag := range cmd.Flags {
		if !flag.Hidden {
			names = append(names, "--"+flag.Name)
			if flag.Shorthand != "" {
				names = append(names, "-"+flag.Shorthand)
			}
		}
	}
	return names
}

func escapeZshString(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

func escapeFishString(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func escapePowerShellString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func InstallCompletion(shell string, progName string) error {
	generator := NewGenerator(GetDefaultCommand(), progName)

	var script strings.Builder
	var err error

	switch shell {
	case "bash":
		err = generator.GenerateBash(&script)
	case "zsh":
		err = generator.GenerateZsh(&script)
	case "fish":
		err = generator.GenerateFish(&script)
	case "powershell", "pwsh":
		err = generator.GeneratePowerShell(&script)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	if err != nil {
		return err
	}

	fmt.Println(script.String())
	return nil
}

func GetDefaultCommand() Command {
	return Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize BIOMETRICS repository",
				Flags: []Flag{
					{Name: "force", Shorthand: "f", Usage: "Force initialization"},
				},
			},
			{
				Name:  "config",
				Short: "Manage configuration",
				Subcommands: []Command{
					{Name: "init", Short: "Create default configuration"},
					{Name: "validate", Short: "Validate existing configuration"},
					{Name: "show", Short: "Display current configuration"},
				},
			},
			{
				Name:  "audit",
				Short: "Query and manage audit logs",
				Subcommands: []Command{
					{Name: "query", Short: "Query audit events"},
					{Name: "export", Short: "Export audit events"},
					{Name: "stats", Short: "Show audit statistics"},
					{Name: "cleanup", Short: "Remove old audit logs"},
					{Name: "rotate", Short: "Rotate audit log files"},
				},
				Flags: []Flag{
					{Name: "start-time", Usage: "Start time filter"},
					{Name: "end-time", Usage: "End time filter"},
					{Name: "event-types", Usage: "Comma-separated event types"},
					{Name: "actors", Usage: "Comma-separated actor IDs"},
					{Name: "format", Usage: "Output format", EnumValues: []string{"table", "json", "csv"}},
					{Name: "output", Shorthand: "o", Usage: "Output file path"},
				},
			},
			{
				Name:  "version",
				Short: "Show version information",
			},
		},
		Flags: []Flag{
			{Name: "help", Shorthand: "h", Usage: "Show help"},
			{Name: "verbose", Shorthand: "v", Usage: "Verbose output"},
			{Name: "config", Shorthand: "c", Usage: "Config file path"},
		},
	}
}

func PrintCompletionHelp() {
	fmt.Print(`Shell Completion for BIOMETRICS CLI

Usage:
  biometrics completion <shell>

Shells:
  bash        Bash completion
  zsh         Zsh completion
  fish        Fish completion
  powershell  PowerShell completion

Installation:

Bash:
  biometrics completion bash > /etc/bash_completion.d/biometrics
  source ~/.bashrc

Zsh:
  biometrics completion zsh > "${fpath[1]}/_biometrics"
  autoload -U compinit && compinit

Fish:
  biometrics completion fish > ~/.config/fish/completions/biometrics.fish

PowerShell:
  biometrics completion powershell > $PROFILE
  . $PROFILE
`)
}

func WriteCompletionToFile(shell, progName, filePath string) error {
	generator := NewGenerator(GetDefaultCommand(), progName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch shell {
	case "bash":
		return generator.GenerateBash(file)
	case "zsh":
		return generator.GenerateZsh(file)
	case "fish":
		return generator.GenerateFish(file)
	case "powershell", "pwsh":
		return generator.GeneratePowerShell(file)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}
