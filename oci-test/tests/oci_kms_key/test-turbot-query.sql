select title, tenant_id, region
from oci.oci_kms_key
where id = '{{ output.resource_id.value }}';