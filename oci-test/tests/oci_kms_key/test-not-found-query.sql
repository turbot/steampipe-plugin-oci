select name, id
from oci.oci_kms_key
where id = '{{ output.resource_id.value }}nf';