select title, tenant_id
from oci_core_vnic_attachment
where id = '{{ output.resource_id.value }}';
