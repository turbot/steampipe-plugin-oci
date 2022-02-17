select display_name, id, lifecycle_state, instance_id, vnic_id
from oci_core_vnic_attachment
where id = '{{ output.resource_id.value }}';
