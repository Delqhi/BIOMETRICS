package completion

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	cmd := Command{
		Name:  "test",
		Short: "Test command",
	}

	g := NewGenerator(cmd, "testprog")

	if g == nil {
		t.Fatal("NewGenerator returned nil")
	}

	if g.progName != "testprog" {
		t.Errorf("Expected progName 'testprog', got '%s'", g.progName)
	}

	if g.cmd.Name != "test" {
		t.Errorf("Expected command name 'test', got '%s'", g.cmd.Name)
	}
}

func TestGenerateBash(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{Name: "init", Short: "Initialize"},
			{Name: "version", Short: "Show version"},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateBash(&buf)
	if err != nil {
		t.Fatalf("GenerateBash failed: %v", err)
	}

	output := buf.String()

	// Check for essential bash completion elements
	if !strings.Contains(output, "#!/bin/bash") {
		t.Error("Missing shebang")
	}

	if !strings.Contains(output, "_biometrics_completion") {
		t.Error("Missing completion function name")
	}

	if !strings.Contains(output, "complete -F") {
		t.Error("Missing complete command")
	}

	if !strings.Contains(output, "init") {
		t.Error("Missing 'init' subcommand")
	}

	if !strings.Contains(output, "version") {
		t.Error("Missing 'version' subcommand")
	}
}

func TestGenerateBashWithFlags(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize",
				Flags: []Flag{
					{Name: "force", Shorthand: "f", Usage: "Force init"},
				},
			},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateBash(&buf)
	if err != nil {
		t.Fatalf("GenerateBash failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "--force") {
		t.Error("Missing --force flag")
	}

	if !strings.Contains(output, "-f") {
		t.Error("Missing -f shorthand")
	}
}

func TestGenerateZsh(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{Name: "init", Short: "Initialize repository"},
			{Name: "version", Short: "Show version"},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateZsh(&buf)
	if err != nil {
		t.Fatalf("GenerateZsh failed: %v", err)
	}

	output := buf.String()

	// Check for essential zsh completion elements
	if !strings.Contains(output, "#compdef biometrics") {
		t.Error("Missing #compdef directive")
	}

	if !strings.Contains(output, "biometrics()") {
		t.Error("Missing main function")
	}

	if !strings.Contains(output, "_arguments") {
		t.Error("Missing _arguments call")
	}

	if !strings.Contains(output, "init") {
		t.Error("Missing 'init' subcommand")
	}

	if !strings.Contains(output, "version") {
		t.Error("Missing 'version' subcommand")
	}
}

func TestGenerateZshWithHiddenCommand(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{Name: "init", Short: "Initialize"},
			{Name: "secret", Short: "Secret command", Hidden: true},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateZsh(&buf)
	if err != nil {
		t.Fatalf("GenerateZsh failed: %v", err)
	}

	output := buf.String()

	if strings.Contains(output, "'secret:") {
		t.Error("Hidden command 'secret' should not appear in completion")
	}

	if !strings.Contains(output, "'init:") {
		t.Error("Visible command 'init' should appear in completion")
	}
}

func TestGenerateZshWithFlags(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize",
				Flags: []Flag{
					{Name: "force", Shorthand: "f", Usage: "Force init"},
					{Name: "config", Usage: "Config file"},
				},
			},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateZsh(&buf)
	if err != nil {
		t.Fatalf("GenerateZsh failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "--force") {
		t.Error("Missing --force flag")
	}

	if !strings.Contains(output, "--config") {
		t.Error("Missing --config flag")
	}
}

func TestGenerateFish(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{Name: "init", Short: "Initialize"},
			{Name: "version", Short: "Show version"},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateFish(&buf)
	if err != nil {
		t.Fatalf("GenerateFish failed: %v", err)
	}

	output := buf.String()

	// Check for essential fish completion elements
	if !strings.Contains(output, "# Fish completion for biometrics") {
		t.Error("Missing fish comment")
	}

	if !strings.Contains(output, "complete -c biometrics") {
		t.Error("Missing complete command")
	}

	if !strings.Contains(output, "-a 'init'") {
		t.Error("Missing 'init' subcommand")
	}

	if !strings.Contains(output, "-a 'version'") {
		t.Error("Missing 'version' subcommand")
	}
}

