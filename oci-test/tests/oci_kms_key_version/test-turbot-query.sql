select title, tenant_id, region
from oci.oci_kms_key_version
where id = '{{ output.resource_id.value }}';