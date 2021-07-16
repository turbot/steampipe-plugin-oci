select
  id,
  display_name,
  lifecycle_state,
  scope,
  lifetime,
  freeform_tags,
  compartment_id,
  tenant_id
from
  oci.oci_core_public_ip
where
  id = '{{ output.resource_id.value }}';