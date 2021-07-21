select title, tenant_id
from oci.oci_core_volume_attachment
where id = '{{ output.resource_id.value }}';