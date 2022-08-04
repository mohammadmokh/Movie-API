package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Add(key string, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

func (v *Validator) Check(ok bool, key string, value string) {
	if !ok {
		v.Add(key, value)
	}
}

func Unique(s []string) bool {

	uniqueValues := make(map[string]bool)
	for _, value := range s {
		uniqueValues[value] = true
	}
	return len(s) == len(uniqueValues)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
