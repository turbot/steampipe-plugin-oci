select name, id, lifecycle_state
from oci.oci_dns_zone
where id = '{{ output.resource_id.value }}';