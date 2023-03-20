package render

import "github.com/jedib0t/go-pretty/v6/table"

var balanceTableHeader = table.Row{"账号", "项目", "云提供商", "余额", "上月总支出", "当前余额可使用月份（以上个月支出为标准）", "昨日支出", "前日支出", "日支出增长率"}

var balanceStyleCSS = `
    table {
    	text-align: center;
    	font-family: verdana, arial, sans-serif;
    	font-size: 11px;
    	color: #333333;
    	border-width: 1px;
    	border-color: #666666;
    	border-collapse: collapse;
    }
    
    table th {
    	border-width: 1px;
    	padding: 8px;
    	border-style: solid;
    	border-color: #666666;
    	background-color: #668F99;
    }
    
    table td {
    	border-width: 1px;
    	padding: 8px;
    	border-style: solid;
    	border-color: #666666;
    	background-color: #f2f2f2;
  }`

var sortedByBalance = table.SortBy{
	Name: "余额",
	Mode: table.AscNumeric,
}

var sortedByProject = table.SortBy{
	Name: "项目",
	Mode: table.Asc,
}
