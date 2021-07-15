select display_name, id, lifecycle_state, freeform_tags, is_enabled
from oci.oci_mysql_channel
where display_name = '{{ resourceName }}';