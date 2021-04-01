select title, tenant_id
from oci.oci_core_virtual_circuit
where id = '{{ output.resource_id.value }}';