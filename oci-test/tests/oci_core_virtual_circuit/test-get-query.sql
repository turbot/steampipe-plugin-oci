select display_name, id, freeform_tags, type, bandwidth_shape_name, provider_name, provider_state
from oci.oci_core_virtual_circuit
where id = '{{ output.resource_id.value }}';