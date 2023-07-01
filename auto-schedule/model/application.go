package model

type Application struct {
	Name         string       `json:"name"`
	Resources    AppResources `json:"resources"`    // The resources information of this application
	Dependencies []Dependency `json:"dependencies"` // The information of all applications that this application depends on.
}

// In Dependency, only AppName is enough, because:
// 1. Bandwidth is not considered in this model;
// 2. RTT is not a hard requirement, but soft, which means that the smaller RTT the better, but high RTT is also OK.
// A high RTT will only make the response slow, but the application will still work.
type Dependency struct {
	AppName string `json:"appName"` // the name of the dependent application
}
