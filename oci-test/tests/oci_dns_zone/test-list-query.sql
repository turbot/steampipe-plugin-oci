select name, id, lifecycle_state
from oci.oci_dns_zone
where name = '{{ resourceName }}';