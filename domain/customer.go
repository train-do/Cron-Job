package domain

type Customer struct {
	ID      uint   `gorm:"primaryKey" json:"-"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func CustomerSeed() []Customer {
	return []Customer{
		{Name: "Customer Satu", Address: "Alamat Satu"},
		{Name: "Customer Dua", Address: "Alamat Dua"},
		{Name: "Customer Tiga", Address: "Alamat Tiga"},
		{Name: "Customer Empat", Address: "Alamat Empat"},
		{Name: "Customer Lima", Address: "Alamat Lima"},
		{Name: "Customer Enam", Address: "Alamat Enam"},
		{Name: "Customer Tujuh", Address: "Alamat Tujuh"},
		{Name: "Customer Delapan", Address: "Alamat Delapan"},
		{Name: "Customer Sembilan", Address: "Alamat Sembilan"},
		{Name: "Customer Sepuluh", Address: "Alamat Sepuluh"},
	}
}
