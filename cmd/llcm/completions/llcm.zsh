#compdef llcm

_llcm() {
  local -a opts
  local cur
  cur="${words[-1]}"

  if [[ "$cur" == "-"* ]]; then
    opts=($(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:${#words[@]}-1} ${cur} --generate-bash-completion))
  else
    opts=($(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:${#words[@]}-1} --generate-bash-completion))
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

command -v llcm >/dev/null 2>&1 && compdef _llcm llcm
