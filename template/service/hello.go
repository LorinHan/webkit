package service

var HelloSvc = &HelloService{}

type HelloService struct {
}

type Company struct {
	ID string `json:"id"`
}

func (h *HelloService) SayHi() (string, error) {
	// return "", errors.New("test normal error")
	// return "", errork.New("test normal error")
	// var count int64
	// if err := model.DB().Raw("select count(id) from companies").Scan(&count).Error; err != nil {
	// 	return "", err
	// }
	// if err := model.DB().Exec("insert into companies values(1, 2)").Error; err != nil {
	// 	return "", err
	// }
	// var comp map[string]interface{}
	// if err := model.DB().Model(Company{}).Limit(1).First(&comp).Error; err != nil {
	// 	return "", err
	// }
	// fmt.Println(comp)
	// zap.S().Info(count)
	return "Hello", nil
}
