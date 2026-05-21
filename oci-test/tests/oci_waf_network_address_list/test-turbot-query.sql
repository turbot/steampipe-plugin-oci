select title, tenant_id
from oci_waf_network_address_list
where id = '{{ output.resource_id.value }}';
