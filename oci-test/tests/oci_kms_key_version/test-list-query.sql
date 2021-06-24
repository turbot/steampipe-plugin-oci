select id, key_id, lifecycle_state
from oci.oci_kms_key_version
where id = '{{ output.resource_id.value }}';