select display_name
from oci.oci_advanced_resource_query_search
where query = '{{ output.query.value }}' and identifier = '{{ output.resource_id.value }}aa';