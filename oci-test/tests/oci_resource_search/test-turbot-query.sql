select title, tenant_id
from oci.oci_resource_search
where query = '{{ output.query.value }}' and identifier = '{{ output.resource_id.value }}' and search_region = '{{ output.region.value }}';