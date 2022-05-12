select name, id
from oci.oci_limits_quota
where id = '{{ output.resource_id.value }}nf';