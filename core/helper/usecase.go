package helper

func NormalizePageLimit(count int64, param *map[string]interface{}) (page int64, limit int) {
	page = 1
	limit = -1
	offset := int64(0)
	if offsetParam, ok := (*param)["offset"]; ok {
		offset = offsetParam.(int64)
	}

	if limitParam, ok := (*param)["limit"]; ok {
		limit = limitParam.(int)
	} else {
		// limit = -1
		// it mean search will has no limit
		return
	}

	page = (offset / int64(limit)) + 1
	maxPage := count/int64(limit) + 1
	if page > maxPage {
		page = maxPage
		// edit offset value in param
		(*param)["offset"] = (page - 1) * int64(limit)
	}
	return
}
