select name, id, freeform_tags, description, lifecycle_state
from oci.oci_identity_dynamic_group
where id = '{{ output.resource_id.value }}aa';