select display_name, id, description, lifecycle_state, freeform_tags, is_enabled
from oci.oci_mysql_channel
where id = '{{ output.resource_id.value }}';