select name, id, freeform_tags, lifecycle_state
from oci.oci_core_local_peering_gateway
where id = '{{ output.resource_id.value }}aa';