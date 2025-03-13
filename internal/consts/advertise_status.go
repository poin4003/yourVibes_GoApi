package consts

type AdvertiseStatus int16

const (
	NOT_ADVERTISE AdvertiseStatus = 0
	IS_ADVERTISE  AdvertiseStatus = 1
	WAS_ADVERTISE AdvertiseStatus = 2
)

var AdvertiseStatusList = []interface{}{
	NOT_ADVERTISE,
	IS_ADVERTISE,
	WAS_ADVERTISE,
}
