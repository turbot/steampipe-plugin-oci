select id
from oci.oci_core_instance
where id = '{{ output.resource_id.value }}::dummy';