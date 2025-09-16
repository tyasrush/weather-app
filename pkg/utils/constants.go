package utils

import "time"

const DefaultDBTimeout time.Duration = 5 * time.Second

const (
	OrderByCreatedAtAsc  = "created_at_ascend"
	OrderByCreatedAtDesc = "created_at_descend"
	OrderByNameAsc       = "name_ascend"
	OrderByNameDesc      = "name_descend"
)

var OrderByMaps = map[string]bool{
	OrderByCreatedAtAsc:  true,
	OrderByCreatedAtDesc: true,
	OrderByNameAsc:       true,
	OrderByNameDesc:      true,
}
