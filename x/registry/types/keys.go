package types

const (
	// ModuleName is the name of the staking module
	ModuleName = "registry"

	// RecordRentModuleAccountName is the name of the module account that keeps track of record rents paid.
	RecordRentModuleAccountName = "record_rent"

	// AuthorityRentModuleAccountName is the name of the module account that keeps track of authority rents paid.
	AuthorityRentModuleAccountName = "authority_rent"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)
