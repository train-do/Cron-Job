package domain

type Banner struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	PathPage  string
	StartDate string `gorm:"type:date"`
	EndDate   string `gorm:"type:date"`
	IsPublish bool
	ImageUrl  string
}

func BannerSeed() []Banner {
	return []Banner{
		{
			Title:     "Promo Akhir Tahun",
			PathPage:  "www.lumoshive.com",
			StartDate: "2024-11-09",
			EndDate:   "2024-11-12",
			IsPublish: false,
			ImageUrl:  "https://ik.imagekit.io/lumoshiveAcademy/agency-api/profile_pictures/abc_1733426377229_whatb7wXM.jpg?updatedAt=1733426379031",
		},
		{
			Title:     "Produk Baru",
			PathPage:  "www.lumoshive.com",
			StartDate: "2024-11-08",
			EndDate:   "2024-11-11",
			IsPublish: true,
			ImageUrl:  "https://ik.imagekit.io/lumoshiveAcademy/agency-api/profile_pictures/abc_1733426377229_whatb7wXM.jpg?updatedAt=1733426379031",
		},
		{
			Title:     "Diskon 30%",
			PathPage:  "www.lumoshive.com",
			StartDate: "2024-11-07",
			EndDate:   "2024-11-10",
			IsPublish: false,
			ImageUrl:  "https://ik.imagekit.io/lumoshiveAcademy/agency-api/profile_pictures/abc_1733426377229_whatb7wXM.jpg?updatedAt=1733426379031",
		},
		{
			Title:     "Giveaway",
			PathPage:  "www.lumoshive.com",
			StartDate: "2024-11-03",
			EndDate:   "2024-11-09",
			IsPublish: false,
			ImageUrl:  "https://ik.imagekit.io/lumoshiveAcademy/agency-api/profile_pictures/abc_1733426377229_whatb7wXM.jpg?updatedAt=1733426379031",
		},
	}
}
