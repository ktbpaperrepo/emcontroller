/*
I presume that:
The requests that need "auto schedule" will go the API functionality in the package "controller", then go to "this package", and then "this package" will call the package functionality in "models" to get the needed information and execute the scheduling schemes decided by "this package".

Thus, "controllers" will import "this package" and "models", "this package" will import "models", and "models" will not import the other 2 packages, so there will not be "Import Cycles".
*/

/*
For migration,
*/

package auto_schedule
