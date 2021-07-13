select name, id, lifecycle_state
from oci.oci_dns_tsig_key
where id = '{{ output.resource_id.value }}';