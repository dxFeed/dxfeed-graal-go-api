* **\[MDAPI-580]\[GO]** Fixed the end of a program by `SIGPIPE` (exit status 141) on Linux when the peer of a
  connection closed it while the SDK wrote to it, for example a publisher whose client disconnected at once.
  * The API registers `SIGPIPE` and `SIGXFSZ` with `os/signal` on Linux and macOS, so a write of the program to a
    closed standard output returns `EPIPE` instead of ending the program (docs/configuration.md, "Signals").
