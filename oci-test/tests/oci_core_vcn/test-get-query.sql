select display_name, id, cidr_block, lifecycle_state
from oci.oci_core_vcn
where id = '{{ output.resource_id.value }}';