select title, tenant_id, region
from oci.oci_core_dhcp_options
where id = '{{ output.resource_id.value }}';