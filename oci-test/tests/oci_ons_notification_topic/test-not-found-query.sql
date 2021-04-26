select name, topic_id, freeform_tags, description, lifecycle_state
from oci.oci_ons_notification_topic
where topic_id = '{{ output.resource_id.value }}nf';