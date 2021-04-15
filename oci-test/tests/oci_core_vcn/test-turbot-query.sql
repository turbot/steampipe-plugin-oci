select title, tenant_id
from oci.oci_core_vcn
where id = '{{ output.resource_id.value }}';