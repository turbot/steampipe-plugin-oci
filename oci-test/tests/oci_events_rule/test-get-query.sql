select id, display_name, freeform_tags, lifecycle_state, is_enabled, description, actions
from oci.oci_events_rule
where id = '{{ output.resource_id.value }}';