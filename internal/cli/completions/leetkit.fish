# fish completion for leetkit. Install with:
#   leetkit completion fish > ~/.config/fish/completions/leetkit.fish

set -l commands list run check progress reset doctor completion version help

complete -c leetkit -f
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a list -d "List modules and exercises"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a run -d "Open the dashboard on one exercise"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a check -d "Run tests"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a progress -d "Show completion by module"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a reset -d "Restore starter code"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a doctor -d "Check Python, pytest, and editor setup"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a completion -d "Print a shell completion script"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a version -d "Print the version"
complete -c leetkit -n "not __fish_seen_subcommand_from $commands" -a help -d "Show help"

complete -c leetkit -n "__fish_seen_subcommand_from run check reset" -a "(leetkit list --ids 2>/dev/null)" -d exercise
complete -c leetkit -n "__fish_seen_subcommand_from check" -s v -d "Show full pytest output"
complete -c leetkit -n "__fish_seen_subcommand_from reset" -s y -l yes -d "Skip the confirmation prompt"
complete -c leetkit -n "__fish_seen_subcommand_from list" -l ids -d "Print only exercise ids"
complete -c leetkit -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"
