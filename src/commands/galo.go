package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"bytes"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"galo", "galolevel", "meugalo"},
		Run:       runGalo,
		Available: true,
		Cooldown:  7,
		Usage:     "j!galo",
		Help:      "Informação sobre seu galo",
		Category:  1,
	})
}

func runGalo(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg, args, session)
	if user.Bot {
		return
	}
	galo, _ := rinha.GetGaloDB(user.ID)
	if galo.Type == 0 {
		galoType := rinha.GetRandByType(rinha.Common)
		galo.Type = galoType
		rinha.SaveGaloDB(user.ID, galo)
	}
	skills := rinha.GetEquipedSkills(&galo)
	avatar, err := utils.DownloadImage(rinha.Sprites[0][galo.Type-1])

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}

	name := "Galo do " + user.Username
	if galo.Name != "" {
		name = galo.Name
	}

	radius := 50.0
	// Resize the images
	avatar = resize.Resize(uint(radius*2), uint(radius*2), avatar, resize.Lanczos3)

	file, err := os.Open("resources/wall.jpg")

	if err != nil {
		return
	}

	defer file.Close()
	img, err := jpeg.Decode(file)

	if err != nil {
		return
	}

	dc := gg.NewContext(320, 450)

	dc.DrawRoundedRectangle(0, 0, 320, 630, 10)
	dc.Fill()

	dc.DrawRoundedRectangle(0, 0, 320, 150, 10)
	dc.Clip()
	dc.DrawImage(img, 0, 0)
	dc.ResetClip()

	dc.DrawRoundedRectangle(0, 0, 320, 630, 10)
	dc.Clip()

	dc.SetRGB(0.8, 0.31, 0.31)
	dc.DrawRectangle(0, 140, 320, 40)
	dc.Fill()

	dc.DrawCircle(160, 70, radius+3)
	//	dc.SetRGB(1, 1, 1)
	color := rinha.Classes[galo.Type].Rarity.Color()
	dc.SetHexColor(fmt.Sprintf("%06x", color))
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(160, 70, radius)
	dc.Clip()
	dc.Fill()

	dc.DrawImage(avatar, int(160-radius), int(70-radius))

	dc.ResetClip()
	dc.SetRGB(0.3, 0.3, 0.3)

	err = dc.LoadFontFace("./resources/Raleway-Bold.ttf", 13)
	dc.SetRGB255(196, 196, 196)
	dc.DrawRoundedRectangle(10, 195, 300, 20, 10)
	dc.Fill()
	dc.SetRGB255(208, 80, 80)
	level := rinha.CalcLevel(galo.Xp)
	curLevelXP := float64(rinha.CalcXP(level))
	nextLevelXp := float64(rinha.CalcXP(level + 1))
	percentage := (float64(galo.Xp-int(curLevelXP)) * 100) / (nextLevelXp - curLevelXP)
	if percentage >= 3 {
		dc.DrawRoundedRectangle(10, 195, 300*(percentage/100), 20, 10)
		dc.Fill()
	}
	dc.SetRGB(0.3, 0.3, 0.3)
	dc.DrawStringAnchored("HABILIDADES EQUIPADAS", 320/2, 295, 0.5, 0.5)
	dc.DrawLine(10, 310, 310, 310)
	dc.Stroke()

	err = dc.LoadFontFace("./resources/Raleway-Light.ttf", 14)
	dc.SetRGB255(255, 255, 255)
	dc.DrawStringAnchored(fmt.Sprintf("%d/%d", galo.Xp-int(curLevelXP), int(nextLevelXp-curLevelXP)), 320/2, 203.5, 0.5, 0.5)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Tipo", 10, 240)
	dc.DrawString("Level", 10, 255)
	dc.DrawString("Item", 10, 270)

	dc.DrawStringAnchored(rinha.Classes[galo.Type].Name, 310, 240, 1, 0)
	dc.DrawStringAnchored(strconv.Itoa(rinha.CalcLevel(galo.Xp)), 310, 255, 1, 0)
	if len(galo.Items) > 0 {
		dc.DrawStringAnchored(rinha.Items[galo.Items[0]].Name, 310, 270, 1, 0)
	} else {
		dc.DrawStringAnchored("Nenhum", 310, 270, 1, 0)

	}

	for i, skill := range skills {
		rinhaSkill := rinha.Skills[galo.Type-1][skill]
		margin := float64(335 + (25 * i))
		text, _ := rinha.SkillToString(rinhaSkill)
		dc.DrawString(text, 10, margin)
		dc.DrawStringAnchored(fmt.Sprintf("Dano: %d - %d", rinhaSkill.Damage[0], rinhaSkill.Damage[1]-1), 310, margin, 1, 0)
	}

	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace("./resources/Raleway-Light.ttf", 22)
	dc.DrawStringAnchored(name, 320/2, 160, 0.5, 0.5)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(b.Bytes()), "galo.jpg", false},
		},
	})
}
