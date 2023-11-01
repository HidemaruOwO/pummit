#! /bin/env fish

function __pummit_auto_complete
    printf '%s\t%s\n' 'alias'   ''
end

function __pummit_alias_auto_complete
    printf '%s\t%s\n' 'add'   ''
    printf '%s\t%s\n' 'list'   ''
    printf '%s\t%s\n' 'reset'   ''
    printf '%s\t%s\n' 'delete'   ''
end

complete -c pummit -x
complete -c pummit -x -l help -d 'help'

complete -c pummit -n '__fish_use_subcommand' -xa '(__pummit_auto_complete)'
complete -c pummit -n '__fish_seen_subcommand_from alias' -xa '(__pummit_alias_auto_complete)'
