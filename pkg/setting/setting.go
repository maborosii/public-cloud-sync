package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

func NewSetting(vp *viper.Viper) *Setting {
	return &Setting{vp}
}

func (s *Setting) ReadConfig(value interface{}) error {
	err := s.vp.Unmarshal(value)
	if err != nil {
		return err
	}
	return nil
}
