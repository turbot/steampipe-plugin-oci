# Table: oci_application_migration_source

The properties that define a source. Source refers to the source environment from which you migrate an application to Oracle Cloud Infrastructure.

## Examples

### Basic info

```sql
select
    id,
    display_name,
    description,
    lifecycle_details,
    lifecycle_state as state
from
oci_application_migration_source;
```