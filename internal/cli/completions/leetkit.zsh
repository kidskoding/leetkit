#compdef leetkit
# zsh completion for leetkit. Install with either:
#   leetkit completion zsh > "${fpath[1]}/_leetkit"
#   source <(leetkit completion zsh)

_leetkit() {
    local -a commands ids
    commands=(
        'list:list modules and exercises'
        'run:open the dashboard on one exercise'
        'check:run tests'
        'progress:show completion by module'
        'reset:restore starter code'
        'doctor:check Python, pytest, and editor setup'
        'completion:print a shell completion script'
        'help:show help'
    )
    if (( CURRENT == 2 )); then
        _describe 'command' commands
        return
    fi
    ids=(${(f)"$(leetkit list --ids 2>/dev/null)"})
    case $words[2] in
        run) _values 'exercise' $ids ;;
        check) _arguments '-v[show full pytest output]' "*:exercise:(${ids})" ;;
        reset) _arguments '(-y --yes)'{-y,--yes}'[skip the confirmation prompt]' "*:exercise:(${ids})" ;;
        list) _arguments '--ids[print only exercise ids]' ;;
        completion) _values 'shell' bash zsh fish ;;
    esac
}

if [[ $funcstack[1] == _leetkit ]]; then
    _leetkit "$@"
else
    compdef _leetkit leetkit
fi
