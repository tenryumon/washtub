package helper

type Criteria map[string]interface{}

var validSortType = map[string]string{
	"ASC":  "ASC",
	"DESC": "DESC",
}

var validLimit = map[int]bool{
	10:  true,
	20:  true,
	30:  true,
	50:  true,
	100: true,
}

func (c Criteria) SetSortBy(sortBy, sortType string, valid map[string]string) {
	c["sort_by"] = "id"
	if value, ok := valid[sortBy]; ok {
		c["sort_by"] = value
	}

	c["sort_type"] = "ASC"
	if value, ok := validSortType[sortType]; ok {
		c["sort_type"] = value
	}
}

func (c Criteria) SetLimit(page int64, limit int) {
	if _, ok := validLimit[limit]; ok {
		c["limit"] = limit
	} else {
		c["limit"] = 10
		limit = 10
	}

	// page start from 1 on client UI
	c["offset"] = int64(0)
	if page > 0 {
		c["offset"] = int64(limit) * (page - 1)
	}
}
