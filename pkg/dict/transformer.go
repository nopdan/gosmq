package dict

type Transformer interface {
	Read(*Dict) []Entry
}

// 转换赛码表
func NewTransformer(format string) Transformer {
	var t Transformer
	switch format {
	case "jisu", "js":
		t = jisu{}
	case "duoduo", "dd":
		t = duoduo{}
	case "jidian", "jd":
		t = jidian{}
	case "bingling", "bl":
		t = duoduo{true}
	case "default":
		t = smb{}
	default:
		t = smb{}
	}
	return t
}
