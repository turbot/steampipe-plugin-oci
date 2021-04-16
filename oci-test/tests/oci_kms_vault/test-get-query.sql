select display_name, id, management_endpoint, vault_type, lifecycle_state
from oci.oci_kms_vault
where id = '{{ output.resource_id.value }}';