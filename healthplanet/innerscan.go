package healthplanet

import "fmt"

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

func (i *InnerscanData) TagSearch(tag string) InnerscanData {
	inner := InnerscanData{
		BirthDate: i.BirthDate,
		Height:    i.Height,
		Sex:       i.Sex,
	}
	data := make([]Data, 0, 20)

	for _, val := range i.Data {
		fmt.Println("tag search ...")
		fmt.Printf("%+v\n", val)
		if val.Tag == tag {
			data = append(data, val)
		}
	}
	inner.Data = data
	return inner
}

func (i *InnerscanData) NewerData(date string) InnerscanData {
	inner := InnerscanData{
		BirthDate: i.BirthDate,
		Height:    i.Height,
		Sex:       i.Sex,
	}
	data := make([]Data, 0, 20)

	for _, val := range i.Data {
		if val.Date > date {
			data = append(data, val)
		}
	}
	inner.Data = data
	return inner
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
