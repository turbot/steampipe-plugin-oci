select id, display_name
from oci.oci_apigateway_api
where display_name = '{{ resourceName }}';