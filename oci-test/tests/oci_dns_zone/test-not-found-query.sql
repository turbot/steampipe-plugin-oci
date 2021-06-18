select name, id
from oci.oci_dns_zone
where id = '{{ output.resource_id.value }}aa';