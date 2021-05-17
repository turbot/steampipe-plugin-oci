select display_name, id, freeform_tags, vcn_id, route_table_id, block_traffic
from oci_core_service_gateway
where display_name = '{{ resourceName }}';