select
  id,
  display_name
from
  oci.oci_core_public_ip
where
  id = '{{ output.resource_id.value }}::dummy';