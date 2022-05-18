package telemetry

// Config is general configuration for instrumentation.
type Config struct {
	// BuildVersion is a SemVer string indicating
	// the version of what is being instrumented.
	BuildVersion string

	// BuildHash is the commit hash of what is being instrumented.
	BuildHash string

	// Name is the name of what is being instrumented.
	Name string
}
