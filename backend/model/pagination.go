// Package model pagination.go 分页数据模型
package model

type StudentPageResult struct {
	List     []Student // 当前页的学生数据列表
	Total    int64     // 满足查询条件的学生总记录数
	Page     int       // 当前页码
	PageSize int       // 每页显示的记录数
}
