select id, freeform_tags, endpoint, protocol, etag
from oci.oci_ons_subscription
where id = '{{ output.resource_id.value }}';