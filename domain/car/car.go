package car

type Car struct {
	id            string // 车牌号
	ownerOpenID   string
	ownerPhoneNum int64
}

func NewCar(id string, ownerOpenID string, ownerPhoneNum int64) *Car {
	return &Car{
		id:            id,
		ownerOpenID:   ownerOpenID,
		ownerPhoneNum: ownerPhoneNum,
	}
}

func (c *Car) ID() string {
	return c.id
}
