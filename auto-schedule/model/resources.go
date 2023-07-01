/*
There are 2 types of resources: (If no people defined them previously, we can define them here.)
1. soft resource: we can choose to meet the requirement or not, if we do not meet the requirement, the response time will be longer, but the application will still work, such as CPU and RTT.
2. hard resource: If we accept an application, we have to meet its hard requirements, or otherwise this application will not work, such as Memory and Storage (If memory is not enough, there will be OOM (Out-of memory) error).
*/

package model

type AppResources struct {
	GenericResources `json:",inline"`

	// number of CPU logical cores allocated to this application, only useful when the application is deployed. We have this resource because CPU core is a soft requirement, which means we do not have to allocate all required CPU cores to an application.
	AllocatedCpuCore float64 `json:"cpuCore"`
}

type GenericResources struct {
	CpuCore float64 `json:"cpuCore"` // number of CPU logical cores that this application needs, this is a soft requirement
	Memory  float64 `json:"memory"`  // unit Byte (B), this is a hard requirement
	Storage float64 `json:"storage"` // unit Byte (B), this is a hard requirement
}
