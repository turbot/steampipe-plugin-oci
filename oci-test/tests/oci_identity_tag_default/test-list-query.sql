select id, lifecycle_state
from oci.oci_identity_tag_default
where id = '{{ output.resource_id.value }}';