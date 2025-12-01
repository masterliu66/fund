package util

import (
	"fmt"
	"fund/model"
	"math"
	"strings"
)

// FillValuationFields 根据已有指标为 FundInfoReport 补充估值状态、投资建议、预期收益和数据状态。
func FillValuationFields(r *model.FundInfoReport) {
	// 基本数据校验：缺少关键值时标记为未知
	if r == nil || r.Gsz <= 0 || r.Tp80MinDwjz <= 0 || r.Tp80MaxDwjz <= 0 {
		r.ValuationStatus = "未知"
		r.InvestAdvice = "数据不足，暂无法给出明确投资建议，请谨慎决策。"
		r.ExpectedReturnMin = 0
		r.ExpectedReturnMax = 0
		r.ExpectedReturnNote = "缺少必要估值分位数据，无法估算预期收益。"
		r.DataStatus = "MISSING"
		if r != nil && r.DataStatusNote == "" {
			r.DataStatusNote = "缺少 TP80 区间或当日估值数据。"
		}
		r.ValuationScore = 0
		return
	}

	r.DataStatus = "OK"
	if r.DataStatusNote == "" {
		r.DataStatusNote = ""
	}

	g := r.Gsz
	low80 := r.Tp80MinDwjz
	high80 := r.Tp80MaxDwjz
	low85 := r.Tp85MinDwjz
	high85 := r.Tp85MaxDwjz

	// 如果 TP85 未提供，则退化为使用 TP80
	if low85 <= 0 {
		low85 = low80
	}
	if high85 <= 0 {
		high85 = high80
	}

	// 估值状态判断：先判断严重低估/严重高估，再判断普通低估/高估
	switch {
	case g <= low85:
		r.ValuationStatus = "严重低估"
	case g >= high85:
		r.ValuationStatus = "严重高估"
	case g <= low80:
		r.ValuationStatus = "低估"
	case g >= high80:
		r.ValuationStatus = "高估"
	default:
		r.ValuationStatus = "正常"
	}

	// 投资建议与简单评分
	switch r.ValuationStatus {
	case "严重低估":
		r.InvestAdvice = "当前估值明显偏低，可考虑积极但分批建仓，并注意资金管理。"
		r.ValuationScore = 90
	case "低估":
		r.InvestAdvice = "当前估值处于历史偏低区间，可考虑分批建仓。"
		r.ValuationScore = 80
	case "正常":
		// 计算与低估/高估阈值的差距（仅在阈值有效时附加说明）
		baseAdvice := "当前估值处于历史正常区间，可结合自身风险偏好定投或继续持有。"
		var extra []string
		if low80 > 0 {
			gapLow := (g - low80) / low80 * 100
			if gapLow > 0 {
				extra = append(extra, fmt.Sprintf("高于低估阈值 %.1f%%", gapLow))
			}
		}
		if high80 > 0 {
			gapHigh := (high80 - g) / high80 * 100
			if gapHigh > 0 {
				extra = append(extra, fmt.Sprintf("低于高估阈值 %.1f%%", gapHigh))
			}
		}
		if len(extra) > 0 {
			r.InvestAdvice = baseAdvice + " （" + strings.Join(extra, "，") + "）"
		} else {
			r.InvestAdvice = baseAdvice
		}
		r.ValuationScore = 60
	case "高估":
		r.InvestAdvice = "当前估值偏高，短期买入性价比较低，建议谨慎或等待更合适的价格。"
		r.ValuationScore = 40
	case "严重高估":
		r.InvestAdvice = "当前估值明显偏高，建议以减仓或观望为主，避免在高位大量加仓。"
		r.ValuationScore = 20
	default:
		r.InvestAdvice = "数据不足，暂无法给出明确投资建议，请谨慎决策。"
		r.ValuationScore = 0
	}

	// 预期收益区间（简单规则：基准区间 + 状态调整）
	baseMin := 0.06 // 6%
	baseMax := 0.12 // 12%

	switch r.ValuationStatus {
	case "严重低估":
		r.ExpectedReturnMin = baseMin + 0.04 // 10%
		r.ExpectedReturnMax = baseMax + 0.04 // 16%
	case "低估":
		r.ExpectedReturnMin = baseMin + 0.02 // 8%
		r.ExpectedReturnMax = baseMax + 0.02 // 14%
	case "正常":
		r.ExpectedReturnMin = baseMin // 6%
		r.ExpectedReturnMax = baseMax // 12%
	case "高估":
		r.ExpectedReturnMin = math.Max(0, baseMin-0.02) // 4%
		r.ExpectedReturnMax = math.Max(0, baseMax-0.02) // 10%
	case "严重高估":
		r.ExpectedReturnMin = math.Max(0, baseMin-0.04) // 2%
		r.ExpectedReturnMax = math.Max(0, baseMax-0.04) // 8%
	default:
		r.ExpectedReturnMin = 0
		r.ExpectedReturnMax = 0
	}

	if r.ExpectedReturnMin > 0 && r.ExpectedReturnMax > 0 {
		r.ExpectedReturnNote = "基于历史表现与当前估值水平的粗略估算，仅供参考，不构成收益承诺。"
	} else {
		r.ExpectedReturnNote = "预期收益无法可靠估算，仅供参考。"
	}

	// 根据历史数据长度判断数据是否部分缺失
	if r.LastYearMaxDwjz == 0 && r.LastYearMinDwjz == 0 {
		if r.DataStatus == "OK" {
			r.DataStatus = "PARTIAL"
		}
		if r.DataStatusNote == "" {
			r.DataStatusNote = "历史数据少于 1 年，估值判断可靠性有限。"
		}
		if r.RiskNote == "" {
			r.RiskNote = "历史样本较少，短期波动风险较难评估。"
		}
	}
}
