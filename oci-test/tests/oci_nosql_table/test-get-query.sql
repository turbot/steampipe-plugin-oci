select name, id, lifecycle_state
from oci.oci_nosql_table
where id = '{{ output.resource_id.value }}';