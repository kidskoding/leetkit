# bash completion for leetkit. Install with:
#   leetkit completion bash > ~/.local/share/bash-completion/completions/leetkit

_leetkit() {
    local cur=${COMP_WORDS[COMP_CWORD]}
    if [[ $COMP_CWORD -eq 1 ]]; then
        COMPREPLY=($(compgen -W "list run check progress reset doctor completion version help" -- "$cur"))
        return
    fi
    case ${COMP_WORDS[1]} in
        run) COMPREPLY=($(compgen -W "$(leetkit list --ids 2>/dev/null)" -- "$cur")) ;;
        check) COMPREPLY=($(compgen -W "-v $(leetkit list --ids 2>/dev/null)" -- "$cur")) ;;
        reset) COMPREPLY=($(compgen -W "-y --yes $(leetkit list --ids 2>/dev/null)" -- "$cur")) ;;
        list) COMPREPLY=($(compgen -W "--ids" -- "$cur")) ;;
        completion) COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur")) ;;
    esac
}

complete -F _leetkit leetkit
