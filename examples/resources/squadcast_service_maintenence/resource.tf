resource "squadcast_service_maintenance" "test" {
  service_id = "service_id"

  windows {
    from             = "2032-06-01T10:30:00.000Z"
    till             = "2032-06-01T11:30:00.000Z"
    repeat_till      = "2032-06-30T10:30:00.000Z"
    repeat_frequency = "week"
  }

  windows {
    from = "2032-07-01T10:30:00.000Z"
    till = "2032-07-02T10:30:00.000Z"
  }
}