select display_name, id, freeform_tags, description, lifecycle_state
from oci.oci_identity_domain
where id = '{{ output.resource_id.value }}aa';