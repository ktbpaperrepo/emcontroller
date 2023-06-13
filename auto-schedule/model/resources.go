package model

type CloudResources struct {
	CpuCore float64 `json:"cpuCore"` // number of logical CPU cores
	Memory  float64 `json:"memory"`  // memory size unit Byte (B)
	Storage float64 `json:"storage"` // storage size unit Byte (B)
	Vm      float64 `json:"vm"`      // number of virtual machines, negative values, such as -1, means unlimited
	Volume  float64 `json:"volume"`  // number of volumes, negative values, such as -1, means unlimited
	Port    float64 `json:"port"`    // number of network ports, negative values, such as -1, means unlimited
}

type AppResources struct {
	CpuCore float64 `json:"cpuCore"` // number of CPU logical cores that this application needs
	Memory  float64 `json:"memory"`  // unit Byte (B)
	Storage float64 `json:"storage"` // unit Byte (B)
}

type Dependency struct {
	AppIdx string `json:"appIdx"` // the index of the dependent application
}
