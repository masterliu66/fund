package template

import (
	"fmt"
	"fund/model"
	"strconv"
	"strings"
)

func GetFundReportTemplate(reports []model.FundInfoReport) string {

	if len(reports) == 0 {
		return ""
	}

	msgs := make([]string, len(reports))
	for i, report := range reports {
		msgs[i] = report.ToString()
	}
	for _, msg := range msgs {
		fmt.Println(msg)
	}

	builder := &strings.Builder{}
	builder.WriteString("<table>")
	builder.WriteString("<tr>")
	builder.WriteString("<th>代码</th>")
	builder.WriteString("<th>名称</th>")
	builder.WriteString("<th>历史最高</th>")
	//builder.WriteString("<th>历史平均</th>")
	builder.WriteString("<th>历史最低</th>")
	builder.WriteString("<th>TP20</th>")
	builder.WriteString("<th>TP80</th>")
	builder.WriteString("<th>TP15</th>")
	builder.WriteString("<th>TP85</th>")
	builder.WriteString("<th>1年最高</th>")
	builder.WriteString("<th>1年最低</th>")
	builder.WriteString("<th>90天最高</th>")
	builder.WriteString("<th>90天最低</th>")
	builder.WriteString("<th>上月最高</th>")
	builder.WriteString("<th>上月最低</th>")
	builder.WriteString("<th>本月最高</th>")
	// builder.WriteString("<th>本月平均</th>")
	builder.WriteString("<th>本月最低</th>")
	builder.WriteString("<th>当日估值</th>")
	builder.WriteString("<th>当日涨幅</th>")
	builder.WriteString("</tr>")
	for _, msg := range msgs {
		builder.WriteString("<tr>")
		elements := strings.Split(msg, ",")
		for _, element := range elements {
			builder.WriteString("<td align='center'>")
			builder.WriteString(element)
			builder.WriteString("</td>")
		}
		builder.WriteString("</tr>")
	}
	builder.WriteString("</table>")

	return builder.String()
}

func GetFundRecordTemplate(records []model.FundRecordPO) string {

	if len(records) == 0 {
		return ""
	}

	msgs := make([]string, len(records))
	for i, record := range records {
		msgs[i] = record.ToString()
	}

	builder := &strings.Builder{}
	builder.WriteString("<table>")
	builder.WriteString("<tr>")
	builder.WriteString("<th>序号</th>")
	builder.WriteString("<th>日期</th>")
	builder.WriteString("<th>基金</th>")
	builder.WriteString("<th>操作类型</th>")
	builder.WriteString("<th>金额</th>")
	builder.WriteString("<th>成交单价</th>")
	builder.WriteString("<th>份额</th>")
	builder.WriteString("<th>涨幅</th>")
	builder.WriteString("<th>盈利</th>")
	builder.WriteString("<th>累积盈利</th>")
	builder.WriteString("<th>累积投入金额</th>")
	builder.WriteString("</tr>")
	for i, msg := range msgs {
		builder.WriteString("<tr>")
		builder.WriteString("<td align='center'>")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString("</td>")
		elements := strings.Split(msg, ",")
		for _, element := range elements {
			builder.WriteString("<td align='center'>")
			builder.WriteString(element)
			builder.WriteString("</td>")
		}
		builder.WriteString("</tr>")
	}
	builder.WriteString("</table>")

	return builder.String()
}
