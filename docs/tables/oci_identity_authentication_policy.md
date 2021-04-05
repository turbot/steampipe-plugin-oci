# Table: oci_identity_authentication_policy

The authentication policy for local IAM users and network source restrictions for all users in tenancy. For more information about using a authentication policy, go to [Managing Authentication Settings](https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/managingpasswordrules.htm).

## Examples

### Basic info

```sql
select
  minimum_password_length,
  is_lowercase_characters_required,
  is_numeric_characters_required,
  is_special_characters_required,
  is_uppercase_characters_required,
  is_username_containment_allowed,
  network_source_ids
from
  oci_identity_authentication_policy
```
