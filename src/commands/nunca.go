package commands

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"bytes"
	"io"
	"image/png"
	"asura/src/utils"
	"strings"
	"os"
	"github.com/nfnt/resize"
)

func init() {
	handler.Register(handler.Command{
		Aliases:    []string{"test"},
		Run:        runTest,
		Available:  true,
	})
}

// ke
func runTest(session disgord.Session, msg *disgord.Message, args []string) {
	// Loads the template of nunca.png in memory
	file, err := os.Open("resources/nunca.png")
	defer file.Close()

	// Decodes it into a image.Image
	img,err := png.Decode(file)
	if err != nil { return; }

	// Download user image	
	url := utils.GetImageURL(msg,args)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar,err  := utils.DownloadImage(replacer.Replace(url))
	if err != nil{
		msg.Reply(context.Background(), session, "Invalid image")
	}
	// Resize the images
	img  = resize.Resize(500, 500, img, resize.Lanczos3)
	avatar = resize.Resize(500, 350, avatar, resize.Lanczos3)

	// Here we draw the image to the "editor"
	dc := gg.NewContext(500, 500)
	dc.DrawImage(img, 0, 0)
	dc.DrawImage(avatar,0, 150)

	// And here we encode it to send
	var b bytes.Buffer
    pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	session.SendMsg(context.Background(), msg.ChannelID, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(b.Bytes()), "nuncalul.jpg", false},
		},
	})
}
