package options

func (o *Options) Validate() []error {
	var errors []error
	errors = append(errors, o.App.Validate()...)
	errors = append(errors, o.Grpc.Validate()...)
	errors = append(errors, o.Http.Validate()...)
	return errors
}
