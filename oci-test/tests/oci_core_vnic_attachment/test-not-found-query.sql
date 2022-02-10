select display_name, id, lifecycle_state, instance_id, vnic_id
from oci_core_vnic_attachment
where display_name = '{{ output.resource_name.value }}nf';
