select display_name, id, lifecycle_state, freeform_tags, cidr_block
from oci.oci_core_vcn
where display_name = '{{ resourceName }}';