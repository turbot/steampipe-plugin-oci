select display_name, id
from oci.oci_mysql_backup
where id = '{{ output.resource_id.value }}nf';