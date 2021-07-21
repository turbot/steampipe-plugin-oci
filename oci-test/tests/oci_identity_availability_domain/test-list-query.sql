select name, id
from oci.oci_identity_availability_domain
where name = '{{ output.resource_name.value }}';