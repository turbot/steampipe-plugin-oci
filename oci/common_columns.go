package oci

const (
	// Constants for Standard Column Descriptions
	ColumnDescriptionAkas        = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags        = "A map of tags for the resource."
	ColumnDescriptionTitle       = "Title of the resource."
	ColumnDescriptionTenant      = "The OCID of the Tenant in which the resource is located."
	ColumnDescriptionCompartment = "The OCID of the compartment in Tenant in which the resource is located."
	ColumnDescriptionRegion      = "The OCI region in which the resource is located."

	// Other repetitive columns for the provider
	ColumnDescriptionFreefromTags = "Free-form tags for resource. This tags can be applied by any user with permissions on the resource."
	ColumnDescriptionDefinedTags  = "Defined tags for resource. Defined tags are set up in your tenancy by an administrator. Only users granted permission to work with the defined tags can apply them to resources."
)
