select name, id,lifecycle_state
from oci.oci_core_local_peering_gateway
where name = '{{ resourceName }}';