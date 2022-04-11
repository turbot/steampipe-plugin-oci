connection "oci" {
  plugin = "oci"
  # config_file_profile = "DEFAULT"
  # config_file_path = "~/.oci/config"

  # The maximum number of attempts (including the initial call) Steampipe will
  # make for failing API calls. Defaults to 9 and must be greater than or equal to 1.
  #max_error_retry_attempts = 9

  # The minimum retry delay in milliseconds after which retries will be performed.
  # This delay is also used as a base value when calculating the exponential backoff retry times.
  # Defaults to 25ms and must be greater than or equal to 1ms.
  #min_error_retry_delay = 25
}
