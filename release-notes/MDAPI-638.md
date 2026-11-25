* **\[MDAPI-638]\[GO]** Faster publishing: `DXPublisher.Publish` builds the event mappers once and puts a batch with its
  strings into a few blocks of C memory instead of four allocations per event, so the NetTest producer publishes about
  1.25M Trade events/s instead of 0.56M on one core; a Quote copies its symbol once on receive. Benchmarks of publishing
  and receiving run in CI
