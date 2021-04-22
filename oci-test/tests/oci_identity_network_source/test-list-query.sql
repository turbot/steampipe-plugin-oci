select name, id, freeform_tags, description, lifecycle_state, public_source_list
from oci.oci_identity_network_source
where name = '{{ resourceName }}';