package friends

type FriendsResponse struct {
	Name       string `json:"name"`
	Subtotal   int64  `json:"subtotal"`
	TaxPortion int32  `json:"tax_portion"`
	TotalOwed  int32  `json:"total_owed"`
}
