package voice

import "voice/service"

func main() {
	voice := make([]byte, 0)
	//此处放入录音文件的字节码形式
	service.Communication(voice)
}
