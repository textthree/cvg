package app

import "os"

// vim ~/.bashrc
// 文件末尾加入：export ENV=production
// 生效：source ~/.bashrc
func IsDevelop() bool {
	return os.Getenv("ENV") == "development"
}

func Env() string {
	return os.Getenv("ENV")
}
