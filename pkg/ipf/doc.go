// Package ipf reads and writes instrument profiles, the descriptions of instruments, in the Instrument Profile Format
// (IPF).
//
// An [InstrumentProfileReader] reads profiles from a file or a URL, with a user and a password or an authorization
// token; [WriteToFile] writes them. The profiles are [events.InstrumentProfile] values. Their standard fields have
// getters and setters; [GetField], [SetField] and the numeric and date variants read and write any field by name in
// the IPF text form, the custom fields included, and [InstrumentProfileField] lists the standard fields.
package ipf
