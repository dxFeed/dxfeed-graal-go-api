* **\[MDAPI-563]\[GO]** Go samples for the remaining samples of the C++ API.
  * Tape file samples: `write-tape-file` publishes quotes into a tape file, `file-parser` prints the events of a
    tape file, `convert-tape-file` converts one tape file into another. `publish-profiles` serves Profile events of
    `:TEST` symbols, and `request-profile` requests one of them. The tests run these samples, without network, on
    every CI platform.
