package healthplanet

type Data struct {
	Date    string `json:"date"`
	Keydata string `json:"keydata"`
	Model   string `json:"model"`
	Tag     string `json:"tag"`
}

type InnerscanData struct {
	BirthDate string `json:"birth_date"`
	Data      []Data `json:"data"`
	Height    string `json:"height"`
	Sex       string `json:"sex"`
}

func (i *InnerscanData) TagSearch(tag string) []Data {
	data := make([]Data, 0, 20)

	for _, val := range i.Data {
		if val.Tag == tag {
			data = append(data, val)
		}
	}
	return data
}

func (i *InnerscanData) LatestData() Data {
	data := Data{Date: "000000000000"}
	for _, val := range i.Data {
		if val.Date > data.Date {
			data = val
		}
	}

	return data
}
