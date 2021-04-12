select display_name, id, freeform_tags, vcn_id, lifecycle_state
from oci.oci_core_route_table
where display_name = '{{ resourceName }}';