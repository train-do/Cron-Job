package helper

import (
	"fmt"
	"project/config"
	"project/database"
	"project/domain"
	"reflect"

	"github.com/xuri/excelize/v2"
)

func CronExcel(migrateDb bool, seedDb bool) func() {
	return func() {
		appConfig, err := config.LoadConfig(migrateDb, seedDb)
		if err != nil {
			fmt.Println(err)
			return
		}
		// instance database
		db, err := database.ConnectDB(appConfig)
		if err != nil {
			fmt.Println(err)
			return
		}
		var banners []domain.Banner
		if err := db.Find(&banners).Error; err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Printf("%+v\n", banners)
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		var data [][]interface{} = [][]interface{}{}
		for _, banner := range banners {
			var temp []interface{}
			val := reflect.ValueOf(banner)
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				// fmt.Printf("%v ", field)
				temp = append(temp, field)
			}
			// fmt.Println("\n", temp)
			data = append(data, temp)
		}
		// fmt.Println(data)
		for i, row := range data {
			cell, err := excelize.CoordinatesToCellName(1, i+1)
			if err != nil {
				fmt.Println(err)
				return
			}
			f.SetSheetRow("Sheet1", cell, &row)
		}
		if err := f.SaveAs("Book1.xlsx"); err != nil {
			fmt.Println(err)
		}
		// fmt.Println(data)
	}
}
