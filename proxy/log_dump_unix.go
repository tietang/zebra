// +build linux darwin

package proxy

import (
	log "github.com/sirupsen/logrus"
	//"os"
	//"os/exec"
	"os"
	"syscall"
)

func logDump() {
	logFile, err := os.OpenFile("./dump.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		log.Info("服务启动出错", "打开异常日志文件失败", err)
		return
	}
	// 将进程标准出错重定向至文件，进程崩溃时运行时将向该文件记录协程调用栈信息,
	// linux系统中的dup系统调用在windows系统及mac系统下暂时不可用
	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
}
