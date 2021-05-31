select title, tenant_id
from oci.oci_apigateway_api
where id = '{{ output.resource_id.value }}';