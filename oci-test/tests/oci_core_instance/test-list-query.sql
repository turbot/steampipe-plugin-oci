select id, display_name, availability_domain, shape
from oci.oci_core_instance
where display_name = '{{ output.resource_name.value }}';