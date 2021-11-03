select id, tag_definition_id
from oci.oci_identity_tag_default
where tag_definition_id = '{{ output.tag_definition_id.value }}';