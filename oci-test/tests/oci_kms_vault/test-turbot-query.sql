select title, tenant_id, region
from oci.oci_kms_vault
where id = '{{ output.resource_id.value }}';