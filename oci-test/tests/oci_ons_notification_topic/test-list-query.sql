select name, topic_id, freeform_tags, lifecycle_state, description
from oci.oci_ons_notification_topic
where name = '{{ resourceName }}';