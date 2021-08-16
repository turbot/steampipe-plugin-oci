select
  id,
  key_id,
  lifecycle_state
from
  oci.oci_kms_key_version
where
  id = '{{ output.resource_id.value }}aa'
  and key_id = '{{ output.key_id.value }}'
  and management_endpoint = '{{ output.endpoint.value }}'
  and region = '{{ output.region.value }}';