select display_name, id, freeform_tags, cidr_block, lifecycle_state
from oci.oci_core_subnet
where display_name = '{{ resourceName }}';