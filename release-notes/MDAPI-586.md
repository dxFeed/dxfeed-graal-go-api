* **\[MDAPI-586]\[GO]** Seven samples accept `-h` and their address as a flag; Ctrl+C ends each sample with status 0
  * auth, feed-ipf-connect, ondemand-replay, print-candles, print-indexed-events, print-quote-events and quick-start
    get `-address`, feed-ipf-connect also `-ipf` (its symbols come from the instrument profiles), and the others except
    quick-start `-symbol`, with the former values as the defaults; `-h` prints their usage instead of starting them.
  * ondemand-replay takes the credentials of the on-demand service with `-user` and `-password` (default `demo`); the
    on-demand service needs the credentials of an account with access to it.
  * feed-ipf-connect, ondemand-replay, print-candles, print-indexed-events and print-quote-events close the endpoint on
    Ctrl+C and exit with status 0 instead of ending with status 130.
