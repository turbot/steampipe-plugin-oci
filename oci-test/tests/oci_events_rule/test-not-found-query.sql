select id, display_name, lifecycle_state, freeform_tags
from oci.oci_events_rule
where id = '{{ output.resource_id.value }}aa';