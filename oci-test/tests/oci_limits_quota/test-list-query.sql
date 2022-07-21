select name, id
from oci.oci_limits_quota
where name = '{{ resourceName }}';