select display_name, id, lifecycle_state
from oci_waf_web_app_firewall_policy
where display_name = '{{ output.resource_name.value }}';
