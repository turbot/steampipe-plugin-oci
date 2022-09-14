select name, id
from oci.oci_streaming_stream
where id = '{{ output.resource_id.value }}aa';