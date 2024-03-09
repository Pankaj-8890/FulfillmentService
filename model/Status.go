package model

type DeliveryStatus string

const (
    CREATED  DeliveryStatus = "CREATED"
    ASSIGNED  DeliveryStatus = "ASSIGNED"
    DELIVERED DeliveryStatus = "DELIVERED"
)

