package app

type CliOptions interface {
	Flags() (fss FlagSets)
	Validate() []error
}

type CompleteableOptions interface {
	Complete() error
}