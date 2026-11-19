* **\[MDAPI-590]\[GO]** Fixed the exit status and the argument errors of the tools
  * The tools exit with status 1 after an error and print it with a newline to stderr, instead of status 0.
  * A missing or malformed option value and an unknown option print the error and the help of the tool to stderr and
    give status 2, instead of a panic or an ignored option; the value of an option is no longer taken as a positional
    argument.
  * `--quiet` is the name of the option that does not print the events; `--quite` stays as an alias.
  * perftest counts the Quote events like the other types and does not print them.
