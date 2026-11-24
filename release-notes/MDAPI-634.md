* **\[MDAPI-634]\[GO]** `tools connect -t` and `tools dump -t` write the tape without `-p dxfeed.wildcard.enable=true`,
  as the tools of the dxFeed C++ API do; `-p dxfeed.wildcard.enable=false` still turns the wildcard off.
