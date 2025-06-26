package system

import (
	"testing"
)

func TestTgUserService_GetAllUserInfoList(t *testing.T) {

}

func TestTgUserService_RestTimes(t *testing.T) {

}

func TestTgUserService_SendFiles(t *testing.T) {

	//botToken := "7668068911:AAFOXuA7KpWOfur0rcoVbZTwGOgsBCjkI3s"
	//chatID := ":-4657809905"
	//filePath := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\log\\2025-03-21\\error.log"
	//bot, err := telego.NewBot(botToken, telego.WithDefaultLogger(true, true))
	//if err != nil {
	//	log.Fatalf("Failed to create bot: %s", err)
	//}
	//
	//// 2. 打开要发送的文件
	//filePath := "path/to/your/file.txt"
	//file, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("Failed to open file: %s", err)
	//}
	//defer file.Close()
	//
	//// 3. 使用 FileUpload 来上传文件
	//uploadedFile := telegofile.FileUpload{
	//	File:     file,
	//	FileName: "file.txt", // 发送时文件的名字
	//}
	//
	//// 4. 构建发送文件请求
	//chatID := telego.ChatID{ID: -123456789} // 群聊 ID，注意负号
	//params := telego.SendDocumentParams{
	//	ChatID:   chatID,
	//	Document: telego.InputFile{File: uploadedFile},
	//	Caption:  "Here is your file!",
	//}
	//
	//// 5. 发送文件
	//message, err := bot.SendDocument(params)
	//if err != nil {
	//	log.Fatalf("Failed to send document: %s", err)
	//}
	//
	//fmt.Printf("Message sent: %+v\n", message)
}
