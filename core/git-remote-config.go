package core

type GitRemoteConfig map[string]string

// Set if value and name is set
func (grc GitRemoteConfig) Set(name, value string) {
	if value == "" || name == "" {
		return
	}

	grc[name] = value
}

// SetIfNotContains set value if the value is not one on the listed ifNot
func (grc GitRemoteConfig) SetIfNotContains(name, value string, ifNot ...string) {
	for _, v := range ifNot {
		if value == v {
			return
		}

	}

	grc.Set(name, value)
}

// SetIf set value if the value is not one on the listed ifNot
func (grc GitRemoteConfig) SetIf(name, value string, ifCase bool) {
	if !ifCase {
		return
	}

	grc.Set(name, value)
}

// Get if value and name is set
func (grc GitRemoteConfig) Get(name string, defaultValues ...string) (value string) {
	if name == "" {
		return
	}

	if v, found := grc[name]; found && v != "" {
		return v
	}

	for _, v := range defaultValues {
		if v != "" {
			return v
		}
	}

	return
}
