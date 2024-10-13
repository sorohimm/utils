package cfg

func Var(v string) string {
	if cfg == nil {
		return ""
	}

	return cfg.prefix + v
}
