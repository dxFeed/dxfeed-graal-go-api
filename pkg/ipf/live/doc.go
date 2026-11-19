// Package live keeps instrument profiles up to date: an [InstrumentProfileConnection] reads the profiles from a file,
// a URL or an instrument profile server into an [InstrumentProfileCollector], which notifies its listeners about the
// new, changed and removed profiles.
//
//	collector, _ := live.NewInstrumentProfileCollector()
//	connection, _ := live.CreateConnection("https://example.com/profiles.ipf.zip", collector)
//	_ = collector.AddUpdateListener(listener)
//	_ = connection.Start()
//	completed, _ := connection.WaitUntilCompleted(time.Minute)
//
// Close the connection first, then the collector.
package live