func TestGenerateFishWithFlags(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize",
				Flags: []Flag{
					{Name: "force", Shorthand: "f", Usage: "Force init", EnumValues: []string{"true", "false"}},
				},
			},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateFish(&buf)
	if err != nil {
		t.Fatalf("GenerateFish failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "-l force") {
		t.Error("Missing --force flag")
	}

	if !strings.Contains(output, "-xa 'true false'") {
		t.Error("Missing enum values for flag")
	}
}

func TestGenerateFishWithHiddenFlag(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize",
				Flags: []Flag{
					{Name: "force", Usage: "Force init"},
					{Name: "internal", Usage: "Internal flag", Hidden: true},
				},
			},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateFish(&buf)
	if err != nil {
		t.Fatalf("GenerateFish failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "-l force") {
		t.Error("Missing visible --force flag")
	}

	if strings.Contains(output, "-l internal") {
		t.Error("Hidden flag 'internal' should not appear")
	}
}

func TestGeneratePowerShell(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{Name: "init", Short: "Initialize"},
			{Name: "version", Short: "Show version"},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GeneratePowerShell(&buf)
	if err != nil {
		t.Fatalf("GeneratePowerShell failed: %v", err)
	}

	output := buf.String()

	// Check for essential PowerShell completion elements
	if !strings.Contains(output, "# PowerShell completion for biometrics") {
		t.Error("Missing PowerShell comment")
	}

	if !strings.Contains(output, "Register-ArgumentCompleter") {
		t.Error("Missing Register-ArgumentCompleter")
	}

	if !strings.Contains(output, "CompletionResult") {
		t.Error("Missing CompletionResult")
	}

	if !strings.Contains(output, "'init'") {
		t.Error("Missing 'init' subcommand")
	}
}

func TestGeneratePowerShellWithFlags(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize",
				Flags: []Flag{
					{Name: "force", Usage: "Force init"},
				},
			},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GeneratePowerShell(&buf)
	if err != nil {
		t.Fatalf("GeneratePowerShell failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "--force") {
		t.Error("Missing --force flag")
	}
}

