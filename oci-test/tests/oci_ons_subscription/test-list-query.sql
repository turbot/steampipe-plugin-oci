select id, freeform_tags, lifecycle_state, protocol, endpoint
from oci.oci_ons_subscription
where created_time = '{{ output.created_time.value }}';