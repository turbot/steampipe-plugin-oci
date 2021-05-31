select id, display_name
from oci.oci_apigateway_api
where id = '{{ output.resource_id.value }}';