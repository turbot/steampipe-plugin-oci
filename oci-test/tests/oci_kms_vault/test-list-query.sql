select display_name, crypto_endpoint, id, lifecycle_state, freeform_tags
from oci.oci_kms_vault
where display_name = '{{ resourceName }}';