select id, display_name, is_enabled, cool_down_in_seconds, freeform_tags
from oci.oci_autoscaling_auto_scaling_configuration
where display_name = '{{ resourceName }}';