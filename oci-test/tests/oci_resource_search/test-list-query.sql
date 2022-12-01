select identifier, display_name
from oci.oci_resource_search
where query = '{{ output.query.value }}' and display_name = '{{ output.resource_name.value }}' and search_region = '{{ output.region.value }}';