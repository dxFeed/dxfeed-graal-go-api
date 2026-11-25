* **\[MDAPI-636]\[GO]** `tools nettest` listens on an address with connector options, such as `:7400[bindAddr=127.0.0.1]`;
  as in the NetTest tool of the SDK, it waits up to 3 s for the port to be bound and still exits with status 1 when the
  port is in use.
