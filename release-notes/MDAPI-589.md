* **\[MDAPI-589]\[GO]** The help of the tools and the examples of the documentation match the code
  * The help of connect, dump, perftest and latencytest lists the options they accept (`-f`, `-t`, `-p`, `-q`,
    `--force-stream`, `--ignore-exchanges`), the event types and `all`, and no longer shows placeholders and doubled
    quotes; dump no longer claims a wildcard by default or a short list of event types.
  * Tests type-check the Go code blocks of the Markdown files and of the doc comments, and check the flags of the
    command lines of the samples against their `-h`.
