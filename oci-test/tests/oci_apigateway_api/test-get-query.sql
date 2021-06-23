select id, display_name, freeform_tags, lifecycle_state
from oci.oci_apigateway_api
where id = '{{ output.resource_id.value }}';