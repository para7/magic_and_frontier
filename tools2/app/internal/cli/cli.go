package cli

import (
	"flag"
	"fmt"
	"io"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
)

func Run(args []string, stdout, stderr io.Writer, cfg config.Config) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "validate":
		return runValidate(args[1:], stdout, stderr, cfg)
	case "export":
		return runExport(args[1:], stdout, stderr, cfg)
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func runValidate(args []string, stdout, stderr io.Writer, cfg config.Config) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc := application.NewService(cfg, application.DefaultDependencies(cfg))
	report, err := svc.ValidateAll()
	if err != nil {
		fmt.Fprintf(stderr, "validate failed: %v\n", err)
		return 1
	}
	if !report.OK {
		fmt.Fprintln(stderr, "savedata validation failed:")
		fmt.Fprintln(stderr, report.String())
		return 1
	}

	fmt.Fprintf(stdout, "savedata validation ok: items=%d grimoire=%d skills=%d enemy_skills=%d enemies=%d treasures=%d loottables=%d\n",
		report.Counts.Items,
		report.Counts.Grimoire,
		report.Counts.Skills,
		report.Counts.EnemySkills,
		report.Counts.Enemies,
		report.Counts.Treasures,
		report.Counts.LootTables,
	)
	return 0
}

func runExport(args []string, stdout, stderr io.Writer, cfg config.Config) int {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc := application.NewService(cfg, application.DefaultDependencies(cfg))
	result := svc.ExportDatapack()
	if !result.OK {
		fmt.Fprintf(stderr, "export failed: %s\n", result.Message)
		if result.Details != "" {
			fmt.Fprintln(stderr, result.Details)
		}
		return 1
	}

	fmt.Fprintf(stdout, "datapack export completed: output=%s total_files=%d item_functions=%d item_loot=%d spell_functions=%d spell_loot=%d skills=%d enemy_skills=%d enemy_functions=%d enemy_loot=%d treasure_loot=%d loottable_loot=%d\n",
		result.OutputRoot,
		result.Generated.TotalFiles,
		result.Generated.ItemFunctions,
		result.Generated.ItemLootTables,
		result.Generated.SpellFunctions,
		result.Generated.SpellLootTables,
		result.Generated.SkillFunctions,
		result.Generated.EnemySkillFunctions,
		result.Generated.EnemyFunctions,
		result.Generated.EnemyLootTables,
		result.Generated.TreasureLootTables,
		result.Generated.LoottableLootTables,
	)
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: tools2-cli <command>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "commands:")
	fmt.Fprintln(w, "  validate   validate savedata and export settings")
	fmt.Fprintln(w, "  export     validate and export datapack")
}
