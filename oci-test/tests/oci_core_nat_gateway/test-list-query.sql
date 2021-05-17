select display_name, id, freeform_tags, vcn_id
from oci.oci_core_nat_gateway
where display_name = '{{ resourceName }}';