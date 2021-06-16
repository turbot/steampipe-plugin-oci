select id, lifecycle_state
from oci.oci_identity_tag_default
where tag_definition_id = '{{ output.tag_definition_id.value }}';