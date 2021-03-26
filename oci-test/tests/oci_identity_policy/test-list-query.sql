select name, id, freeform_tags, description, lifecycle_state
from oci.oci_identity_policy
where name = '{{ resourceName }}';