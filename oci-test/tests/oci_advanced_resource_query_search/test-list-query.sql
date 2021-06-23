select identifier, display_name
from oci.oci_advanced_resource_query_search
where query = '{{ output.query.value }}' and display_name = '{{ output.resource_name.value }}';