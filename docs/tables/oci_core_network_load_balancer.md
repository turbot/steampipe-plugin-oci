# Table: oci_core_network_load_balancer

Network Load Balancers can provide automated traffic distribution from one entry point to multiple servers in a backend set. It ensure that your services remain available by directing traffic only to healthy servers. Network Load Balancer provides the benefits of flow high availability, source and destination IP addresses, and port preservation. It is designed to handle volatile traffic patterns and millions of flows, offering high throughput while maintaining ultra low latency. It is the ideal load balancing solution for latency sensitive workloads. It is also optimized for long-running connections in the order of days or months.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  subnet_id,
  lifecycle_state as state,
  health_status,
  time_created
from
  oci_core_network_load_balancer;
```

### List NLBs assigns with public IP address

```sql
select
  display_name,
  id,
  is_private
from
  oci_core_network_load_balancer
where
  not is_private;
```

### List critical NLBs

```sql
select
  display_name,
  id,
  network_load_balancer_health -> 'status' as health_status
from
  oci_core_network_load_balancer
where
  network_load_balancer_health ->> 'status' = 'CRITICAL';
```