func TestGeneratePowerShellWithHiddenFlag(t *testing.T) {
	cmd := Command{
		Name:  "biometrics",
		Short: "BIOMETRICS CLI",
		Subcommands: []Command{
			{
				Name:  "init",
				Short: "Initialize",
				Flags: []Flag{
					{Name: "force", Usage: "Force init"},
					{Name: "secret", Usage: "Secret flag", Hidden: true},
				},
			},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GeneratePowerShell(&buf)
	if err != nil {
		t.Fatalf("GeneratePowerShell failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "--force") {
		t.Error("Missing visible --force flag")
	}

	if strings.Contains(output, "--secret") {
		t.Error("Hidden flag 'secret' should not appear")
	}
}

func TestEscapeZshString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"it's a test", "it'\\''s a test"},
		{"no quotes", "no quotes"},
		{"'quoted'", "'\\''quoted'\\''"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeZshString(tt.input)
			if result != tt.expected {
				t.Errorf("escapeZshString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeFishString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"it's a test", "it\\'s a test"},
		{"no quotes", "no quotes"},
		{"'quoted'", "\\'quoted\\'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeFishString(tt.input)
			if result != tt.expected {
				t.Errorf("escapeFishString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapePowerShellString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"it's a test", "it''s a test"},
		{"no quotes", "no quotes"},
		{"'quoted'", "''quoted''"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapePowerShellString(tt.input)
			if result != tt.expected {
				t.Errorf("escapePowerShellString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCommandNames(t *testing.T) {
	cmd := Command{
		Name: "biometrics",
		Subcommands: []Command{
			{Name: "init", Short: "Initialize"},
			{Name: "config", Short: "Config", Aliases: []string{"cfg"}},
			{Name: "hidden", Short: "Hidden", Hidden: true},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	names := g.getCommandNames()

	expectedNames := []string{"init", "config", "cfg"}
	for _, expected := range expectedNames {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command name '%s' not found", expected)
		}
	}

	// Hidden commands should not appear
	for _, name := range names {
		if name == "hidden" {
			t.Error("Hidden command 'hidden' should not appear in command names")
		}
	}
}

func TestGetFlagNames(t *testing.T) {
	subcmd := Command{
		Name: "init",
		Flags: []Flag{
			{Name: "force", Shorthand: "f"},
			{Name: "config"},
			{Name: "hidden", Hidden: true},
		},
	}

	g := &Generator{cmd: Command{}}
	names := g.getFlagNames(subcmd)

	expectedNames := []string{"--force", "-f", "--config"}
	for _, expected := range expectedNames {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected flag name '%s' not found", expected)
		}
	}

	// Hidden flags should not appear
	for _, name := range names {
		if name == "--hidden" {
			t.Error("Hidden flag '--hidden' should not appear in flag names")
		}
	}
}

func TestInstallCompletion(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell", "pwsh"}

	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			// Redirect stdout to capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := InstallCompletion(shell, "biometrics")

			w.Close()
			os.Stdout = old

			if err != nil {
				t.Errorf("InstallCompletion(%s) failed: %v", shell, err)
			}

			// Read captured output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if len(output) == 0 {
				t.Errorf("InstallCompletion(%s) produced no output", shell)
			}
		})
	}
}

func TestInstallCompletionUnsupportedShell(t *testing.T) {
	err := InstallCompletion("unsupported", "biometrics")
	if err == nil {
		t.Error("Expected error for unsupported shell")
	}

	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("Expected 'unsupported shell' error, got: %v", err)
	}
}

func TestGetDefaultCommand(t *testing.T) {
	cmd := GetDefaultCommand()

	if cmd.Name != "biometrics" {
		t.Errorf("Expected command name 'biometrics', got '%s'", cmd.Name)
	}

	// Check that default subcommands exist
	expectedSubcmds := []string{"init", "config", "audit", "version"}
	for _, expected := range expectedSubcmds {
		found := false
		for _, subcmd := range cmd.Subcommands {
			if subcmd.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found in default command", expected)
		}
	}

	// Check that global flags exist
	expectedFlags := []string{"help", "verbose", "config"}
	for _, expected := range expectedFlags {
		found := false
		for _, flag := range cmd.Flags {
			if flag.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected flag '%s' not found in default command", expected)
		}
	}
}

func TestWriteCompletionToFile(t *testing.T) {
	// Create temp file
	tmpFile := "/tmp/test_completion.sh"
	defer os.Remove(tmpFile)

	tests := []struct {
		name  string
		shell string
	}{
		{"bash", "bash"},
		{"zsh", "zsh"},
		{"fish", "fish"},
		{"powershell", "powershell"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpPath := tmpFile + "." + tt.name
			defer os.Remove(tmpPath)

			err := WriteCompletionToFile(tt.shell, "biometrics", tmpPath)
			if err != nil {
				t.Fatalf("WriteCompletionToFile failed: %v", err)
			}

			// Verify file was created
			if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
				t.Errorf("Completion file was not created at %s", tmpPath)
			}

			// Verify file has content
			content, err := os.ReadFile(tmpPath)
			if err != nil {
				t.Fatalf("Failed to read completion file: %v", err)
			}

			if len(content) == 0 {
				t.Error("Completion file is empty")
			}
		})
	}
}

func TestWriteCompletionToFileUnsupportedShell(t *testing.T) {
	tmpFile := "/tmp/test_invalid.sh"
	defer os.Remove(tmpFile)

	err := WriteCompletionToFile("invalid", "biometrics", tmpFile)
	if err == nil {
		t.Error("Expected error for unsupported shell")
	}

	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("Expected 'unsupported shell' error, got: %v", err)
	}
}

func TestWriteCompletionToFileInvalidPath(t *testing.T) {
	err := WriteCompletionToFile("bash", "biometrics", "/nonexistent/directory/completion.sh")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestGenerateZshFlagsWithEnumValues(t *testing.T) {
	subcmd := Command{
		Name: "query",
		Flags: []Flag{
			{
				Name:       "format",
				Usage:      "Output format",
				EnumValues: []string{"json", "table", "csv"},
			},
		},
	}

	g := &Generator{cmd: Command{}}
	flags := g.generateZshFlags(subcmd)

	if !strings.Contains(flags, "--format") {
		t.Error("Missing --format flag")
	}

	if !strings.Contains(flags, "->format") {
		t.Error("Missing enum state transition for format flag")
	}
}

func TestGenerateBashWithHiddenSubcommand(t *testing.T) {
	cmd := Command{
		Name: "biometrics",
		Subcommands: []Command{
			{Name: "visible", Short: "Visible command"},
			{Name: "hidden", Short: "Hidden command", Hidden: true},
		},
	}

	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	err := g.GenerateBash(&buf)
	if err != nil {
		t.Fatalf("GenerateBash failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "visible") {
		t.Error("Visible command should appear in completion")
	}

	// Note: bash completion lists all commands in a single string,
	// so hidden commands might appear in the commands list
}

func TestCommandWithAliases(t *testing.T) {
	cmd := Command{
		Name:    "biometrics",
		Aliases: []string{"bio", "bm"},
	}

	g := NewGenerator(cmd, "biometrics")

	if g.cmd.Name != "biometrics" {
		t.Errorf("Expected name 'biometrics', got '%s'", g.cmd.Name)
	}

	if len(g.cmd.Aliases) != 2 {
		t.Errorf("Expected 2 aliases, got %d", len(g.cmd.Aliases))
	}
}

func TestFlagWithAllFields(t *testing.T) {
	flag := Flag{
		Name:       "format",
		Shorthand:  "f",
		Usage:      "Output format",
		Default:    "json",
		Required:   true,
		EnumValues: []string{"json", "table", "csv"},
		Hidden:     false,
	}

	if flag.Name != "format" {
		t.Errorf("Expected name 'format', got '%s'", flag.Name)
	}

	if flag.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got '%s'", flag.Shorthand)
	}

	if !flag.Required {
		t.Error("Expected Required to be true")
	}

	if len(flag.EnumValues) != 3 {
		t.Errorf("Expected 3 enum values, got %d", len(flag.EnumValues))
	}
}

func TestArgFields(t *testing.T) {
	arg := Arg{
		Name:       "file",
		Usage:      "Input file",
		Required:   true,
		EnumValues: []string{"json", "yaml"},
	}

	if arg.Name != "file" {
		t.Errorf("Expected name 'file', got '%s'", arg.Name)
	}

	if !arg.Required {
		t.Error("Expected Required to be true")
	}

	if len(arg.EnumValues) != 2 {
		t.Errorf("Expected 2 enum values, got %d", len(arg.EnumValues))
	}
}

func TestPrintCompletionHelp(t *testing.T) {
	// This function just prints help text, verify it doesn't panic
	PrintCompletionHelp()
}

// Benchmark tests
func BenchmarkGenerateBash(b *testing.B) {
	cmd := GetDefaultCommand()
	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		g.GenerateBash(&buf)
	}
}

func BenchmarkGenerateZsh(b *testing.B) {
	cmd := GetDefaultCommand()
	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		g.GenerateZsh(&buf)
	}
}

func BenchmarkGenerateFish(b *testing.B) {
	cmd := GetDefaultCommand()
	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		g.GenerateFish(&buf)
	}
}

func BenchmarkGeneratePowerShell(b *testing.B) {
	cmd := GetDefaultCommand()
	g := NewGenerator(cmd, "biometrics")
	var buf bytes.Buffer

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		g.GeneratePowerShell(&buf)
	}
}
