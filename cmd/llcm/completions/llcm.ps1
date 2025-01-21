Register-ArgumentCompleter -Native -CommandName "llcm" -ScriptBlock {
    param($commandName, $wordToComplete, $cursorPosition)
    (Invoke-Expression "$wordToComplete --generate-bash-completion").ForEach{
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
