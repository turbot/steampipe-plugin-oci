select name, id, lifecycle_state
from oci.oci_kms_key
where name = '{{ resourceName }}';