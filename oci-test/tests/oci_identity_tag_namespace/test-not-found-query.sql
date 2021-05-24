select name, id, description, lifecycle_state
from oci.oci_identity_tag_namespace
where id = '{{ output.resource_id.value }}aa';