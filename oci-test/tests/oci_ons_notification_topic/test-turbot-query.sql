select title, tenant_id
from oci.oci_ons_notification_topic
where topic_id = '{{ output.resource_id.value }}';