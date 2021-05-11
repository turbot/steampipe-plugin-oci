select display_name, id, freeform_tags
from oci.oci_logging_log_group
where display_name = '{{ output.display_name.value }}::dummy';