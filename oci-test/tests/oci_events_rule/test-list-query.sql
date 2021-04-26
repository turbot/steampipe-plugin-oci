select display_name, id, freeform_tags, lifecycle_state
from oci.oci_events_rule
where display_name = '{{ resourceName }}';