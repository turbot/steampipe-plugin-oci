select display_name, id, freeform_tags, time_created, lifecycle_state
from oci.oci_core_nat_gateway
where id = '{{ output.resource_id.value }}aa';