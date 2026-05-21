select display_name, id
from oci_waf_web_app_firewall_policy
where id = '{{ output.resource_id.value }}nf';
