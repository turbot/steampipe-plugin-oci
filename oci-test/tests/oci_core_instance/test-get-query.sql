select id, display_name, availability_domain, shape
from oci.oci_core_instance
where id = '{{ output.resource_id.value }}';