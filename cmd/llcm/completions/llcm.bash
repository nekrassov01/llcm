#! /bin/bash

_cli_init_completion() {
  COMPREPLY=()
  _get_comp_words_by_ref "$@" cur prev words cword
}

_llcm() {
  [[ "${COMP_WORDS[0]}" == "source" ]] && return 0
  local cur words opts comp
  if declare -F _init_completion >/dev/null 2>&1; then
    _init_completion -n "=:" || return
  else
    _cli_init_completion -n "=:" || return
  fi
  cur="${COMP_WORDS[COMP_CWORD]}"
  if [[ "$cur" == "-"* ]]; then
    comp="${words[*]} ${cur} --generate-bash-completion"
  else
    comp="${words[*]} --generate-bash-completion"
  fi
  opts=$(eval "${comp}" 2>/dev/null)
  # shellcheck disable=SC2207
  COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
}

command -v llcm >/dev/null 2>&1 && complete -o bashdefault -o default -o nospace -F _llcm llcm
