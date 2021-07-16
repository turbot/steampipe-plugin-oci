select
  id,
  display_name
from
  oci.oci_core_public_ip
where
  display_name = '{{ resourceName }}';