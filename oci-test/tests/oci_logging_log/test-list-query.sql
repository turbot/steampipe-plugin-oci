select name, id
from oci.oci_logging_log
where name = '{{ resourceName }}';