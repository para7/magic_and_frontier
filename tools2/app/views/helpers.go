package views

import (
	"fmt"
	"strings"

	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/webui"
)

func FieldError(errs map[string]string, key string) string {
	if errs == nil {
		return ""
	}
	return errs[key]
}

func HasValue(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func NoticeClass(notice *webui.Notice) string {
	if notice == nil || notice.Kind == "" {
		return "notice-info"
	}
	switch notice.Kind {
	case "success":
		return "notice-success"
	case "error":
		return "notice-error"
	default:
		return "notice-info"
	}
}

func FloatText(value *float64) string {
	if value == nil {
		return "-"
	}
	return trimFloat(*value)
}

func trimFloat(value float64) string {
	text := fmt.Sprintf("%.4f", value)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" {
		return "0"
	}
	return text
}

func BoolText(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func SubmitLabel(editing bool) string {
	if editing {
		return "Update"
	}
	return "Create"
}

func FormTitle(label string, editing bool) string {
	if editing {
		return "Edit " + label
	}
	return "New " + label
}

func FormAction(basePath string, editing bool) string {
	if editing {
		return basePath + "/edit"
	}
	return basePath + "/new"
}

func TriggerText(entry enemyskills.EnemySkillEntry) string {
	if entry.Trigger == nil {
		return "-"
	}
	return string(*entry.Trigger)
}

func CooldownText(entry enemyskills.EnemySkillEntry) string {
	return FloatText(entry.Cooldown)
}

func JoinLines(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}
