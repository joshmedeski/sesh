# Migration Plan: urfave/cli → Cobra + Fang

## Overview
Convert Sesh CLI from urfave/cli/v2 to Cobra with Fang enhancements while maintaining full backward compatibility and improving user experience.

## Phase 1: Dependencies & Setup

### 1.1 Update Dependencies
```bash
go mod edit -dropreplace github.com/urfave/cli/v2
go get github.com/spf13/cobra@latest
go get github.com/charmbracelet/fang@latest
```

### 1.2 Create New App Structure
- Replace `seshcli/seshcli.go` with Cobra root command
- Use Fang's `Execute()` instead of urfave's `app.Run()`
- Maintain existing dependency injection pattern

## Phase 2: Command Structure Migration

### 2.1 Root Command (`main.go` → `seshcli/root.go`)
```go
// Before (urfave)
app := seshcli.App(version)
app.Run(os.Args)

// After (Cobra + Fang)
cmd := seshcli.NewRootCommand(version)
fang.Execute(context.TODO(), cmd)
```

### 2.2 Command Mapping
| urfave/cli | Cobra Equivalent | Changes |
|------------|------------------|---------|
| `cli.Command{}` | `cobra.Command{}` | Struct field mapping |
| `Action: func(*cli.Context)` | `RunE: func(*cobra.Command, []string)` | Function signature change |
| `c.Args().Get(0)` | `args[0]` | Direct argument access |
| `c.Bool("flag")` | `cmd.Flags().GetBool("flag")` | Flag access pattern |

### 2.3 Flat Command Structure (No Changes)
- Keep all 6 commands at root level: `list`, `last`, `connect`, `clone`, `root`, `preview`
- Maintain existing aliases and descriptions

## Phase 3: Flag & Argument Migration

### 3.1 Flag Conversion Pattern
```go
// Before (urfave)
Flags: []cli.Flag{
    &cli.BoolFlag{Name: "config", Aliases: []string{"c"}},
}

// After (Cobra)
cmd.Flags().BoolP("config", "c", false, "Show configured sessions")
```

### 3.2 Command-Specific Migrations

#### List Command (8 flags)
- Convert all boolean flags to `BoolP()`
- Maintain short and long flag names
- Keep existing flag descriptions

#### Connect Command (4 flags + args)
- `--switch/-s`, `--command/-c`, `--tmuxinator/-T`, `--root/-r`
- Handle multi-word session name argument parsing
- Maintain `cobra.ExactArgs(1)` or `cobra.MinimumNArgs(1)`

#### Clone Command (2 flags + 1 arg)
- String flags: `--cmdDir/-c`, `--dir/-d`
- Repository URL validation in `PreRunE`

## Phase 4: Implementation Strategy

### 4.1 File Structure
```
seshcli/
├── root.go          # Root command + Fang setup
├── list.go          # List command (convert from urfave)
├── connect.go       # Connect command
├── clone.go         # Clone command  
├── last.go          # Last command
├── root_cmd.go      # Root session command
├── preview.go       # Preview command
└── seshcli.go       # Legacy file (remove after migration)
```

### 4.2 Dependency Injection Preservation
- Keep existing DI pattern in root command creation
- Pass dependencies through command constructors
- Maintain interface-based architecture

### 4.3 Error Handling
```go
// Before (urfave)
return cli.Exit("error message", 1)

// After (Cobra)
return fmt.Errorf("error message")  // Fang handles exit codes
```

## Phase 5: Enhanced Features with Fang

### 5.1 Automatic Enhancements
- Styled help output (automatic)
- Enhanced error messages (automatic)  
- Version flag (automatic)
- Shell completion support

### 5.2 Optional Additions
- Manpage generation capability
- Themed interface customization
- Silent usage output for scripts

## Phase 6: Testing & Validation

### 6.1 Compatibility Testing
- All existing flag combinations work
- Command aliases function correctly
- Argument parsing matches existing behavior
- Exit codes remain consistent

### 6.2 Integration Testing
- Tmux integration unchanged
- Zoxide integration unchanged
- Configuration file handling preserved
- JSON output format maintained

## Phase 7: Migration Benefits

### 7.1 User Experience Improvements
- Better help formatting and readability
- Enhanced error messages with styling
- Automatic shell completion
- Consistent CLI patterns with other Cobra tools

### 7.2 Developer Benefits
- More intuitive command structure
- Better flag handling and validation
- Improved testing patterns
- Automatic manpage generation

## Implementation Timeline
1. **Week 1**: Dependencies, root command, basic structure
2. **Week 2**: Migrate 3 simple commands (last, root, preview)
3. **Week 3**: Migrate complex commands (list, connect, clone)
4. **Week 4**: Testing, validation, documentation updates